package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/aungmyozaw92/go-api-setup/internal/config"
	"github.com/aungmyozaw92/go-api-setup/internal/handler"
	"github.com/aungmyozaw92/go-api-setup/internal/repository"
	"github.com/aungmyozaw92/go-api-setup/internal/routes"
	"github.com/aungmyozaw92/go-api-setup/internal/usecase"
	"github.com/aungmyozaw92/go-api-setup/pkg/database"
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

	// Start background user count monitoring
	go startUserCountMonitoring(userRepo)

	// Initialize use cases
	userUsecase := usecase.NewUserUsecase(userRepo, config.JWT.SecretKey)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userUsecase)
	userHandler := handler.NewUserHandler(userUsecase)

	// Setup routes using the routes package
	router := routes.SetupRoutes(authHandler, userHandler, config.JWT.SecretKey)

	// Log server information
	logServerInfo(config.Server.Port)

	// Start server
	if err := http.ListenAndServe(":"+config.Server.Port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// logServerInfo logs the server startup information and available endpoints
func logServerInfo(port string) {
	log.Printf("Server starting on port %s", port)
	log.Printf("🚀 Go REST API Server")
	log.Printf("📍 Available endpoints:")
	log.Printf("")
	log.Printf("🌐 General:")
	log.Printf("  GET    /                    - API welcome message")
	log.Printf("  GET    /health              - Health check")
	log.Printf("")
	log.Printf("🔐 Authentication (Public):")
	log.Printf("  POST   /api/auth/register   - Register a new user")
	log.Printf("  POST   /api/auth/login      - Login user")
	log.Printf("")
	log.Printf("👤 User Profile (Protected):")
	log.Printf("  GET    /api/profile         - Get current user profile")
	log.Printf("  PUT    /api/profile         - Update current user profile")
	log.Printf("  DELETE /api/profile         - Delete current user account")
	log.Printf("")
	log.Printf("👥 User Management (Protected):")
	log.Printf("  POST   /api/users           - Create a new user")
	log.Printf("  GET    /api/users           - Get all users (with pagination)")
	log.Printf("  GET    /api/users/{id}      - Get user by ID")
	log.Printf("  PUT    /api/users/{id}      - Update user by ID")
	log.Printf("  DELETE /api/users/{id}      - Delete user by ID")
	log.Printf("")
	log.Printf("📖 Documentation: https://github.com/aungmyozaw92/go-api-setup")
	log.Printf("🎯 Ready to accept requests!")
}

// startUserCountMonitoring runs a background goroutine that logs user count every 10 seconds
func startUserCountMonitoring(userRepo repository.UserRepository) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	log.Println("📊 Starting user count monitoring (every 10 seconds)")
	
	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			count, err := userRepo.Count(ctx)
			cancel()
			
			if err != nil {
				log.Printf("❌ Error getting user count: %v", err)
			} else {
				log.Printf("👥 Current user count: %d", count)
			}
		}
	}
} 