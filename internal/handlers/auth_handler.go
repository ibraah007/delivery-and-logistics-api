package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ibrah/logistics-tracking-api/internal/repository"
	"github.com/ibrah/logistics-tracking-api/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest holds the incoming credentials from the user
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse sends the secure token back to the user
type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

// LoginUserHandler handles authenticating users and issuing JWT tokens
func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		http.Error(w, "Invalid email or password request body", http.StatusBadRequest)
		return
	}

	// 1. Find the user in the database
	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// 2. Use Bcrypt to safely compare the stored hash against the incoming plain password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// 3. Create their secure digital entry ticket
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		http.Error(w, "Failed to generate security token", http.StatusInternalServerError)
		return
	}

	// 4. Return the token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{
		Token: token,
		Role:  user.Role,
	})
}
