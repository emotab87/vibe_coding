# RealWorld Conduit Backend Makefile
# Go 1.21+ required

.PHONY: help build run test clean dev deps lint fmt vet

# Variables
BINARY_NAME=conduit
BINARY_PATH=./cmd/main.go
BUILD_DIR=./build
GO_FILES=$(shell find . -type f -name '*.go')

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download and install dependencies
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

build: deps ## Build the application
	@echo "🔨 Building application..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(BINARY_PATH)

run: ## Run the application
	@echo "🚀 Running application..."
	go run $(BINARY_PATH)

dev: ## Run the application with hot reload (requires Air)
	@echo "🔥 Running with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

test: ## Run tests
	@echo "🧪 Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

bench: ## Run benchmarks
	@echo "⚡ Running benchmarks..."
	go test -bench=. -benchmem ./...

lint: ## Run linter (requires golangci-lint)
	@echo "🔍 Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Please install it: https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

fmt: ## Format Go code
	@echo "🎨 Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "🔍 Running go vet..."
	go vet ./...

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	go clean

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t realworld-conduit-backend .

docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	docker run -p 8080:8080 --env-file ../.env realworld-conduit-backend

migrate-up: ## Run database migrations (up)
	@echo "⬆️  Running migrations up..."
	@if [ -f "./cmd/migrate/main.go" ]; then \
		go run ./cmd/migrate/main.go up; \
	else \
		echo "Migration tool not found. Create ./cmd/migrate/main.go first"; \
	fi

migrate-down: ## Run database migrations (down)
	@echo "⬇️  Running migrations down..."
	@if [ -f "./cmd/migrate/main.go" ]; then \
		go run ./cmd/migrate/main.go down; \
	else \
		echo "Migration tool not found. Create ./cmd/migrate/main.go first"; \
	fi

install-tools: ## Install development tools
	@echo "🔧 Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

check: fmt vet lint test ## Run all checks (format, vet, lint, test)
	@echo "✅ All checks passed!"

# Development workflow
dev-setup: deps install-tools ## Setup development environment
	@echo "🎉 Development environment setup complete!"

.PHONY: all
all: clean deps check build ## Run clean, deps, check, and build