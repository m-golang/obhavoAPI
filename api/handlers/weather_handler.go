package handlers

import (
	"errors"
	"fmt"
	"havoAPI/api/helpers"
	"havoAPI/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WeatherHandler struct {
	weather services.WeatherAPIServiceInterface
}

func NewWeatherHandler(weather services.WeatherAPIServiceInterface) *WeatherHandler {
	return &WeatherHandler{weather: weather}
}

func (service *WeatherHandler) WeatherData(c *gin.Context) {
	apiKey, query, err := helpers.GetParametersFromUrl(c)
	if err != nil {
		helpers.ClientError(c, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	_, err = service.weather.APIKeyAuthorization(apiKey)
	if err != nil {
		if errors.Is(err, services.ErrAPIKeyNotFound) {
			helpers.ClientError(c, http.StatusUnauthorized, "API key has been disabled.")
			return
		}
		helpers.ServerError(c, err)
		return
	}

	weatherData, err := service.weather.FetchWeatherData(query)
	if err != nil {
		if errors.Is(err, services.ErrNoLocationFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("%v", err))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"location": weatherData,
	})
}
