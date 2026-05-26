package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ibrah/logistics-tracking-api/internal/models"
	"github.com/ibrah/logistics-tracking-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUserHandler handles POST requests for new user registration
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"` // customer, driver, admin
	}

	// 1. Decode request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Hash the user's password securely
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 3. Map to model
	user := models.User{
		FullName:     input.FullName,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Role:         input.Role,
	}

	// 4. Save to PostgreSQL
	id, err := repository.InsertUser(user)
	if err != nil {
		http.Error(w, "Database insertion failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// CORRECT
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"user_id": id,
	})
}