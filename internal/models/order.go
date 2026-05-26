package models

import "time"

type Order struct {
	ID             uint      `json:"id"`
	CustomerID     uint      `json:"customer_id"`
	DriverID       *uint     `json:"driver_id,omitempty"` // Pointer because it can be null initially
	Status         string    `json:"status"`
	PickupAddress  string    `json:"pickup_address"`
	DropoffAddress string    `json:"dropoff_address"`
	CustomerPrice  float64   `json:"customer_price"`
	CompanyExpense float64   `json:"company_expense"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}