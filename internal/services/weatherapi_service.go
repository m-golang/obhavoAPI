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

	var formattedData FormattedWeatherData
	formattedData.Name = weatherData.Location.Name
	formattedData.Country = weatherData.Location.Country
	formattedData.Lat = weatherData.Location.Lat
	formattedData.Lon = weatherData.Location.Lon

	formattedData.TempC = weatherData.Current.TempC
	if formattedData.TempC < -20 {
		formattedData.TempColor = "#003366" // Deep Blue
	} else if formattedData.TempC >= -20 && formattedData.TempC < -10 {
		formattedData.TempColor = "#4A90E2" // Ice Blue
	} else if formattedData.TempC >= -10 && formattedData.TempC < 0 {
		formattedData.TempColor = "#B3DFFD" // Light Blue
	} else if formattedData.TempC >= 0 && formattedData.TempC < 10 {
		formattedData.TempColor = "#E6F7FF" // Pale Grayish Blue
	} else if formattedData.TempC >= 10 && formattedData.TempC < 20 {
		formattedData.TempColor = "#D1F2D3" // Light Green
	} else if formattedData.TempC >= 20 && formattedData.TempC < 30 {
		formattedData.TempColor = "#FFFACD" // Soft Yellow
	} else if formattedData.TempC >= 30 && formattedData.TempC < 40 {
		formattedData.TempColor = "#FFCC80" // Light Orange
	} else if formattedData.TempC >= 40 && formattedData.TempC < 50 {
		formattedData.TempColor = "#FF7043" // Deep Orange
	} else if formattedData.TempC >= 50 {
		formattedData.TempColor = "#D32F2F" // Bright Red
	}

	formattedData.WindKph = weatherData.Current.WindKph
	if formattedData.WindKph >= 0 && formattedData.WindKph < 10 {
		formattedData.WindColor = "#E0F7FA" // Light Cyan
	} else if formattedData.WindKph >= 10 && formattedData.WindKph < 20 {
		formattedData.WindColor = "#B2EBF2" // Pale Blue
	} else if formattedData.WindKph >= 20 && formattedData.WindKph < 40 {
		formattedData.WindColor = "#4DD0E1" // Soft Teal
	} else if formattedData.WindKph >= 40 && formattedData.WindKph < 60 {
		formattedData.WindColor = "#0288D1" // Bright Blue
	} else if formattedData.WindKph >= 60 {
		formattedData.WindColor = "#01579B" // Deep Navy Blue
	}

	formattedData.Cloud = weatherData.Current.Cloud
	if formattedData.Cloud >= 0 && formattedData.Cloud < 10 {
		formattedData.CloudColor = "#FFF9C4" // Light Yellow
	} else if formattedData.Cloud >= 10 && formattedData.Cloud < 30 {
		formattedData.CloudColor = "#FFF176" // Soft Yellow
	} else if formattedData.Cloud >= 30 && formattedData.Cloud < 60 {
		formattedData.CloudColor = "#E0E0E0" // Light Gray
	} else if formattedData.Cloud >= 60 && formattedData.Cloud < 90 {
		formattedData.CloudColor = "#9E9E9E" // Gray
	} else if formattedData.Cloud >= 90 && formattedData.Cloud <= 100 {
		formattedData.CloudColor = "#616161" // Dark Gray
	}

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
