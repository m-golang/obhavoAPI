package main

import (
	"fmt"
	"havoAPI/api/config"
	"havoAPI/api/handlers"
	"havoAPI/api/routes"
	"havoAPI/internal/models"
	"havoAPI/internal/services"
	"log"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Load environment variables from the .env file
	// If this fails, log the error and terminate the program
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("failed to load .env file in main.go: %v", err)
	}

	// Load DB connection parameters from environment variables
	// If any variable is missing, log the error and terminate the program
	dbUserName, err := config.LoadEnvironmentVariable("DB_USER_NAME")
	if err != nil {
		log.Fatalf("%v", err)
	}

	dbUserPassword, err := config.LoadEnvironmentVariable("DB_USER_PASSWORD")
	if err != nil {
		log.Fatalf("%v", err)
	}

	dbName, err := config.LoadEnvironmentVariable("DB_NAME")
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Construct the Data Source Name (DSN) for the database connection
	// The DSN will be used to connect to the MySQL database
	dsn := fmt.Sprintf("%v:%v@/%v?parseTime=true", dbUserName, dbUserPassword, dbName)

	// Open a connection to the database
	// If the connection fails, log the error and terminate the program
	db, err := models.OpenDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // Ensure that the DB connection is closed when the program exits

	// Initialize the UserService with the database connection
	usersService := services.NewUsersService(db)
	// Initialize the UserHandler with the UserService
	usersHandler := handlers.NewUsersHandler(usersService)

	// Initialize the WeatherAPIService with the database connection
	weatherAPIService := services.NewWeatherAPIService(db)
	// Initialize the WeatherHandler with the WeatherAPIService
	weatherapiHandler := handlers.NewWeatherHandler(weatherAPIService)

	// Create the ServeHandlerWrapper to group UserHandler and WeatherHandler
	// This will be used to route requests to the appropriate handler
	serveHandlerWrapper := &routes.ServeHandlerWrapper{
		UserHandler:    usersHandler,
		WeatherHandler: weatherapiHandler,
	}

	// Initialize a new cron job to periodically update weather data in the Redis cache every 30 minutes
	cronJob := cron.New()
	_, err = cronJob.AddFunc("@every 30m", func() {
		// Update the weather data in the cache
		err := weatherAPIService.UpdateWeatherDataInTheRedisCache()
		if err != nil {
			// Log the error if the update fails
			log.Printf("Error updating weather data in cache: %v", err)
		} else {
			// Log a success message if the update is successful
			log.Println("Weather data updated successfully!")
		}
	})
	if err != nil {
		log.Fatal(err) // If adding the cron job fails, log the error and terminate
	}

	// Start the cron job in a separate goroutine to run it periodically
	go cronJob.Start()

	// Initialize the Gin router with the routes defined in the ServeHandlerWrapper
	router := routes.Route(serveHandlerWrapper)

	// Start the HTTP server in a separate goroutine to handle incoming requests
	go func() {
		if err := router.Run(); err != nil {
			// If there is an error starting the server, log the error and terminate
			log.Fatal("error running the server")
		}
	}()

	// Block the main goroutine indefinitely so that the application keeps running
	select {}
}
