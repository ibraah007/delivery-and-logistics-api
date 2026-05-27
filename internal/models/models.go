package models

import "time"

// User represents anyone in the system (Admin, Driver, or Customer)
type User struct {
	ID        uint      `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	PasswordHash string `json:"-"` // Never send passwords over the network
	Role      string    `json:"role"` // "admin", "driver", "customer"
	CreatedAt time.Time `json:"created_at"`
}

// DriverProfile tracks availability and metrics for company expenses
type DriverProfile struct {
	UserID          uint    `json:"user_id"`
	VehicleType     string  `json:"vehicle_type"`    // "bike", "car", "truck"
	IsAvailable     bool    `json:"is_available"`
	CurrentVehicle  string  `json:"current_vehicle"` // License plate number
	FuelExpenseRate float64 `json:"fuel_expense_rate"`
}

// Order represents a single delivery request with tracking and financial metrics
type Order struct {
	ID             uint      `json:"id"`
	CustomerID     uint      `json:"customer_id"`
	DriverID       *uint     `json:"driver_id,omitempty"` // Pointer handles NULL state smoothly before assignment
	Status         string    `json:"status"`              // "pending", "in_transit", "delivered"
	PickupAddress  string    `json:"pickup_address"`
	DropoffAddress string    `json:"dropoff_address"`
	CustomerPrice  float64   `json:"customer_price"`
	CompanyExpense float64   `json:"company_expense"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Location represents the live coordinates stream from a driver
type Location struct {
	DriverID  uint      `json:"driver_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
}
