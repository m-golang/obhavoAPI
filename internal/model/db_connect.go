package model

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Mysql representes a connection to a MYSQL database.
type Mysql struct {
	DB *sql.DB
}

func OpenDB(dsn string) (*Mysql, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &Mysql{DB: db}, nil
}

func (mysql *Mysql) Close() {
	// Attempt to close the database connection
	err := mysql.DB.Close()
	if err != nil {
		log.Fatalf("mysql close failer: %v", err) // Fatal log if closing fails
	}
}
