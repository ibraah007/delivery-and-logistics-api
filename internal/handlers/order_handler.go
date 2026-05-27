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

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.OrderID == 0 || input.DriverID == 0 {
		http.Error(w, "Missing order_id or driver_id", http.StatusBadRequest)
		return
	}

	updated, err := repository.AssignDriver(input.OrderID, input.DriverID)
	if err != nil {
		http.Error(w, "Database error during assignment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !updated {
		http.Error(w, "Assignment failed: Order not found or already assigned", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Driver assigned successfully",
		"order_id": input.OrderID,
		"status":   "in_transit",
	})
}

// CompleteOrderHandler handles PUT requests to mark an order as delivered securely
func CompleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Extract the user_id from the secure JWT context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized: Identity missing from session", http.StatusUnauthorized)
		return
	}

	var input struct {
		OrderID        uint    `json:"order_id"`
		CompanyExpense float64 `json:"company_expense"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.OrderID == 0 || input.CompanyExpense < 0 {
		http.Error(w, "Missing order_id or invalid expense value", http.StatusBadRequest)
		return
	}

	// 2. Real-App Security Validation: Check who this order belongs to
	assignedDriverID, err := repository.GetOrderDriverID(input.OrderID)
	if err != nil {
		http.Error(w, "Order not found or database mismatch", http.StatusNotFound)
		return
	}

	// 3. Enforce order ownership rules
	if assignedDriverID != userID {
		http.Error(w, "Forbidden: You cannot complete an order assigned to a different driver", http.StatusForbidden)
		return
	}

	// 4. Update state machine in PostgreSQL if validation passes
	completed, err := repository.CompleteOrder(input.OrderID, input.CompanyExpense)
	if err != nil {
		http.Error(w, "Database error during order completion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !completed {
		http.Error(w, "Completion failed: Order is not currently marked as in_transit", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":         "Order marked as delivered successfully",
		"order_id":        input.OrderID,
		"status":          "delivered",
		"company_expense": input.CompanyExpense,
	})
}

// GetAnalyticsMarginsHandler handles GET requests for platform financial overview
func GetAnalyticsMarginsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	margins, err := repository.GetProfitMargins()
	if err != nil {
		http.Error(w, "Failed to calculate profit margins: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(margins)
}
