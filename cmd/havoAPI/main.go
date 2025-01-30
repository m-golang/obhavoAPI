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
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("failed to load .env file in mian.go: %v", err)
	}

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

	dsn := fmt.Sprintf("%v:%v@/%v?parseTime=true", dbUserName, dbUserPassword, dbName)

	db, err := models.OpenDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	usersService := services.NewUsersService(db)
	usersHandler := handlers.NewUsersHandler(usersService)

	weatherAPIService := services.NewWeatherAPIService(db)
	weatherapiHandler := handlers.NewWeatherHandler(weatherAPIService)

	serveHandlerWrapper := &routes.ServeHandlerWrapper{
		UserHandler:    usersHandler,
		WeatherHandler: weatherapiHandler,
	}
	router := routes.Route(serveHandlerWrapper)

	router.Run()

}
