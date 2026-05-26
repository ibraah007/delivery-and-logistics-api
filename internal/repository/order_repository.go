package repository

import (
	"context"
	"time"

	"github.com/ibrah/logistics-tracking-api/configs"
	"github.com/ibrah/logistics-tracking-api/internal/models"
)

// InsertOrder writes a new delivery order into the database
func InsertOrder(order models.Order) (uint, error) {
	query := `
		INSERT INTO orders (customer_id, status, pickup_address, dropoff_address, customer_price, company_expense, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;
	`

	var lastInsertID uint
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := configs.DBInstance.QueryRowContext(ctx, query,
		order.CustomerID,
		"pending", // New orders always start as pending
		order.PickupAddress,
		order.DropoffAddress,
		order.CustomerPrice,
		0.00, // Expense starts at zero until a driver is assigned
		time.Now(),
		time.Now(),
	).Scan(&lastInsertID)

	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}