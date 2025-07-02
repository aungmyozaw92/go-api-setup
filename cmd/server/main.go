package main

import (
	"log"
	"net/http"

	"github.com/aungmyozaw92/go-api-setup/internal/config"
	"github.com/aungmyozaw92/go-api-setup/internal/handler"
	"github.com/aungmyozaw92/go-api-setup/internal/middleware"
	"github.com/aungmyozaw92/go-api-setup/internal/repository"
	"github.com/aungmyozaw92/go-api-setup/internal/usecase"
	"github.com/aungmyozaw92/go-api-setup/pkg/database"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	config := config.Load()
	log.Println("Configuration loaded successfully")

	// Connect to database
	db, err := database.NewMySQLConnection(&config.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize use cases
	userUsecase := usecase.NewUserUsecase(userRepo, config.JWT.SecretKey)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userUsecase)
	userHandler := handler.NewUserHandler(userUsecase)

	// Setup routes
	router := mux.NewRouter()

	// Apply CORS middleware to all routes
	router.Use(middleware.CORSMiddleware)

	// Public routes (no authentication required)
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST", "OPTIONS")

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	}).Methods("GET")

	// Protected routes (authentication required)
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware(config.JWT.SecretKey))

	// User profile routes
	protected.HandleFunc("/profile", userHandler.GetProfile).Methods("GET", "OPTIONS")
	protected.HandleFunc("/profile", userHandler.UpdateUser).Methods("PUT", "OPTIONS")

	// User CRUD routes
	protected.HandleFunc("/users", userHandler.CreateUser).Methods("POST", "OPTIONS")
	protected.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET", "OPTIONS")
	protected.HandleFunc("/users/{id:[0-9]+}", userHandler.GetUser).Methods("GET", "OPTIONS")
	protected.HandleFunc("/users/{id:[0-9]+}", userHandler.UpdateUserByID).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/users/{id:[0-9]+}", userHandler.DeleteUserByID).Methods("DELETE", "OPTIONS")

	log.Printf("Server starting on port %s", config.Server.Port)
	log.Printf("Available endpoints:")
	log.Printf("Authentication:")
	log.Printf("  POST /api/auth/register     - Register a new user")
	log.Printf("  POST /api/auth/login        - Login user")
	log.Printf("User Profile (Protected):")
	log.Printf("  GET    /api/profile         - Get current user profile")
	log.Printf("  PUT    /api/profile         - Update current user profile")
	log.Printf("  DELETE /api/profile         - Delete current user account")
	log.Printf("User Management (Protected):")
	log.Printf("  POST   /api/users           - Create a new user")
	log.Printf("  GET    /api/users           - Get all users (with pagination)")
	log.Printf("  GET    /api/users/{id}      - Get user by ID")
	log.Printf("  PUT    /api/users/{id}      - Update user by ID")
	log.Printf("  DELETE /api/users/{id}      - Delete user by ID")
	log.Printf("Other:")
	log.Printf("  GET    /health              - Health check")

	if err := http.ListenAndServe(":"+config.Server.Port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
} 