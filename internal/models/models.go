package models

import "time"


type User struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"` // "admin", "driver", "customer"
}


type Order struct {
	ID            uint      `json:"id"`
	CustomerID    uint      `json:"customer_id"`
	DriverID      uint      `json:"driver_id"`
	Status        string    `json:"status"` // "pending", "in_transit", "delivered"
	PickupLocation  string  `json:"pickup_location"`
	DropoffLocation string  `json:"dropoff_location"`
	CreatedAt     time.Time `json:"created_at"`
}