package repository

import (
	"context"
	"time"

	"github.com/ibrah/logistics-tracking-api/configs"
	"github.com/ibrah/logistics-tracking-api/internal/models"
)

// InsertUser writes a new user record into the PostgreSQL database
func InsertUser(user models.User) (uint, error) {
	query := `
		INSERT INTO users (full_name, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	var lastInsertID uint
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := configs.DBInstance.QueryRowContext(ctx, query,
		user.FullName,
		user.Email,
		user.PasswordHash,
		user.Role,
		time.Now(),
	).Scan(&lastInsertID)

	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}