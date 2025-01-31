package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"havoAPI/api/config"
	"havoAPI/internal/models"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// WeatherAPIServiceInterface defines the methods for interacting with weather data.
// These methods include fetching individual or bulk weather data, authorizing API keys,
// and updating weather data in a Redis cache.
type WeatherAPIServiceInterface interface {
	// FetchBulkWeatherData retrieves weather data for multiple locations.
	// It returns an array of formatted weather data and an array of locations not found.
	FetchBulkWeatherData(queries []string) ([]FormattedWeatherData, []string, error)

	// FetchWeatherData retrieves weather data for a single location.
	// It returns the formatted weather data or an error if the location is not found or the request fails.
	FetchWeatherData(query string) (FormattedWeatherData, error)

	// APIKeyAuthorization checks if the provided API key is valid for a user.
	// It returns true if the API key is valid, otherwise false along with an error if any.
	APIKeyAuthorization(apiKey string) (bool, error)

	// UpdateWeatherDataInTheRedisCache updates all weather data in the Redis cache.
	// This involves deleting the current cache and fetching new data for predefined locations.
	UpdateWeatherDataInTheRedisCache() error
}

// WeatherAPIService is a concrete implementation of the WeatherAPIServiceInterface.
// It interacts with both a database and a Redis client to fetch, cache, and manage weather data.
type WeatherAPIService struct {
	// db is an instance of the DBContractWeatherapi interface that handles database operations related to weather data.
	db models.DBContractWeatherapi

	// redisClient is a Redis client used for caching weather data.
	redisClient *redis.Client
}

// NewWeatherAPIService initializes a new instance of WeatherAPIService.
// It connects to a Redis instance using credentials loaded from environment variables.
func NewWeatherAPIService(db models.DBContractWeatherapi) *WeatherAPIService {
	// Load Redis address from the environment.
	redisAddr, err := config.LoadEnvironmentVariable("REDIS_ADDR")
	if err != nil {
		log.Fatal("failed to receive redis address from .env file")
	}

	// Load Redis password from the environment.
	redisPass, err := config.LoadEnvironmentVariable("REDIS_PASS")
	if err != nil {
		log.Fatal("failed to receive redis password from .env file")
	}

	// Initialize Redis client with the loaded credentials.
	rdb := redis.NewClient(&redis.Options{
		Addr:        redisAddr,
		Password:    redisPass,
		DB:          0,
		DialTimeout: 5 * time.Second,
	})

	// Return the newly created WeatherAPIService instance.
	return &WeatherAPIService{
		db:          db,
		redisClient: rdb,
	}
}

// FetchWeatherData retrieves weather data for a single location, either from the Redis cache or by querying the weather API.
// If data is not in the cache, it makes a request to the weather API and caches the result.
func (s *WeatherAPIService) FetchWeatherData(q string) (FormattedWeatherData, error) {
	// Capitalize the first letter of the location for consistent formatting.
	q = capitalizeFirstLetter(q)

	// Attempt to retrieve the weather data from Redis cache.
	cachedData, err := s.retrieveWeatherDataFromRedisCache(q)
	if errors.Is(err, nil) {
		// If data is found in the cache, return it.
		return cachedData, nil
	}

	// If no data is found in the cache, attempt to fetch it from the weather API.
	if errors.Is(err, ErrNoDataCache) {
		// Load the Weather API key from the environment.
		apiKeyForWeatherAPI, err := config.LoadEnvironmentVariable("API_KEY_FOR_WEATHERAPI")
		if err != nil {
			return FormattedWeatherData{}, err
		}

		// Format the query for the API request.
		query := strings.Replace(q, " ", "%20", -1)
		url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKeyForWeatherAPI, query)

		// Make the request to the weather API.
		resBody, err := requestToWeatherApi(url)
		if err != nil {
			// Return specific error if no location is found.
			if errors.Is(err, ErrNoLocationFound) {
				return FormattedWeatherData{}, ErrNoLocationFound
			}
			return FormattedWeatherData{}, err
		}

		// Parse the response body into a Weather struct.
		var weatherData Weather
		err = json.Unmarshal(resBody, &weatherData)
		if err != nil {
			// Handle JSON parsing errors.
			if _, ok := err.(*json.SyntaxError); ok {
				return FormattedWeatherData{}, ErrUnexpectedEndOfJSONInput
			}
			return FormattedWeatherData{}, fmt.Errorf("error occurred while unmarshaling JSON: %w", err)
		}

		// Format the weather data and cache it in Redis.
		formattedData := formatWeatherData(weatherData)
		err = s.cacheTheWeatherDataToRedis(query, formattedData)
		if err != nil {
			log.Fatalf("Error caching weather data: %v", err)
		}

		// Return the formatted weather data.
		return formattedData, nil
	}

	// Return an error if something else went wrong.
	return FormattedWeatherData{}, err
}

// FetchBulkWeatherData retrieves weather data for multiple locations, handling both found and not found locations.
func (s *WeatherAPIService) FetchBulkWeatherData(queries []string) ([]FormattedWeatherData, []string, error) {
	var bulkWeatherData []FormattedWeatherData
	var notFound []string

	// Loop through each query and attempt to fetch its weather data.
	for _, q := range queries {
		weatherData, err := s.FetchWeatherData(q)
		if err != nil {
			// If no location is found, add it to the notFound list.
			if errors.Is(err, ErrNoLocationFound) {
				notFound = append(notFound, fmt.Sprintf("'%s' not found", q))
				continue
			} else {
				return nil, nil, err
			}
		}
		// Append the found weather data to the result.
		bulkWeatherData = append(bulkWeatherData, weatherData)
	}

	// Return the bulk weather data and any locations that were not found.
	if len(notFound) > 0 {
		return bulkWeatherData, notFound, nil
	}
	return bulkWeatherData, nil, nil
}

// APIKeyAuthorization checks whether the provided API key is valid.
func (s *WeatherAPIService) APIKeyAuthorization(apiKey string) (bool, error) {
	// Check the validity of the API key by querying the database.
	isKeyTrue, err := s.db.CheckUserAPIKey(apiKey)
	if err != nil {
		// Return an error if the key is not found or another issue occurs.
		if errors.Is(err, models.ErrAPIKeyNotFound) {
			return false, ErrAPIKeyNotFound
		}
		return false, fmt.Errorf("error occurred while checking user API key: %w", err)
	}

	// Return true if the API key is valid, otherwise false.
	return isKeyTrue, nil
}

// requestToWeatherApi sends a GET request to the Weather API and returns the response body.
func requestToWeatherApi(url string) ([]byte, error) {
	// Send a GET request to the given URL.
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request to the given URL: %w", err)
	}
	defer response.Body.Close()

	// Check if the response indicates an error or invalid location.
	if response.StatusCode == http.StatusBadRequest {
		return nil, ErrNoLocationFound
	}

	// If the response status is not OK, return an error.
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error occurred: weatherapi response status code is not 200: %w", err)
	}

	// Read the response body.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error occurred while reading response body of weatherapi: %w", err)
	}

	// Return the response body.
	return body, nil
}

// cacheTheWeatherDataToRedis stores the weather data for a specific location in Redis.
func (s *WeatherAPIService) cacheTheWeatherDataToRedis(location string, weatherData FormattedWeatherData) error {
	// Marshal the weather data into JSON format.
	jsonData, err := json.Marshal(weatherData)
	if err != nil {
		return fmt.Errorf("failed to marshal weatherData: %w", err)
	}

	// Set the cached data in Redis with a 30-minute expiration time.
	err = s.redisClient.Set(context.Background(), location, jsonData, 30*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to set data in Redis: %w", err)
	}

	// Return nil if the operation was successful.
	return nil
}

// retrieveWeatherDataFromRedisCache attempts to fetch weather data from Redis cache for a location.
func (s *WeatherAPIService) retrieveWeatherDataFromRedisCache(location string) (FormattedWeatherData, error) {
	// Capitalize the first letter of the location for consistent formatting.
	location = capitalizeFirstLetter(location)

	// Attempt to get cached data from Redis.
	jsonData, err := s.redisClient.Get(context.Background(), location).Result()
	if err != nil {
		// Return an error if data is not found in the cache.
		if errors.Is(err, redis.Nil) {
			return FormattedWeatherData{}, ErrNoDataCache
		}
		return FormattedWeatherData{}, fmt.Errorf("failed to get data from Redis: %w", err)
	}

	// Unmarshal the cached data into a FormattedWeatherData object.
	var weatherData FormattedWeatherData
	err = json.Unmarshal([]byte(jsonData), &weatherData)
	if err != nil {
		return FormattedWeatherData{}, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	// Return the cached weather data.
	return weatherData, nil
}

// deleteAllWeatherDataFromRedisCache clears all weather data from the Redis cache.
func (s *WeatherAPIService) deleteAllWeatherDataFromRedisCache() error {
	// Flush the entire Redis database.
	err := s.redisClient.FlushDB(context.Background()).Err()
	if err != nil {
		return fmt.Errorf("failed to flush Redis database: %v", err)
	}
	return nil
}

// UpdateWeatherDataInTheRedisCache deletes the current weather data in Redis and updates it with new data
// for a predefined list of countries.
func (s *WeatherAPIService) UpdateWeatherDataInTheRedisCache() error {
	// Delete all existing weather data from Redis.
	err := s.deleteAllWeatherDataFromRedisCache()
	if err != nil {
		return err
	}

	// List of predefined countries to update weather data for.
	var country_list = []string{"Afghanistan", "Albania", "Algeria", "Andorra", "Angola", "Anguilla", "Antigua &amp; Barbuda", "Argentina", "Armenia", "Aruba", "Australia", "Austria", "Azerbaijan", "Bahamas", "Bahrain", "Bangladesh", "Barbados", "Belarus", "Belgium", "Belize", "Benin", "Bermuda", "Bhutan", "Bolivia", "Bosnia &amp; Herzegovina", "Botswana", "Brazil", "British Virgin Islands", "Brunei", "Bulgaria", "Burkina Faso", "Burundi", "Cambodia", "Cameroon", "Cape Verde", "Cayman Islands", "Chad", "Chile", "China", "Colombia", "Congo", "Cook Islands", "Costa Rica", "Cote D Ivoire", "Croatia", "Cruise Ship", "Cuba", "Cyprus", "Czech Republic", "Denmark", "Djibouti", "Dominica", "Dominican Republic", "Ecuador", "Egypt", "El Salvador", "Equatorial Guinea", "Estonia", "Ethiopia", "Falkland Islands", "Faroe Islands", "Fiji", "Finland", "France", "French Polynesia", "French West Indies", "Gabon", "Gambia", "Georgia", "Germany", "Ghana", "Gibraltar", "Greece", "Greenland", "Grenada", "Guam", "Guatemala", "Guernsey", "Guinea", "Guinea Bissau", "Guyana", "Haiti", "Honduras", "Hong Kong", "Hungary", "Iceland", "India", "Indonesia", "Iran", "Iraq", "Ireland", "Isle of Man", "Israel", "Italy", "Jamaica", "Japan", "Jersey", "Jordan", "Kazakhstan", "Kenya", "Kuwait", "Kyrgyz Republic", "Laos", "Latvia", "Lebanon", "Lesotho", "Liberia", "Libya", "Liechtenstein", "Lithuania", "Luxembourg", "Macau", "Macedonia", "Madagascar", "Malawi", "Malaysia", "Maldives", "Mali", "Malta", "Mauritania", "Mauritius", "Mexico", "Moldova", "Monaco", "Mongolia", "Montenegro", "Montserrat", "Morocco", "Mozambique", "Namibia", "Nepal", "Netherlands", "Netherlands Antilles", "New Caledonia", "New Zealand", "Nicaragua", "Niger", "Nigeria", "Norway", "Oman", "Pakistan", "Palestine", "Panama", "Papua New Guinea", "Paraguay", "Peru", "Philippines", "Poland", "Portugal", "Puerto Rico", "Qatar", "Reunion", "Romania", "Russia", "Rwanda", "Saint Pierre &amp; Miquelon", "Samoa", "San Marino", "Satellite", "Saudi Arabia", "Senegal", "Serbia", "Seychelles", "Sierra Leone", "Singapore", "Slovakia", "Slovenia", "South Africa", "South Korea", "Spain", "Sri Lanka", "St Kitts &amp; Nevis", "St Lucia", "St Vincent", "St. Lucia", "Sudan", "Suriname", "Swaziland", "Sweden", "Switzerland", "Syria", "Taiwan", "Tajikistan", "Tanzania", "Thailand", "Timor L'Este", "Togo", "Tonga", "Trinidad &amp; Tobago", "Tunisia", "Turkey", "Turkmenistan", "Turks &amp; Caicos", "Uganda", "Ukraine", "United Arab Emirates", "United Kingdom", "Uruguay", "Uzbekistan", "Venezuela", "Vietnam", "Virgin Islands (US)", "Yemen", "Zambia", "Zimbabwe"}

	// Fetch weather data for each country and cache it.
	for _, location := range country_list {
		_, err := s.FetchWeatherData(location)
		if err != nil {
			log.Printf("Error fetching data for %s: %v", location, err)
			continue
		}

		// Throttle the requests to avoid overwhelming the API.
		time.Sleep(500 * time.Millisecond)
	}

	// Return nil when the update process is complete.
	return nil
}
 