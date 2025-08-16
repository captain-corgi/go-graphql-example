# Technology Stack

## Core Technologies

- **Go**: Version 1.24.3 (as defined in go.mod)
- **GraphQL**: Schema-first approach using gqlgen for code generation
- **Web Framework**: Gin for HTTP routing and middleware
- **Database**: PostgreSQL with database/sql (no ORM)
- **Configuration**: Viper for environment-based config management
- **Logging**: Standard library slog with structured JSON logging
- **Migrations**: golang-migrate for database schema management
- **Containerization**: Docker with multi-stage builds and Docker Compose
- **Build Automation**: Comprehensive Makefile for development workflow

## Key Dependencies

```go
// Core GraphQL and web
github.com/99designs/gqlgen v0.17.78
github.com/gin-gonic/gin v1.10.1

// Database and migrations
github.com/lib/pq v1.10.9
github.com/golang-migrate/migrate/v4 v4.18.3

// Configuration and utilities
github.com/spf13/viper v1.20.1
github.com/google/uuid v1.6.0

// Testing and mocking
github.com/stretchr/testify v1.10.0
github.com/golang/mock v1.6.0
```

## Build and Development Commands

### Quick Start (Recommended)

```bash
# Setup development environment
make setup

# Start all services with Docker Compose
make docker-run

# Run database migrations
make migrate-up

# Start development server
make dev

# Open GraphQL Playground
make playground
```

### Code Generation

```bash
# Generate all code (GraphQL + mocks)
make generate

# Generate GraphQL resolvers and models only
make generate-graphql

# Generate mocks for interfaces only
make generate-mocks
```

### Running the Application

```bash
# Run with development config
make dev

# Run the server (production mode)
make run

# Build binary
make build

# Build for Linux (Docker)
make build-linux
```

### Docker Operations

```bash
# Start all services (development)
make docker-run

# Build and start services
make docker-run-build

# Stop all services
make docker-stop

# View logs
make docker-logs

# Clean Docker resources
make docker-clean

# Production deployment
make docker-prod-build
```

### Database Operations

```bash
# Run migrations up
make migrate-up

# Run migrations down
make migrate-down

# Create new migration
make migrate-create NAME=migration_name

# Force migration version
make migrate-force VERSION=1
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run integration tests only
make test-integration

# Run benchmarks
make benchmark
```

### Code Quality

```bash
# Run all checks (lint, vet, test)
make check

# Format code
make format

# Run linters
make lint

# Run go vet
make vet

# Run security checks
make security
```

### Dependencies

```bash
# Download and tidy dependencies
make deps

# Update all dependencies
make deps-update

# Create vendor directory
make deps-vendor
```

## Configuration

The application uses environment-based configuration with YAML files:

- Default: `configs/config.yaml`
- Development: `configs/config.development.yaml`
- Docker: `configs/config.docker.yaml`
- Production: `configs/config.production.yaml`
- Override via `CONFIG_FILE` environment variable

## Docker Deployment

### Development

```bash
# Start all services (includes database)
make docker-run

# View logs
make docker-logs

# Stop services
make docker-stop
```

### Production

```bash
# Build and start production services
make docker-prod-build

# Set required environment variables
export POSTGRES_PASSWORD=your_secure_password
export DATABASE_URL=postgres://user:pass@host:5432/dbname?sslmode=require
```

## GraphQL Development Workflow

1. Design schema in `api/graphql/*.graphqls` files
2. Run `make generate-graphql` to generate code
3. Implement business logic in application services
4. Keep resolvers thin - delegate to application layer only
5. Test via GraphQL Playground: `make playground`
