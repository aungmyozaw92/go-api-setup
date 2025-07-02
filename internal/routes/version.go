package routes

import (
	"encoding/json"
	"net/http"

	"github.com/aungmyozaw92/go-api-setup/internal/handler"
	"github.com/gorilla/mux"
)

// SetupV1Routes configures routes for API version 1
func SetupV1Routes(router *mux.Router, authHandler *handler.AuthHandler, userHandler *handler.UserHandler, jwtSecret string) {
	// V1 API routes
	v1 := router.PathPrefix("/api/v1").Subrouter()

	// Authentication routes
	auth := v1.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", authHandler.Register).Methods("POST", "OPTIONS")
	auth.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")

	// Version info endpoint
	v1.HandleFunc("/version", versionHandler).Methods("GET", "OPTIONS")

	// Future: V1 protected routes can be added here
	// setupV1ProtectedRoutes(v1, userHandler, jwtSecret)
}

// versionHandler returns API version information
func versionHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"api_version": "1.0.0",
		"service":     "go-api-setup",
		"status":      "active",
		"features": []string{
			"user_management",
			"jwt_authentication",
			"crud_operations",
		},
		"documentation": "https://github.com/aungmyozaw92/go-api-setup",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Future: setupV1ProtectedRoutes would configure protected routes for V1
// func setupV1ProtectedRoutes(router *mux.Router, userHandler *handler.UserHandler, jwtSecret string) {
//     protected := router.NewRoute().Subrouter()
//     protected.Use(middleware.AuthMiddleware(jwtSecret))
//     // V1 specific protected routes...
// } 