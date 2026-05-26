package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ibrah/logistics-tracking-api/configs"
	"github.com/ibrah/logistics-tracking-api/internal/handlers" // Imported handlers package
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting the Delivery & Logistics Engine...")

	// 1. Load the .env file from the root directory
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, relying on system environment variables")
	}

	// 2. Initialize our database connection pool and run migrations
	configs.ConnectDB()

	// 3. Set up a basic health check route to test our API
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "service": "delivery-api"}`))
	})

	// 4. Register the user registration endpoint
	http.HandleFunc("/register", handlers.RegisterUserHandler)

	// 5. Start the server on port 8080
	port := ":8080"
	fmt.Printf("Backend API server listening on port %s 🚀\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}