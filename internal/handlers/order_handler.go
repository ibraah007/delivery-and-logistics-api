package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ibrah/logistics-tracking-api/internal/models"
	"github.com/ibrah/logistics-tracking-api/internal/repository"
)

// OrdersHandler routes incoming traffic for /orders based on the HTTP Method (GET or POST)
func OrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createOrder(w, r)
	case http.MethodGet:
		getPendingOrders(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Private helper function to handle order creation (POST)
func createOrder(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CustomerID     uint    `json:"customer_id"`
		PickupAddress  string  `json:"pickup_address"`
		DropoffAddress string  `json:"dropoff_address"`
		CustomerPrice  float64 `json:"customer_price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.CustomerID == 0 || input.PickupAddress == "" || input.DropoffAddress == "" || input.CustomerPrice <= 0 {
		http.Error(w, "Missing required order fields or invalid price", http.StatusBadRequest)
		return
	}

	order := models.Order{
		CustomerID:     input.CustomerID,
		PickupAddress:  input.PickupAddress,
		DropoffAddress: input.DropoffAddress,
		CustomerPrice:  input.CustomerPrice,
	}

	id, err := repository.InsertOrder(order)
	if err != nil {
		http.Error(w, "Failed to book order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Delivery order booked successfully",
		"order_id": id,
		"status":   "pending",
	})
}

// Private helper function to handle order fetching (GET)
func getPendingOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := repository.GetPendingOrders()
	if err != nil {
		http.Error(w, "Failed to retrieve orders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

// AssignDriverHandler handles PUT requests to assign a driver to an order
func AssignDriverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		OrderID  uint `json:"order_id"`
		DriverID uint `json:"driver_id"`
	}

	// 1. Decode incoming assignment JSON
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Simple validation
	if input.OrderID == 0 || input.DriverID == 0 {
		http.Error(w, "Missing order_id or driver_id", http.StatusBadRequest)
		return
	}

	// 3. Update database state
	updated, err := repository.AssignDriver(input.OrderID, input.DriverID)
	if err != nil {
		http.Error(w, "Database error during assignment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !updated {
		http.Error(w, "Assignment failed: Order not found or already assigned", http.StatusNotFound)
		return
	}

	// 4. Return success status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Driver assigned successfully",
		"order_id": input.OrderID,
		"status":   "in_transit",
	})
}