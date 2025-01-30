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

type WeatherAPIServiceInterface interface {
	FetchBulkWeatherData(queries []string) ([]FormattedWeatherData, []string, error)
	FetchWeatherData(query string) (FormattedWeatherData, error)
	APIKeyAuthorization(apiKey string) (bool, error)
	UpdateWeatherDataInTheRedisCache() error
}

type WeatherAPIService struct {
	db          models.DBContractWeatherapi
	redisClient *redis.Client
}

func NewWeatherAPIService(db models.DBContractWeatherapi) *WeatherAPIService {
	redisAddr, err := config.LoadEnvironmentVariable("REDIS_ADDR")
	if err != nil {
		log.Fatal("failed to recieve redis addres from .env file")
	}

	redisPass, err := config.LoadEnvironmentVariable("REDIS_PASS")
	if err != nil {
		log.Fatal("failed to recieve redis password from .env file")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:        redisAddr,
		Password:    redisPass,
		DB:          0,
		DialTimeout: 5 * time.Second,
	})

	return &WeatherAPIService{
		db:          db,
		redisClient: rdb}
}

func (s *WeatherAPIService) FetchWeatherData(q string) (FormattedWeatherData, error) {
	cachedData, err := s.retrieveWeatherDataFromRedisCache(q)
	if errors.Is(err, nil) {
		return cachedData, nil
	}

	if errors.Is(err, ErrNoDataCache) {
		apiKeyForWeatherAPI, err := config.LoadEnvironmentVariable("API_KEY_FOR_WEATHERAPI")
		if err != nil {
			return FormattedWeatherData{}, err
		}
		log.Println("No cache found")
		query := strings.Replace(q, " ", "%20", -1)

		url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKeyForWeatherAPI, query)
		resBody, err := requestToWeatherApi(url)
		if err != nil {
			if errors.Is(err, ErrNoLocationFound) {
				return FormattedWeatherData{}, ErrNoLocationFound
			}
			return FormattedWeatherData{}, err
		}
		var weatherData Weather

		err = json.Unmarshal(resBody, &weatherData)
		if err != nil {
			if _, ok := err.(*json.SyntaxError); ok {
				return FormattedWeatherData{}, ErrUnexpectedEndOfJSONInput
			}
			return FormattedWeatherData{}, fmt.Errorf("error occured while unmarsheling json: %w", err)
		}

		formattedData := formatWeatherData(weatherData)

		err = s.cacheTheWeatherDataToRedis(formattedData.Name, formattedData)
		if err != nil {
			log.Fatalf("Error caching weather data: %v", err)
		}

		return formattedData, nil
	}

	return FormattedWeatherData{}, err
}

func (s *WeatherAPIService) FetchBulkWeatherData(queries []string) ([]FormattedWeatherData, []string, error) {

	var bulkWeatherData []FormattedWeatherData
	var notFound []string

	for _, q := range queries {
		weatherData, err := s.FetchWeatherData(q)
		if err != nil {
			if errors.Is(err, ErrNoLocationFound) {
				notFound = append(notFound, fmt.Sprintf("'%s' not found", q))
				continue
			} else {
				return nil, nil, err
			}
		}
		bulkWeatherData = append(bulkWeatherData, weatherData)
	}
	if len(notFound) > 0 {
		return bulkWeatherData, notFound, nil

	}
	return bulkWeatherData, nil, nil
}

func (s *WeatherAPIService) APIKeyAuthorization(apiKey string) (bool, error) {
	isKeyTrue, err := s.db.CheckUserAPIKey(apiKey)
	if err != nil {
		if errors.Is(err, models.ErrAPIKeyNotFound) {
			return false, ErrAPIKeyNotFound
		}
		return false, fmt.Errorf("error occured while checking user API key: %w", err)
	}

	return isKeyTrue, nil
}

func requestToWeatherApi(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request to the given URL: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusBadRequest {
		return nil, ErrNoLocationFound
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error occured: weatherapi response status code is not 200: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error occured while reading response body of weatherapi: %w", err)
	}

	return body, nil
}

func (s *WeatherAPIService) cacheTheWeatherDataToRedis(location string, weatherData FormattedWeatherData) error {

	jsonData, err := json.Marshal(weatherData)
	if err != nil {
		return fmt.Errorf("failed to marshal weatherData: %w", err)
	}

	err = s.redisClient.Set(context.Background(), location, jsonData, 30*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to set data in Redis: %w", err)
	}

	return nil
}

func (s *WeatherAPIService) retrieveWeatherDataFromRedisCache(location string) (FormattedWeatherData, error) {

	location = capitalizeFirstLetter(location)

	jsonData, err := s.redisClient.Get(context.Background(), location).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return FormattedWeatherData{}, ErrNoDataCache
		}
		return FormattedWeatherData{}, fmt.Errorf("failed to get data from Redis: %w", err)
	}

	var weatherData FormattedWeatherData
	err = json.Unmarshal([]byte(jsonData), &weatherData)
	if err != nil {
		return FormattedWeatherData{}, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return weatherData, nil
}

func (s *WeatherAPIService) deleteAllWeatherDataFromRedisCache() error {
	err := s.redisClient.FlushDB(context.Background()).Err()
	if err != nil {
		return fmt.Errorf("failed to flush Redis database: %v", err)
	}
	log.Printf("Successfully deleted weather data for all locations")
	return nil
}

func (s *WeatherAPIService) UpdateWeatherDataInTheRedisCache() error {
	err := s.deleteAllWeatherDataFromRedisCache()
	if err != nil {
		return err
	}

	var country_list = []string{"Afghanistan", "Albania", "Algeria", "Andorra", "Angola", "Anguilla", "Antigua &amp; Barbuda", "Argentina", "Armenia", "Aruba", "Australia", "Austria", "Azerbaijan", "Bahamas", "Bahrain", "Bangladesh", "Barbados", "Belarus", "Belgium", "Belize", "Benin", "Bermuda", "Bhutan", "Bolivia", "Bosnia &amp; Herzegovina", "Botswana", "Brazil", "British Virgin Islands", "Brunei", "Bulgaria", "Burkina Faso", "Burundi", "Cambodia", "Cameroon", "Cape Verde", "Cayman Islands", "Chad", "Chile", "China", "Colombia", "Congo", "Cook Islands", "Costa Rica", "Cote D Ivoire", "Croatia", "Cruise Ship", "Cuba", "Cyprus", "Czech Republic", "Denmark", "Djibouti", "Dominica", "Dominican Republic", "Ecuador", "Egypt", "El Salvador", "Equatorial Guinea", "Estonia", "Ethiopia", "Falkland Islands", "Faroe Islands", "Fiji", "Finland", "France", "French Polynesia", "French West Indies", "Gabon", "Gambia", "Georgia", "Germany", "Ghana", "Gibraltar", "Greece", "Greenland", "Grenada", "Guam", "Guatemala", "Guernsey", "Guinea", "Guinea Bissau", "Guyana", "Haiti", "Honduras", "Hong Kong", "Hungary", "Iceland", "India", "Indonesia", "Iran", "Iraq", "Ireland", "Isle of Man", "Israel", "Italy", "Jamaica", "Japan", "Jersey", "Jordan", "Kazakhstan", "Kenya", "Kuwait", "Kyrgyz Republic", "Laos", "Latvia", "Lebanon", "Lesotho", "Liberia", "Libya", "Liechtenstein", "Lithuania", "Luxembourg", "Macau", "Macedonia", "Madagascar", "Malawi", "Malaysia", "Maldives", "Mali", "Malta", "Mauritania", "Mauritius", "Mexico", "Moldova", "Monaco", "Mongolia", "Montenegro", "Montserrat", "Morocco", "Mozambique", "Namibia", "Nepal", "Netherlands", "Netherlands Antilles", "New Caledonia", "New Zealand", "Nicaragua", "Niger", "Nigeria", "Norway", "Oman", "Pakistan", "Palestine", "Panama", "Papua New Guinea", "Paraguay", "Peru", "Philippines", "Poland", "Portugal", "Puerto Rico", "Qatar", "Reunion", "Romania", "Russia", "Rwanda", "Saint Pierre &amp; Miquelon", "Samoa", "San Marino", "Satellite", "Saudi Arabia", "Senegal", "Serbia", "Seychelles", "Sierra Leone", "Singapore", "Slovakia", "Slovenia", "South Africa", "South Korea", "Spain", "Sri Lanka", "St Kitts &amp; Nevis", "St Lucia", "St Vincent", "St. Lucia", "Sudan", "Suriname", "Swaziland", "Sweden", "Switzerland", "Syria", "Taiwan", "Tajikistan", "Tanzania", "Thailand", "Timor L'Este", "Togo", "Tonga", "Trinidad &amp; Tobago", "Tunisia", "Turkey", "Turkmenistan", "Turks &amp; Caicos", "Uganda", "Ukraine", "United Arab Emirates", "United Kingdom", "Uruguay", "Uzbekistan", "Venezuela", "Vietnam", "Virgin Islands (US)", "Yemen", "Zambia", "Zimbabwe"}
	for _, location := range country_list {
		_, err := s.FetchWeatherData(location)
		if err != nil {
			log.Printf("Error fetching data for %s: %v", location, err)
			continue
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}
