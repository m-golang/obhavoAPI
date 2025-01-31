package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// DBContractUsers defines the contract (interface) for database operations
// related to users. Any struct that implements this interface must provide
// implementations for inserting users and retrieving user credentials.
type DBContractUsers interface {
	InsertUser(name, surname, username string, password_hash []byte) (int, error)
	RetrieveUserCredentials(username string) (int, string, error)
	InsertUserAPIKey(userID int, apiKey string) error
	CheckUserAPIKey(apiKey string) (bool, error)
	RetriveUserAPIKey(userID int) (string, error)
}

// UsersModel represents the struct that holds the database connection
// and provides methods for user-related operations in the database.
type UsersModel struct {
	db DBContractUsers // db represents the database connection (dependency injection)
}

// NewUsersModel initializes a new UsersModel instance with the provided
// database connection. This ensures that any database-related operations
// are encapsulated within this model.
func NewUsersModel(db DBContractUsers) *UsersModel {
	return &UsersModel{db: db}
}

// InsertUser inserts a new user into the database. It checks for duplicate
// usernames to ensure that no two users can have the same username.
// Returns the newly created user's ID, or an error if the operation fails.
func (msql *MySQL) InsertUser(name, surname, username string, password_hash []byte) (int, error) {
	// SQL query to insert a new user into the 'users' table
	stmt := `INSERT INTO users (name, surname, username, password_hash) VALUES(?, ?, ?, ?)`

	// Execute the insert operation, returning an error if it fails
	req, err := msql.DB.Exec(stmt, name, surname, username, password_hash)
	if err != nil {
		// Check for MySQL-specific error: duplicate username
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 { // 1062 is MySQL's error code for duplicate entry
				// Return a custom error indicating the username already exists
				return 0, ErrDuplicatedUsername
			}
		}
		// Return a wrapped error with the original failure reason
		return 0, fmt.Errorf("failed to insert the new user to the database: %w", err)
	}

	// Retrieve the ID of the newly inserted user
	userId, err := req.LastInsertId()
	if err != nil {
		// Return an error if unable to retrieve the last inserted ID
		return 0, fmt.Errorf("failed to retrieve last inserted user id: %w", err)
	}

	// Return the user ID and nil indicating a successful insert
	return int(userId), nil
}

// RetrieveUserCredentials retrieves the credentials (user ID and password hash)
// for a given username. If the user is not found, it returns an error.
// This method assumes the 'users' table contains 'id' and 'password_hash' columns.
func (msql *MySQL) RetrieveUserCredentials(username string) (int, string, error) {
	// SQL query to retrieve user credentials based on the username
	stmt := `SELECT id, password_hash FROM users WHERE username = ?`

	// Variables to store the retrieved user ID and password hash
	var userID int
	var password_hash string

	// Query the database and scan the result into userID and password_hash
	err := msql.DB.QueryRow(stmt, username).Scan(&userID, &password_hash)
	if err != nil {
		// If no rows are returned (user not found), return a custom error
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", ErrUserNotFound
		}
		// Return a wrapped error if any other error occurs during the query
		return 0, "", fmt.Errorf("failed to scan user credentials: %w", err)
	}

	// Return the user ID and password hash if found
	return userID, password_hash, nil
}

// InsertUserAPIKey inserts a new API key into the `api_keys` table for the specified user.
// It associates the provided user ID with the given API key in the database.
func (msql *MySQL) InsertUserAPIKey(userID int, apiKey string) error {
	// SQL query to insert the user ID and API key into the api_keys table
	stmt := `INSERT INTO api_keys (user_id, api_key) VALUES (?, ?)`

	// Execute the insert statement with the userID and apiKey values
	_, err := msql.DB.Exec(stmt, userID, apiKey)
	if err != nil {
		// Return a wrapped error indicating failure to insert the API key
		return fmt.Errorf("failed to insert new API key into the database: %w", err)
	}

	// Return nil if the insert operation is successful
	return nil
}

// RetriveUserAPIKey retrieves the API key for a given user ID from the `api_keys` table.
// If no API key is found for the user, it returns an error.
func (msql *MySQL) RetriveUserAPIKey(userID int) (string, error) {
	// SQL query to retrieve the API key for the given user ID
	stmt := `SELECT api_key FROM api_keys WHERE user_id = ?`

	// Variable to store the retrieved API key
	var apiKey string

	// Query the database and scan the result into apiKey
	err := msql.DB.QueryRow(stmt, userID).Scan(&apiKey)
	if err != nil {
		// Return a wrapped error if the retrieval fails
		return "", fmt.Errorf("failed to retrieve user API key: %w", err)
	}

	// Return the retrieved API key
	return apiKey, nil
}
