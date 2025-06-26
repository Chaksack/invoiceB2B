package database

import (
	"database/sql"
	"fmt"
	"log"

	"invoiceB2B/internal@v2/config"

	_ "github.com/lib/pq"
)

// DB holds the database connection
type DB struct {
	*sql.DB
}

// InitDB initializes the database connection
func InitDB(cfg *config.Config) (*DB, error) {
	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connection established successfully")

	return &DB{db}, nil
}

// Query executes a query and returns the result
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.DB.Query(query, args...)
}

// QueryRow executes a query and returns a single row
func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.DB.QueryRow(query, args...)
}

// Exec executes a query without returning rows
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.DB.Exec(query, args...)
}

// Begin starts a new transaction
func (db *DB) Begin() (*sql.Tx, error) {
	return db.DB.Begin()
} 