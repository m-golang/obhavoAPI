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
	cachedData, err := s.retrieveWeatherDataFromCache(q)
	if errors.Is(err, nil) {
		return cachedData, nil
	}

	if errors.Is(err, ErrNoDataCache) {
		apiKeyForWeatherAPI, err := config.LoadEnvironmentVariable("API_KEY_FOR_WEATHERAPI")
		if err != nil {
			return FormattedWeatherData{}, err
		}

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

		err = s.cacheTheWeatherData(formattedData.Name, formattedData)
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

func (s *WeatherAPIService) cacheTheWeatherData(location string, weatherData FormattedWeatherData) error {

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

func (s *WeatherAPIService) retrieveWeatherDataFromCache(location string) (FormattedWeatherData, error) {

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
