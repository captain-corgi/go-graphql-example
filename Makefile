# Go GraphQL Example - Makefile
# Comprehensive build and development automation

# Variables
BINARY_NAME=server
BINARY_PATH=bin/$(BINARY_NAME)
MAIN_PATH=cmd/server/main.go
DOCKER_IMAGE=go-graphql-example
DOCKER_TAG=latest

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOGENERATE=$(GOCMD) generate

# Database variables
DB_URL=postgres://postgres:postgres@localhost:5432/graphql_service?sslmode=disable
MIGRATIONS_PATH=migrations

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help build clean test coverage run dev docker-build docker-run docker-stop generate migrate-up migrate-down migrate-create lint format deps check install

# Default target
all: clean deps generate test build

# Help target
help: ## Show this help message
	@echo "$(BLUE)Go GraphQL Example - Available Commands$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the application binary
	@echo "$(YELLOW)Building application...$(NC)"
	@mkdir -p bin
	$(GOBUILD) -o $(BINARY_PATH) -v $(MAIN_PATH)
	@echo "$(GREEN)Build completed: $(BINARY_PATH)$(NC)"

build-linux: ## Build for Linux (useful for Docker)
	@echo "$(YELLOW)Building for Linux...$(NC)"
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -a -installsuffix cgo -o $(BINARY_PATH)-linux $(MAIN_PATH)
	@echo "$(GREEN)Linux build completed: $(BINARY_PATH)-linux$(NC)"

# Clean targets
clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	$(GOCLEAN)
	rm -rf bin/
	rm -f coverage.out
	@echo "$(GREEN)Clean completed$(NC)"

# Development targets
run: ## Run the application locally
	@echo "$(YELLOW)Running application...$(NC)"
	$(GOCMD) run $(MAIN_PATH)

dev: ## Run with development configuration
	@echo "$(YELLOW)Running in development mode...$(NC)"
	CONFIG_FILE=configs/config.development.yaml $(GOCMD) run $(MAIN_PATH)

install: build ## Install the binary to GOPATH/bin
	@echo "$(YELLOW)Installing binary...$(NC)"
	cp $(BINARY_PATH) $(GOPATH)/bin/$(BINARY_NAME)
	@echo "$(GREEN)Binary installed to $(GOPATH)/bin/$(BINARY_NAME)$(NC)"

# Testing targets
test: ## Run all tests
	@echo "$(YELLOW)Running tests...$(NC)"
	$(GOTEST) -v ./...

test-short: ## Run tests with short flag
	@echo "$(YELLOW)Running short tests...$(NC)"
	$(GOTEST) -short -v ./...

test-integration: ## Run integration tests
	@echo "$(YELLOW)Running integration tests...$(NC)"
	$(GOTEST) -tags=integration -v ./...

coverage: ## Run tests with coverage
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

benchmark: ## Run benchmarks
	@echo "$(YELLOW)Running benchmarks...$(NC)"
	$(GOTEST) -bench=. -benchmem ./...

# Code generation targets
generate: ## Generate code (GraphQL, mocks, etc.)
	@echo "$(YELLOW)Generating code...$(NC)"
	$(GOGENERATE) ./...
	@echo "$(GREEN)Code generation completed$(NC)"

generate-graphql: ## Generate GraphQL resolvers and models
	@echo "$(YELLOW)Generating GraphQL code...$(NC)"
	$(GOCMD) run github.com/99designs/gqlgen generate
	@echo "$(GREEN)GraphQL code generation completed$(NC)"

generate-mocks: ## Generate mocks for interfaces
	@echo "$(YELLOW)Generating mocks...$(NC)"
	$(GOGENERATE) ./internal/domain/...
	$(GOGENERATE) ./internal/application/...
	@echo "$(GREEN)Mock generation completed$(NC)"

# Database targets
migrate-up: ## Run database migrations up
	@echo "$(YELLOW)Running migrations up...$(NC)"
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up
	@echo "$(GREEN)Migrations completed$(NC)"

migrate-down: ## Run database migrations down
	@echo "$(YELLOW)Running migrations down...$(NC)"
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down
	@echo "$(GREEN)Migrations rolled back$(NC)"

migrate-create: ## Create a new migration (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)Error: NAME is required. Usage: make migrate-create NAME=migration_name$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Creating migration: $(NAME)$(NC)"
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(NAME)
	@echo "$(GREEN)Migration created$(NC)"

migrate-force: ## Force migration version (usage: make migrate-force VERSION=1)
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Error: VERSION is required. Usage: make migrate-force VERSION=1$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Forcing migration to version $(VERSION)$(NC)"
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(VERSION)
	@echo "$(GREEN)Migration forced to version $(VERSION)$(NC)"

# Docker targets
docker-build: ## Build Docker image
	@echo "$(YELLOW)Building Docker image...$(NC)"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)$(NC)"

docker-run: ## Run application with Docker Compose
	@echo "$(YELLOW)Starting services with Docker Compose...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)Services started$(NC)"

docker-run-build: ## Build and run with Docker Compose
	@echo "$(YELLOW)Building and starting services...$(NC)"
	docker-compose up -d --build
	@echo "$(GREEN)Services built and started$(NC)"

docker-stop: ## Stop Docker Compose services
	@echo "$(YELLOW)Stopping Docker Compose services...$(NC)"
	docker-compose down
	@echo "$(GREEN)Services stopped$(NC)"

docker-logs: ## Show Docker Compose logs
	docker-compose logs -f

docker-clean: ## Clean Docker images and containers
	@echo "$(YELLOW)Cleaning Docker resources...$(NC)"
	docker-compose down -v --remove-orphans
	docker system prune -f
	@echo "$(GREEN)Docker cleanup completed$(NC)"

# Production Docker targets
docker-prod: ## Run production Docker Compose
	@echo "$(YELLOW)Starting production services...$(NC)"
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
	@echo "$(GREEN)Production services started$(NC)"

docker-prod-build: ## Build and run production Docker Compose
	@echo "$(YELLOW)Building and starting production services...$(NC)"
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
	@echo "$(GREEN)Production services built and started$(NC)"

# Code quality targets
lint: ## Run linters
	@echo "$(YELLOW)Running linters...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(RED)golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
		$(GOCMD) vet ./...; \
	fi
	@echo "$(GREEN)Linting completed$(NC)"

format: ## Format code
	@echo "$(YELLOW)Formatting code...$(NC)"
	$(GOCMD) fmt ./...
	@echo "$(GREEN)Code formatted$(NC)"

vet: ## Run go vet
	@echo "$(YELLOW)Running go vet...$(NC)"
	$(GOCMD) vet ./...
	@echo "$(GREEN)Vet completed$(NC)"

# Dependency targets
deps: ## Download and tidy dependencies
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "$(GREEN)Dependencies updated$(NC)"

deps-update: ## Update all dependencies
	@echo "$(YELLOW)Updating dependencies...$(NC)"
	$(GOGET) -u ./...
	$(GOMOD) tidy
	@echo "$(GREEN)Dependencies updated$(NC)"

deps-vendor: ## Create vendor directory
	@echo "$(YELLOW)Creating vendor directory...$(NC)"
	$(GOMOD) vendor
	@echo "$(GREEN)Vendor directory created$(NC)"

# Security targets
security: ## Run security checks
	@echo "$(YELLOW)Running security checks...$(NC)"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "$(RED)gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest$(NC)"; \
	fi
	@echo "$(GREEN)Security check completed$(NC)"

# Utility targets
check: lint vet test ## Run all checks (lint, vet, test)
	@echo "$(GREEN)All checks completed$(NC)"

ci: clean deps generate check build ## Run CI pipeline locally
	@echo "$(GREEN)CI pipeline completed$(NC)"

playground: ## Open GraphQL Playground
	@echo "$(YELLOW)Opening GraphQL Playground...$(NC)"
	@echo "$(BLUE)GraphQL Playground: http://localhost:8080/playground$(NC)"
	@if command -v open >/dev/null 2>&1; then \
		open http://localhost:8080/playground; \
	elif command -v xdg-open >/dev/null 2>&1; then \
		xdg-open http://localhost:8080/playground; \
	else \
		echo "$(YELLOW)Please open http://localhost:8080/playground in your browser$(NC)"; \
	fi

# Development setup
setup: ## Setup development environment
	@echo "$(YELLOW)Setting up development environment...$(NC)"
	@echo "$(BLUE)Installing development tools...$(NC)"
	$(GOGET) github.com/99designs/gqlgen@latest
	$(GOGET) github.com/golang/mock/mockgen@latest
	@if ! command -v migrate >/dev/null 2>&1; then \
		echo "$(BLUE)Installing migrate tool...$(NC)"; \
		$(GOGET) -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	fi
	@echo "$(GREEN)Development environment setup completed$(NC)"
	@echo "$(BLUE)Next steps:$(NC)"
	@echo "  1. Start database: make docker-run"
	@echo "  2. Run migrations: make migrate-up"
	@echo "  3. Start development server: make dev"
	@echo "  4. Open GraphQL Playground: make playground"