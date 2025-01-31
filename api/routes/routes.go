package routes

import (
	"havoAPI/api/handlers"
	"havoAPI/api/middlewares"

	"github.com/gin-gonic/gin"
)

// ServeHandlerWrapper wraps the UserHandler and WeatherHandler to provide HTTP handler functionality.
// By embedding these handlers, the wrapper allows easy access to user and weather-related routes in the application.
type ServeHandlerWrapper struct {
	*handlers.UserHandler    // Embeds the UserHandler to handle user-related actions (signup, login, etc.)
	*handlers.WeatherHandler // Embeds the WeatherHandler to handle weather-related actions (weather data retrieval, bulk queries, etc.)
}

// Route sets up the routes and handlers for the application.
// It accepts a ServeHandlerWrapper, which contains the logic for user-related actions like signup, login, and logout,
// as well as weather data retrieval and bulk requests.
func Route(h *ServeHandlerWrapper) *gin.Engine {
	// Create a new Gin router with default middleware (logging, recovery, etc.)
	router := gin.Default()

	// Apply middleware for panic recovery, secure headers, and rate limiting
	router.Use(middlewares.RecoverPanic())  // Handles panics during request processing
	router.Use(middlewares.SecureHeaders()) // Adds security-related headers to the response
	router.Use(middlewares.RateLimiter())   // Limits the rate of incoming requests

	// Define version 1 of the API routes with the /v1 prefix
	v1 := router.Group("/api/v1")
	{
		// POST /v1/signup: Route for user signup
		// This route accepts user details, validates them, and creates a new user.
		v1.POST("/signup", h.Signup)

		// POST /v1/login: Route for user login
		// This route validates the user credentials and generates a JWT token upon successful authentication.
		v1.POST("/login", h.Login)

		// POST /v1/logout: Route for user logout, requires JWT authorization middleware
		// This route allows the user to log out and clear their session by removing the JWT token.
		v1.POST("/logout", middlewares.UserAuthorizationJWT(), h.Logout)

		// GET /v1/user/dashboard: Route to fetch user dashboard details, requires JWT authorization
		// This route provides user-specific data (e.g., API key) for the logged-in user.
		v1.GET("/user/dashboard", middlewares.UserAuthorizationJWT(), h.UserDashboard)

		// GET /v1/weather: Route for fetching weather data based on query parameter
		// This route returns weather data for a given location.
		v1.GET("/weather.current", h.WeatherData)

		// POST /v1/weather: Route for bulk weather data requests
		// This route accepts a list of locations and fetches weather data for each location.
		v1.POST("/weather.current", h.BulkWeatherData)
	}

	// Return the configured router to be used by the web server
	// This allows the Gin engine to process requests according to the defined routes and handlers.
	return router
}
