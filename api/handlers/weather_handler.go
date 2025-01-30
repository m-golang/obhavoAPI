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

func (service *WeatherHandler) BulkWeatherData(c *gin.Context) {
	apiKey, err := helpers.GetParametersFromUrlForBulk(c)
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

	var locations LocationsForm

	if err := c.ShouldBindJSON(&locations); err != nil {
		helpers.RespondWithValidationErrors(c, err, locations)
		return
	}

	qValues := helpers.FilterValidQValues(locations)

	bulkWeatherData, notFoundList, err := service.weather.FetchBulkWeatherData(qValues)
	if err != nil {
		helpers.ServerError(c, err)
		return
	}

	if len(notFoundList) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"bulk":      bulkWeatherData,
			"not_found": notFoundList,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"bulk": bulkWeatherData,
	})
}
