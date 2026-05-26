package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ibrah/logistics-tracking-api/internal/models"
	"github.com/ibrah/logistics-tracking-api/internal/repository"
)

// CreateOrderHandler handles POST requests for booking a new delivery
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		CustomerID     uint    `json:"customer_id"`
		PickupAddress  string  `json:"pickup_address"`
		DropoffAddress string  `json:"dropoff_address"`
		CustomerPrice  float64 `json:"customer_price"`
	}

	// 1. Decode incoming order JSON
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Quick validation
	if input.CustomerID == 0 || input.PickupAddress == "" || input.DropoffAddress == "" || input.CustomerPrice <= 0 {
		http.Error(w, "Missing required order fields or invalid price", http.StatusBadRequest)
		return
	}

	// 3. Map to our order structural model
	order := models.Order{
		CustomerID:     input.CustomerID,
		PickupAddress:  input.PickupAddress,
		DropoffAddress: input.DropoffAddress,
		CustomerPrice:  input.CustomerPrice,
	}

	// 4. Save order to PostgreSQL
	id, err := repository.InsertOrder(order)
	if err != nil {
		http.Error(w, "Failed to book order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Send back success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Delivery order booked successfully",
		"order_id": id,
		"status":   "pending",
	})
}