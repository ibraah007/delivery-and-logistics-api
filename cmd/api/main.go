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

	// 3. Unified Order Routes (Public or validated inside handler)
	http.HandleFunc("/orders", handlers.OrdersHandler)

	// 4. Driver Action Routes (Secured: Only Admins can assign drivers)
	http.HandleFunc("/orders/assign", middleware.AuthMiddleware(middleware.AdminOnly(handlers.AssignDriverHandler)))
	
	// 5. Order Completion Route (Secured: Must be logged in)
	http.HandleFunc("/orders/complete", middleware.AuthMiddleware(handlers.CompleteOrderHandler))

	// 6. Financial Analytics Route (Secured: Must be logged in AND an Admin)
	http.HandleFunc("/analytics/margins", middleware.AuthMiddleware(middleware.AdminOnly(handlers.GetAnalyticsMarginsHandler)))

	port := ":8080"
	fmt.Printf("Backend API server listening on port %s 🚀\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
