package routes

import (
	"encoding/json"
	"net/http"

	"github.com/aungmyozaw92/go-api-setup/internal/handler"
	"github.com/aungmyozaw92/go-api-setup/internal/middleware"
	"github.com/gorilla/mux"
)

// SetupRoutes configures and returns the main router with all routes
func SetupRoutes(authHandler *handler.AuthHandler, userHandler *handler.UserHandler, jwtSecret string) *mux.Router {
	// Create main router
	router := mux.NewRouter()

	// Apply CORS middleware to all routes
	router.Use(middleware.CORSMiddleware)

	// Setup route groups
	setupPublicRoutes(router, authHandler)
	setupProtectedRoutes(router, userHandler, jwtSecret)
	setupHealthRoutes(router)

	// Setup versioned API routes (for future expansion)
	SetupV1Routes(router, authHandler, userHandler, jwtSecret)

	return router
}

// setupPublicRoutes configures routes that don't require authentication
func setupPublicRoutes(router *mux.Router, authHandler *handler.AuthHandler) {
	// Authentication routes (current/default version)
	auth := router.PathPrefix("/api/auth").Subrouter()
	auth.HandleFunc("/register", authHandler.Register).Methods("POST", "OPTIONS")
	auth.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")
}

// setupProtectedRoutes configures routes that require JWT authentication
func setupProtectedRoutes(router *mux.Router, userHandler *handler.UserHandler, jwtSecret string) {
	// Protected routes group
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware(jwtSecret))

	// User profile routes (current user)
	setupProfileRoutes(protected, userHandler)

	// User management routes (CRUD operations)
	setupUserManagementRoutes(protected, userHandler)
}

// setupProfileRoutes configures routes for current user profile management
func setupProfileRoutes(router *mux.Router, userHandler *handler.UserHandler) {
	router.HandleFunc("/profile", userHandler.GetProfile).Methods("GET", "OPTIONS")
	router.HandleFunc("/profile", userHandler.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/profile", userHandler.DeleteUser).Methods("DELETE", "OPTIONS")
}

// setupUserManagementRoutes configures routes for user CRUD operations
func setupUserManagementRoutes(router *mux.Router, userHandler *handler.UserHandler) {
	// User collection routes
	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET", "OPTIONS")

	// Individual user routes
	router.HandleFunc("/users/{id:[0-9]+}", userHandler.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/{id:[0-9]+}", userHandler.UpdateUserByID).Methods("PUT", "OPTIONS")
	router.HandleFunc("/users/{id:[0-9]+}", userHandler.DeleteUserByID).Methods("DELETE", "OPTIONS")
}

// setupHealthRoutes configures health check and utility routes
func setupHealthRoutes(router *mux.Router) {
	router.HandleFunc("/health", healthCheckHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/", rootHandler).Methods("GET", "OPTIONS")
}

// healthCheckHandler handles health check requests
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "healthy",
		"service": "go-api-setup",
		"version": "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// rootHandler handles requests to the root path
func rootHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Welcome to Go REST API",
		"version": "1.0.0",
		"endpoints": map[string]interface{}{
			"health":        "/health",
			"auth":          "/api/auth/*",
			"profile":       "/api/profile",
			"users":         "/api/users",
			"versioned_api": "/api/v1/*",
			"api_version":   "/api/v1/version",
			"documentation": "https://github.com/aungmyozaw92/go-api-setup",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
} 