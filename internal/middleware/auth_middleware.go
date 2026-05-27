package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ibrah/logistics-tracking-api/internal/utils"
)

// AuthMiddleware protects routes from unauthorized access
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Grab the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		// 2. Expecting format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token format. Use 'Bearer <token>'", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		secretKey := []byte(os.Getenv("JWT_SECRET"))

		// 3. Parse and validate the token string
		claims := &utils.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 4. Attach user context information to the request so handlers can read it
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_role", claims.Role)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// AdminOnly protects routes to ensure only Admins can execute them
func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value("user_role").(string)
		if !ok || role != "admin" {
			http.Error(w, "Access denied: Admins only", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}
