package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aungmyozaw92/go-api-setup/pkg/utils"
)

// AuthMiddleware creates a JWT authentication middleware
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeErrorResponse(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Check if it starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeErrorResponse(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			// Extract token
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				writeErrorResponse(w, "Token required", http.StatusUnauthorized)
				return
			}

			// Validate token
			claims, err := utils.ValidateJWT(token, jwtSecret)
			if err != nil {
				writeErrorResponse(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "user_email", claims.Email)

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// writeErrorResponse writes an error response in JSON format
func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := map[string]string{
		"error": message,
	}
	
	json.NewEncoder(w).Encode(response)
} 