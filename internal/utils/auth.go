package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims defines the structured data we want to pack inside our digital entry pass
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a secure, signed JWT ticket for a logged-in user
func GenerateToken(userID uint, email string, role string) (string, error) {
	// 1. Grab our secret key from the environment variables
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	// 2. Pack the pass with user details and set it to expire in 24 hours
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 3. Create the token using the claims and sign it with our secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
