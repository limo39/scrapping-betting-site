# Betting Odds Scraper Makefile

.PHONY: build run test clean install deps docker-build docker-run

# Variables
BINARY_NAME=betting-odds-scraper
DOCKER_IMAGE=betting-odds-scraper
PORT=8080

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o bin/$(BINARY_NAME) main.go

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	go run main.go

# Run with live reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@echo "Starting development server with live reload..."
	air

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Quick test without Chrome noise
test-simple:
	@echo "Running simple scraper test..."
	go run cmd/test-simple/main.go

# Run in demo mode (fast, no Chrome required)
demo:
	@echo "Starting demo mode..."
	cp .env.demo .env
	go run main.go

# Test demo mode
test-demo:
	@echo "Testing demo mode..."
	LOG_LEVEL=demo go run cmd/test-simple/main.go

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	cp bin/$(BINARY_NAME) /usr/local/bin/

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run

# Security scan (requires gosec)
security:
	@echo "Running security scan..."
	gosec ./...

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p $(PORT):$(PORT) --rm $(DOCKER_IMAGE)

# Start all services
start: deps build
	@echo "Starting betting odds scraper..."
	./bin/$(BINARY_NAME)

# Quick setup for new developers
setup:
	@echo "Setting up development environment..."
	go mod tidy
	cp .env.example .env
	@echo "Setup complete! Run 'make run' to start the server."

# Scrape odds manually
scrape:
	@echo "Triggering manual scrape..."
	curl -X POST http://localhost:$(PORT)/api/v1/scrape/trigger

# Check service health
health:
	@echo "Checking service health..."
	curl http://localhost:$(PORT)/api/v1/health

# Get best odds
odds:
	@echo "Getting best odds..."
	curl http://localhost:$(PORT)/api/v1/odds/best | jq .

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  dev           - Run with live reload (requires air)"
	@echo "  deps          - Install dependencies"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  install       - Install binary to /usr/local/bin"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code (requires golangci-lint)"
	@echo "  security      - Run security scan (requires gosec)"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  start         - Setup and start the application"
	@echo "  setup         - Setup development environment"
	@echo "  scrape        - Trigger manual scrape"
	@echo "  health        - Check service health"
	@echo "  odds          - Get best odds (requires jq)"
	@echo "  help          - Show this help message"