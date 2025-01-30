package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// MySQL represents a connection to a MySQL database.
type MySQL struct {
	DB *sql.DB // The underlying database connection.
}

// OpenDB initializes and opens a connection to the MySQL database using the provided DSN (Data Source Name).
// It returns a pointer to a MySQL instance or an error if the connection fails.
func OpenDB(dsn string) (*MySQL, error) {
	// Attempt to open a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Ping the database to ensure the connection is valid
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Return a MySQL instance wrapping the connection
	return &MySQL{DB: db}, nil
}

// Close attempts to close the MySQL database connection.
// If an error occurs during closure, it logs the error and terminates the program.
func (mysql *MySQL) Close() {
	// Attempt to close the database connection
	err := mysql.DB.Close()
	if err != nil {
		log.Fatalf("mysql close failure: %v", err) // Fatal log if closing fails
	}
}
