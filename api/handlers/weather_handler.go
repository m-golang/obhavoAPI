package handlers

import (
	"errors"
	"fmt"
	"havoAPI/api/helpers"
	"havoAPI/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WeatherHandler is a struct that handles weather-related operations.
// It interacts with a service layer that fetches weather data from an external API.
type WeatherHandler struct {
	weather services.WeatherAPIServiceInterface // Interface to interact with the weather API service
}

// NewWeatherHandler creates a new instance of WeatherHandler with the provided weather service.
// This function is typically used during handler setup in the routing layer.
func NewWeatherHandler(weather services.WeatherAPIServiceInterface) *WeatherHandler {
	return &WeatherHandler{weather: weather}
}

// WeatherData handles the retrieval of weather data for a specific location.
// It expects an API key and a query parameter (location) from the URL,
// performs authorization and fetches the weather data for the location.
func (service *WeatherHandler) WeatherData(c *gin.Context) {
	// Extract API key and query (location) from the request URL
	apiKey, query, err := helpers.GetParametersFromUrl(c)
	if err != nil {
		// If there is an issue with the parameters, respond with an error
		helpers.ClientError(c, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	// Authorize the API key
	_, err = service.weather.APIKeyAuthorization(apiKey)
	if err != nil {
		// Handle case where the API key is invalid or disabled
		if errors.Is(err, services.ErrAPIKeyNotFound) {
			helpers.ClientError(c, http.StatusUnauthorized, "API key has been disabled.")
			return
		}
		// For other errors, respond with a server error
		helpers.ServerError(c, err)
		return
	}

	// Fetch weather data based on the query (location)
	weatherData, err := service.weather.FetchWeatherData(query)
	if err != nil {
		// Handle case where no location is found
		if errors.Is(err, services.ErrNoLocationFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("%v", err))
			return
		}
		// Respond with a server error if another issue occurs
		helpers.ServerError(c, err)
		return
	}

	// Return the fetched weather data in the response
	c.JSON(http.StatusOK, gin.H{
		"location": weatherData, // Send the weather data for the location
	})
}

// BulkWeatherData handles the retrieval of weather data for multiple locations at once.
// It expects an API key and a list of locations from the request body.
func (service *WeatherHandler) BulkWeatherData(c *gin.Context) {
	// Extract the API key from the URL
	apiKey, err := helpers.GetParametersFromUrlForBulk(c)
	if err != nil {
		// If the API key extraction fails, respond with an error
		helpers.ClientError(c, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	// Authorize the API key
	_, err = service.weather.APIKeyAuthorization(apiKey)
	if err != nil {
		// Handle case where the API key is invalid or disabled
		if errors.Is(err, services.ErrAPIKeyNotFound) {
			helpers.ClientError(c, http.StatusUnauthorized, "API key has been disabled.")
			return
		}
		// For other errors, respond with a server error
		helpers.ServerError(c, err)
		return
	}

	// Parse the request body to extract the list of locations
	var locations LocationsForm
	if err := c.ShouldBindJSON(&locations); err != nil {
		// If binding fails, respond with validation errors
		helpers.RespondWithValidationErrors(c, err, locations)
		return
	}

	// Filter valid location queries to avoid unnecessary API calls
	qValues := helpers.FilterValidQValues(locations)

	// Fetch bulk weather data for the valid locations
	bulkWeatherData, notFoundList, err := service.weather.FetchBulkWeatherData(qValues)
	if err != nil {
		// If there is an error fetching the weather data, respond with a server error
		helpers.ServerError(c, err)
		return
	}

	// If there are locations not found, include them in the response
	if len(notFoundList) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"bulk":      bulkWeatherData, // Weather data for found locations
			"not_found": notFoundList,    // Locations that were not found
		})
		return
	}

	// If all locations were found, return just the weather data
	c.JSON(http.StatusOK, gin.H{
		"bulk": bulkWeatherData, // Send the bulk weather data
	})
}
