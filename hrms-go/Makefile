.PHONY: help install run build test clean migrate seed docker-up docker-down

# Default target
help:
	@echo "Available commands:"
	@echo "  make install       - Install Go dependencies"
	@echo "  make run          - Run the API server"
	@echo "  make run-scheduler - Run the background scheduler"
	@echo "  make build        - Build the application"
	@echo "  make test         - Run tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make migrate      - Run database migrations"
	@echo "  make seed         - Seed database with test data"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make docker-build - Build Docker image"
	@echo "  make lint         - Run linter"
	@echo "  make fmt          - Format code"

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run the API server
run:
	@echo "Starting API server..."
	go run cmd/api/main.go

# Run the scheduler
run-scheduler:
	@echo "Starting scheduler..."
	go run cmd/scheduler/main.go

# Build the application
build:
	@echo "Building application..."
	go build -o bin/api cmd/api/main.go
	go build -o bin/scheduler cmd/scheduler/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run database migrations
migrate:
	@echo "Running migrations..."
	go run cmd/migrate/main.go

# Seed database
seed:
	@echo "Seeding database..."
	go run cmd/migrate/main.go --seed

# Start Docker containers
docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d

# Stop Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker-compose build

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Development watch mode (requires air)
dev:
	@echo "Starting development server with hot reload..."
	air

# Generate mocks (requires mockgen)
mocks:
	@echo "Generating mocks..."
	go generate ./...
