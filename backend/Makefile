.PHONY: help install build run dev clean test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

build: ## Build the server binary
	@echo "🔨 Building server..."
	go build -o bin/server ./cmd/server
	@echo "✅ Build complete: bin/server"

run: build ## Build and run the server
	@echo "🚀 Starting server..."
	./bin/server

dev: ## Run in development mode with hot reload (requires air)
	@echo "🔥 Starting in development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "⚠️  'air' not installed. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

test: ## Run tests
	@echo "🧪 Running tests..."
	go test -v ./...

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	rm -rf bin/
	@echo "✅ Clean complete"

format: ## Format code
	@echo "✨ Formatting code..."
	go fmt ./...

lint: ## Run linter
	@echo "🔍 Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Install from https://golangci-lint.run/usage/install/"; \
	fi

migrate-up: ## Run database migrations
	@echo "📊 Running migrations..."
	@echo "Migrations are auto-applied on server start"

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t inkstudio-backend .

docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	docker run -p 8080:8080 --env-file .env inkstudio-backend
