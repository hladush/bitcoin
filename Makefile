# Bitcoin Tracker Makefile

.PHONY: build run test clean dev help

# Build the application
build:
	@echo "🔨 Building Bitcoin Tracker..."
	go build -o bitcoin-tracker cmd/server/main.go
	@echo "✅ Build complete!"

# Run the application
run: build
	@echo "🚀 Starting Bitcoin Tracker..."
	./bitcoin-tracker

# Run in development mode
dev:
	@echo "🛠️  Starting in development mode..."
	go run cmd/server/main.go

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Test the API (requires server to be running)
test-api: 
	@echo "🌐 Testing API endpoints..."
	./test_api.sh

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "🔍 Linting code..."
	golangci-lint run

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -f bitcoin-tracker
	rm -f bitcoin_tracker.db
	rm -f server.log
	rm -f nohup.out

# Show help
help:
	@echo "Bitcoin Tracker - Available commands:"
	@echo "  build     - Build the application"
	@echo "  run       - Build and run the application"
	@echo "  dev       - Run in development mode"
	@echo "  deps      - Install dependencies"
	@echo "  test      - Run unit tests"
	@echo "  test-api  - Test API endpoints (server must be running)"
	@echo "  fmt       - Format code"
	@echo "  lint      - Lint code"
	@echo "  clean     - Clean build artifacts"
	@echo "  help      - Show this help"

# Default target
.DEFAULT_GOAL := help
