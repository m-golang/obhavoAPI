package model

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// DBContractUsers defines the contract (interface) for database operations
// related to users. This ensures that any struct implementing this interface
// must provide implementations for inserting users and retrieving user credentials.
type DBContractUsers interface {
	InsertUser(name, surname, username string, password_hash []byte) error
	RetrieveUserCredentials(username string) (int, string, error)
}

// UsersModel represents the struct that holds the database connection
// and methods related to user operations.
type UsersModel struct {
	db DBContractUsers // db represents the database connection
}

// NewUsersModel initializes a new UsersModel struct with a provided database connection.
// This function ensures that a UsersModel instance is always created with a valid DB connection.
func NewUsersModel(db DBContractUsers) *UsersModel {
	return &UsersModel{db: db}
}

// InsertUser inserts a new user into the database. If the username already exists,
// it returns an appropriate error. This method assumes that the username is unique
// in the database.
func (msql *MySQL) InsertUser(name, surname, username string, password_hash []byte) error {
	// SQL query to insert a new user into the 'users' table.
	stmt := `INSERT INTO users (name, surname, username, password_hash) VALUES(?, ?, ?, ?)`

	// Check if the error is a MySQL error and if the error code indicates a duplicate username.
	_, err := msql.DB.Exec(stmt, name, surname, username, password_hash)
	if err != nil {
		// Check if the error is due to a duplicate username
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 { // 1062 is the MySQL error code for duplicate entry
				// Return a custom error if the username already exists in the database.
				return ErrDuplicatedUsername
			}
		}
		// Return a wrapped error if the insert operation fails for any other reason.
		return fmt.Errorf("failed to insert the new user to the database: %w", err)
	}

	// Return nil if the user is successfully inserted.
	return nil
}

// RetrieveUserCredentials retrieves the credentials (user ID and hashed password)
// for a given username. If the user is not found, it returns an error indicating that.
// This method assumes the 'users' table contains the columns 'id' and 'password_hash'.
func (msql *MySQL) RetrieveUserCredentials(username string) (int, string, error) {
	// SQL query to retrieve user credentials based on the username.
	stmt := `SELECT id, password_hash FROM users WHERE username = ?`

	// Variables to store the retrieved user ID and password hash.
	var userID int
	var password_hash string

	// Query the database and scan the result into the userID and password_hash variables.
	err := msql.DB.QueryRow(stmt, username).Scan(&userID, &password_hash)
	if err != nil {
		// If no rows are returned (user not found), return a custom error indicating this.
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", ErrUserNotFound
		}
		// Return a wrapped error if any other error occurs during the query or scanning.
		return 0, "", fmt.Errorf("failed to scan user credentials: %w", err)
	}

	// Return the user ID and password hash if the user is found.
	return userID, password_hash, nil
}
