package configs

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DBInstance is our global database connection pool
var DBInstance *sql.DB

// ConnectDB initializes the connection to PostgreSQL
func ConnectDB() {
	// In a production app, these come from your secret .env file
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	var err error
	DBInstance, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	// Verify the connection works
	err = DBInstance.Ping()
	if err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	fmt.Println("Successfully connected to the logistics database!")
}