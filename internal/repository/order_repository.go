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

// GetPendingOrders retrieves all orders from the database that are currently pending
func GetPendingOrders() ([]models.Order, error) {
	query := `
		SELECT id, customer_id, driver_id, status, pickup_address, dropoff_address, customer_price, company_expense, created_at, updated_at
		FROM orders
		WHERE status = 'pending'
		ORDER BY created_at DESC;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := configs.DBInstance.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order

	// Loop through rows returned by PostgreSQL
	for rows.Next() {
		var o models.Order
		// Using pointers or checking for NULL values on nullable fields like driver_id
		err := rows.Scan(
			&o.ID,
			&o.CustomerID,
			&o.DriverID,
			&o.Status,
			&o.PickupAddress,
			&o.DropoffAddress,
			&o.CustomerPrice,
			&o.CompanyExpense,
			&o.CreatedAt,
			&o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return empty slice instead of null if no orders found
	if orders == nil {
		orders = []models.Order{}
	}

	return orders, nil
}
// AssignDriver updates an order's status to 'in_transit' and sets the driver_id
func AssignDriver(orderID uint, driverID uint) (bool, error) {
	query := `
		UPDATE orders
		SET driver_id = $1, status = 'in_transit', updated_at = $2
		WHERE id = $3 AND status = 'pending';
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := configs.DBInstance.ExecContext(ctx, query, driverID, time.Now(), orderID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

// CompleteOrder marks an order as delivered and logs the final company expenses
func CompleteOrder(orderID uint, expense float64) (bool, error) {
	query := `
		UPDATE orders
		SET status = 'delivered', company_expense = $1, updated_at = $2
		WHERE id = $3 AND status = 'in_transit';
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := configs.DBInstance.ExecContext(ctx, query, expense, time.Now(), orderID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

// ProfitMargins represents the financial health summary of the logistics engine
type ProfitMargins struct {
	TotalRevenue  float64 `json:"total_revenue"`
	TotalExpenses float64 `json:"total_expenses"`
	NetProfit     float64 `json:"net_profit"`
	ProfitMargin  float64 `json:"profit_margin_percentage"`
}

// GetProfitMargins runs aggregation analytics across all orders
func GetProfitMargins() (ProfitMargins, error) {
	query := `
		SELECT 
			COALESCE(SUM(customer_price), 0.00) AS total_revenue,
			COALESCE(SUM(company_expense), 0.00) AS total_expenses,
			COALESCE(SUM(customer_price) - SUM(company_expense), 0.00) AS net_profit,
			CASE 
				WHEN SUM(customer_price) > 0 THEN 
					ROUND(((SUM(customer_price) - SUM(company_expense)) / SUM(customer_price)) * 100, 2)
				ELSE 0.00
			END AS profit_margin_percentage
		FROM orders;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var margins ProfitMargins
	err := configs.DBInstance.QueryRowContext(ctx, query).Scan(
		&margins.TotalRevenue,
		&margins.TotalExpenses,
		&margins.NetProfit,
		&margins.ProfitMargin,
	)

	if err != nil {
		return margins, err
	}

	return margins, nil
}
