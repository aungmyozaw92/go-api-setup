.PHONY: test test-unit test-verbose test-coverage test-race clean build run docker-build docker-run docker-dev docker-stop docker-clean help

# Default target
all: test

# Run all tests
test:
	@echo "🧪 Running all tests..."
	go test ./...

# Run tests with verbose output
test-verbose:
	@echo "🧪 Running tests with verbose output..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

# Run tests with race detection
test-race:
	@echo "🏃 Running tests with race detection..."
	go test -race ./...

# Run specific package tests
test-utils:
	@echo "🧪 Testing utils package..."
	go test -v ./pkg/utils/...

test-usecase:
	@echo "🧪 Testing usecase package..."
	go test -v ./internal/usecase/...

test-handler:
	@echo "🧪 Testing handler package..."
	go test -v ./internal/handler/...

test-middleware:
	@echo "🧪 Testing middleware package..."
	go test -v ./internal/middleware/...

# Build the application
build:
	@echo "🔨 Building application..."
	go build -o bin/server ./cmd/server

# Run the application
run:
	@echo "🚀 Starting server..."
	go run ./cmd/server/main.go

# Docker Commands
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t go-api-setup:latest .

docker-run:
	@echo "🐳 Running with Docker Compose (Production)..."
	docker compose up -d

docker-dev:
	@echo "🐳 Running with Docker Compose (Development)..."
	docker compose -f docker-compose.dev.yml up

docker-stop:
	@echo "🛑 Stopping Docker containers..."
	docker compose down
	docker compose -f docker-compose.dev.yml down

docker-clean:
	@echo "🧹 Cleaning Docker containers and volumes..."
	docker compose down -v
	docker compose -f docker-compose.dev.yml down -v
	docker system prune -f

docker-logs:
	@echo "📋 Showing Docker logs..."
	docker compose logs -f

docker-logs-api:
	@echo "📋 Showing API container logs..."
	docker compose logs -f api

docker-shell:
	@echo "🐚 Opening shell in API container..."
	docker compose exec api sh

# Clean build artifacts and test files
clean:
	@echo "🧹 Cleaning up..."
	rm -f coverage.out coverage.html
	rm -rf bin/

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "🔍 Linting code..."
	golangci-lint run

# Generate mocks (if you add mockgen later)
generate-mocks:
	@echo "🔧 Generating mocks..."
	@echo "Mock generation would go here (install mockgen if needed)"

# Run benchmark tests
benchmark:
	@echo "⚡ Running benchmarks..."
	go test -bench=. ./...

# Check for vulnerabilities
security:
	@echo "🔒 Checking for security vulnerabilities..."
	govulncheck ./...

# Full CI pipeline
ci: deps fmt test-race test-coverage
	@echo "✅ CI pipeline completed successfully!"

# Development setup
dev-setup: deps
	@echo "⚙️ Setting up development environment..."
	@echo "Installing development tools..."
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "✅ Development setup complete!"

# Display help
help:
	@echo "Available commands:"
	@echo ""
	@echo "📋 Testing:"
	@echo "  test           - Run all tests"
	@echo "  test-verbose   - Run tests with verbose output"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-race      - Run tests with race detection"
	@echo "  test-utils     - Test utils package only"
	@echo "  test-usecase   - Test usecase package only"
	@echo "  test-handler   - Test handler package only"
	@echo "  test-middleware- Test middleware package only"
	@echo ""
	@echo "🔨 Building & Running:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application locally"
	@echo ""
	@echo "🐳 Docker Commands:"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run with Docker Compose (Production)"
	@echo "  docker-dev     - Run with Docker Compose (Development)"
	@echo "  docker-stop    - Stop Docker containers"
	@echo "  docker-clean   - Clean Docker containers and volumes"
	@echo "  docker-logs    - Show all container logs"
	@echo "  docker-logs-api- Show API container logs only"
	@echo "  docker-shell   - Open shell in API container"
	@echo ""
	@echo "🛠️ Development:"
	@echo "  deps           - Install dependencies"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"
	@echo "  clean          - Clean build artifacts"
	@echo "  dev-setup      - Setup development environment"
	@echo ""
	@echo "🔍 Quality & Security:"
	@echo "  benchmark      - Run benchmark tests"
	@echo "  security       - Check for vulnerabilities"
	@echo "  ci             - Run full CI pipeline"
	@echo ""
	@echo "❓ Help:"
	@echo "  help           - Show this help message" 