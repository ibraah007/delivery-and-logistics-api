package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ibrah/logistics-tracking-api/configs"
	"github.com/joho/godotenv" // 1. Import the env loader
)

func main() {
	fmt.Println("Starting the Delivery & Logistics Engine...")

	// 2. Load the .env file from the root directory
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, relying on system environment variables")
	}

	// 3. Initialize our database connection pool
	configs.ConnectDB()

	// 4. Set up a basic health check route to test our API
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "service": "delivery-api"}`))
	})

	// 5. Start the server on port 8080
	port := ":8080"
	fmt.Printf("Backend API server listening on port %s 🚀\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}