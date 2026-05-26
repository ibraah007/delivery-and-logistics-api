package models

import "time"

type User struct {
	ID           uint      `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Hidden from JSON responses for security
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}