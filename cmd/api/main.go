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

	// 1. Health Check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "service": "delivery-api"}`))
	})

	// 2. User Routes
	http.HandleFunc("/register", handlers.RegisterUserHandler)

	// 3. Unified Order Routes (Handles both GET and POST)
	http.HandleFunc("/orders", handlers.OrdersHandler)

	// 4. Driver Assignment Route (Handles PUT)
	http.HandleFunc("/orders/assign", handlers.AssignDriverHandler)

	// 5. Order Completion Route (Handles PUT)
	http.HandleFunc("/orders/complete", handlers.CompleteOrderHandler)

	port := ":8080"
	fmt.Printf("Backend API server listening on port %s 🚀\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
