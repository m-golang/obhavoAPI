package routes

import (
	"havoAPI/api/handlers"
	"havoAPI/api/middlewares"

	"github.com/gin-gonic/gin"
)

// ServeHandlerWrapper wraps the UserHandler to provide HTTP handler functionality.
type ServeHandlerWrapper struct {
	*handlers.UserHandler // Embeds the UserHandler to use its methods for routing
	*handlers.WeatherHandler
}

// Route sets up the routes and handlers for the application.
// It accepts a ServeHandlerWrapper, which contains the logic for user-related actions like signup, login, and logout.
func Route(h *ServeHandlerWrapper) *gin.Engine {
	// Create a new Gin router with default middleware (logging, recovery, etc.)
	router := gin.Default()

	router.Use(middlewares.RecoverPanic())
	router.Use(middlewares.SecureHeaders())
	
	// Define version 1 of the API routes.
	v1 := router.Group("/v1")
	{
		// POST /v1/signup: Route for user signup
		v1.POST("/signup", h.Signup)

		// POST /v1/login: Route for user login
		v1.POST("/login", h.Login)

		// POST /v1/logout: Route for logging out; requires JWT authorization middleware
		v1.POST("/logout", middlewares.UserAuthorizationJWT(), h.Logout)

		v1.GET("/user/dashboard", middlewares.UserAuthorizationJWT(), h.UserDashboard)

		v1.GET("/weather", h.WeatherData)
		v1.POST("/weather", h.BulkWeatherData)
	}

	// Return the configured router to be used by the web server
	return router
}
