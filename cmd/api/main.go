package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ibrah/logistics-tracking-api/configs"
	"github.com/ibrah/logistics-tracking-api/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting the Delivery & Logistics Engine...")

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, relying on system environment variables")
	}

	configs.ConnectDB()

	// Health Check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "service": "delivery-api"}`))
	})

	// User Routes
	http.HandleFunc("/register", handlers.RegisterUserHandler)

	// Unified Order Routes (Handles both GET and POST)
	http.HandleFunc("/orders", handlers.OrdersHandler)

	port := ":8080"
	fmt.Printf("Backend API server listening on port %s 🚀\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}