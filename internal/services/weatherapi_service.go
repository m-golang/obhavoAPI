package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"havoAPI/api/config"
	"havoAPI/internal/models"
	"io"
	"net/http"
	"strings"
)

type WeatherAPIServiceInterface interface {
	FetchWeatherData(query string) (*FormattedWeatherData, error)
	APIKeyAuthorization(apiKey string) (bool, error)
}

type WeatherAPIService struct {
	db models.DBContractWeatherapi
}

func NewWeatherAPIService(db models.DBContractWeatherapi) *WeatherAPIService {
	return &WeatherAPIService{db: db}
}

func (s *WeatherAPIService) FetchWeatherData(q string) (*FormattedWeatherData, error) {
	apiKeyForWeatherAPI, err := config.LoadEnvironmentVariable("API_KEY_FOR_WEATHERAPI")
	if err != nil {
		return &FormattedWeatherData{}, err
	}

	query := strings.Replace(q, " ", "%20", -1)

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKeyForWeatherAPI, query)
	resBody, err := requestToWeatherApi(url)
	if err != nil {
		if errors.Is(err, ErrNoLocationFound) {
			return &FormattedWeatherData{}, ErrNoLocationFound
		}
		return &FormattedWeatherData{}, err
	}
	var weatherData Weather

	err = json.Unmarshal(resBody, &weatherData)
	if err != nil {
		if _, ok := err.(*json.SyntaxError); ok {
			return nil, ErrUnexpectedEndOfJSONInput
		}
		return nil, fmt.Errorf("error occured while unmarsheling json: %w", err)
	}

	formattedData := formatWeatherData(weatherData)

	return &formattedData, nil
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
