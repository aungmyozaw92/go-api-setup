.PHONY: test test-unit test-verbose test-coverage test-race clean build run

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
	@echo "  test           - Run all tests"
	@echo "  test-verbose   - Run tests with verbose output"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-race      - Run tests with race detection"
	@echo "  test-utils     - Test utils package only"
	@echo "  test-usecase   - Test usecase package only"
	@echo "  test-handler   - Test handler package only"
	@echo "  test-middleware- Test middleware package only"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Install dependencies"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"
	@echo "  benchmark      - Run benchmark tests"
	@echo "  security       - Check for vulnerabilities"
	@echo "  ci             - Run full CI pipeline"
	@echo "  dev-setup      - Setup development environment"
	@echo "  help           - Show this help message" 