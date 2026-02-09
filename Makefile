.PHONY: help build run test clean docker-build docker-run docker-stop deploy-check lint

# Default target
help:
	@echo "Available targets:"
	@echo "  make build         - Build the Go application"
	@echo "  make run           - Run the application locally"
	@echo "  make test          - Run tests"
	@echo "  make lint          - Run linter"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run with Docker Compose"
	@echo "  make docker-stop   - Stop Docker Compose services"
	@echo "  make deploy-check  - Validate deployment configuration"

# Build the application
build:
	@echo "Building application..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./bin/server ./cmd/main.go
	@echo "✅ Build complete: ./bin/server"

# Run the application locally
run:
	@echo "Starting application..."
	go run ./cmd/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "✅ Tests complete"

# Run tests with coverage report
test-coverage: test
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "⚠️  golangci-lint not installed. Install with:"; \
		echo "    brew install golangci-lint"; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf ./bin
	rm -f coverage.out coverage.html
	@echo "✅ Clean complete"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t gold-go-backend:latest .
	@echo "✅ Docker image built: gold-go-backend:latest"

# Run with Docker Compose
docker-run:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d
	@echo "✅ Services started"
	@echo "Backend: http://localhost:8080"
	@echo "Redis: localhost:6379"
	@echo ""
	@echo "View logs: docker-compose logs -f"
	@echo "Stop services: make docker-stop"

# Stop Docker Compose services
docker-stop:
	@echo "Stopping Docker Compose services..."
	docker-compose down
	@echo "✅ Services stopped"

# Validate deployment configuration
deploy-check:
	@echo "Validating deployment configuration..."
	@echo ""
	@echo "Checking required files..."
	@test -f Dockerfile && echo "✅ Dockerfile exists" || echo "❌ Dockerfile missing"
	@test -f render.yaml && echo "✅ render.yaml exists" || echo "❌ render.yaml missing"
	@test -f .env.example && echo "✅ .env.example exists" || echo "❌ .env.example missing"
	@test -f docker-compose.yml && echo "✅ docker-compose.yml exists" || echo "❌ docker-compose.yml missing"
	@echo ""
	@echo "Checking .gitignore..."
	@grep -q "^.env$$" .gitignore && echo "✅ .env in .gitignore" || echo "⚠️  .env not in .gitignore"
	@echo ""
	@echo "Checking Go modules..."
	@go mod verify && echo "✅ Go modules verified" || echo "❌ Go modules invalid"
	@echo ""
	@echo "✅ Deployment configuration check complete"

# Install development dependencies
dev-setup:
	@echo "Installing development dependencies..."
	@echo "Installing golangci-lint..."
	@if ! command -v golangci-lint > /dev/null; then \
		brew install golangci-lint; \
	fi
	@echo "✅ Development setup complete"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted"

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy
	@echo "✅ Dependencies tidied"

# Full local test (before pushing)
pre-push: fmt tidy lint test deploy-check
	@echo ""
	@echo "✅ All pre-push checks passed!"
	@echo "Ready to commit and push to trigger deployment."
