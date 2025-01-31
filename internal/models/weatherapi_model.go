package models

import "fmt"

// DBContractWeatherapi defines the contract (interface) for database operations
// related to weather API keys. This ensures that any struct implementing this
// interface must provide an implementation for checking the validity of an API key.
type DBContractWeatherapi interface {
	CheckUserAPIKey(apiKey string) (bool, error) // Check if the provided API key exists in the database
}

// WeatherapiModel represents the struct that holds the database connection
// for weather-related API operations. It includes methods for interacting with
// the weather API data, specifically for checking the validity of API keys.
type WeatherapiModel struct {
	db DBContractWeatherapi // db represents the database connection (dependency injection)
}

// NewWeatherapiModel initializes a new WeatherapiModel instance with the provided
// database connection, ensuring that the database operations can be performed within the model.
func NewWeatherapiModel(db DBContractWeatherapi) *WeatherapiModel {
	return &WeatherapiModel{db: db}
}

// CheckUserAPIKey checks if the provided API key exists in the `api_keys` table in the database.
// It returns true if the API key is valid, or false if not. If an error occurs, it returns the error.
func (msql *MySQL) CheckUserAPIKey(apiKey string) (bool, error) {
	// SQL query to count how many rows in the api_keys table have the provided api_key
	stmt := `SELECT COUNT(*) FROM api_keys WHERE api_key=?`

	// Variable to store the count of matching rows
	var count int

	// Execute the query and scan the result into the 'count' variable
	err := msql.DB.QueryRow(stmt, apiKey).Scan(&count)
	if err != nil {
		// Return a wrapped error if something goes wrong during the query
		return false, fmt.Errorf("failed to scan count of api key in the database: %w", err)
	}

	// If the count is greater than 0, the API key is valid and exists in the database
	if count > 0 {
		return true, nil
	}
	// If no matching rows are found, return the custom error indicating the API key is not found
	return false, ErrAPIKeyNotFound
}
