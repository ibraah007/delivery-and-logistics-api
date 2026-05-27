package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ibrah/logistics-tracking-api/configs"
	"github.com/ibrah/logistics-tracking-api/internal/handlers"
	"github.com/ibrah/logistics-tracking-api/internal/middleware"
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

	// 2. User Authentication Routes
	http.HandleFunc("/register", handlers.RegisterUserHandler)
	http.HandleFunc("/login", handlers.LoginUserHandler)

	// 3. Unified Order Routes
	http.HandleFunc("/orders", handlers.OrdersHandler)

	// 4. Driver Action Routes
	http.HandleFunc("/orders/assign", middleware.AuthMiddleware(middleware.AdminOnly(handlers.AssignDriverHandler)))
	http.HandleFunc("/orders/complete", middleware.AuthMiddleware(handlers.CompleteOrderHandler))

	// 5. Live Telemetry Route (Secured: Must be logged in as a driver)
	http.HandleFunc("/driver/location", middleware.AuthMiddleware(handlers.TrackLocationHandler))

	// 6. Financial Analytics Route
	http.HandleFunc("/analytics/margins", middleware.AuthMiddleware(middleware.AdminOnly(handlers.GetAnalyticsMarginsHandler)))

	port := ":8080"
	fmt.Printf("Backend API server listening on port %s 🚀\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
