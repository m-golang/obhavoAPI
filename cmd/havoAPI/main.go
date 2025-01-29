package main

import (
	"fmt"
	"havoAPI/api/config"
	"havoAPI/api/handlers"
	"havoAPI/api/routes"
	"havoAPI/internal/model"
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

	db, err := model.OpenDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	usersService := services.NewUsersService(db)
	usersHandler := handlers.NewUsersHandler(usersService)

	serveHandlerWrapper := &routes.ServeHandlerWrapper{
		UserHandler: usersHandler,
	}
	router := routes.Route(serveHandlerWrapper)

	router.Run()

}
