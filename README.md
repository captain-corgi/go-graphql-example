# Go GraphQL Example

[![Go Version](https://img.shields.io/badge/Go-1.24.3-blue.svg)](https://golang.org)
[![GraphQL](https://img.shields.io/badge/GraphQL-gqlgen-e10098.svg)](https://gqlgen.com)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](Dockerfile)

A production-ready GraphQL service built in Go, demonstrating Clean Architecture principles with Domain-Driven Design (DDD) patterns. This project serves as both a functional GraphQL API and an educational reference for building scalable web services.

## 🚀 Features

- **GraphQL API**: Schema-first GraphQL implementation with gqlgen
- **Clean Architecture**: Clear separation of concerns across domain, application, infrastructure, and interface layers
- **User Management**: Complete CRUD operations with pagination support
- **Database Integration**: PostgreSQL with migrations and connection pooling
- **Docker Support**: Multi-stage builds with development and production configurations
- **Configuration Management**: Environment-based configuration with validation
- **Structured Logging**: JSON-formatted logging with configurable levels
- **Health Checks**: Built-in health monitoring and startup validation
- **Testing**: Comprehensive unit and integration tests with mocking
- **Development Tools**: Makefile automation and Docker Compose setup

## 📋 Table of Contents

- [Quick Start](#-quick-start)
- [Architecture](#-architecture)
- [API Documentation](#-api-documentation)
- [Development](#-development)
- [Docker Deployment](#-docker-deployment)
- [Testing](#-testing)
- [Configuration](#-configuration)
- [Contributing](#-contributing)
- [License](#-license)

## 🏃 Quick Start

### Prerequisites

- Go 1.24.3 or later
- PostgreSQL 12+ (or use Docker Compose)
- Make (optional, for convenience commands)

### Using Docker Compose (Recommended)

1. **Clone the repository**

   ```bash
   git clone https://github.com/captain-corgi/go-graphql-example.git
   cd go-graphql-example
   ```

2. **Start all services**

   ```bash
   make docker-run
   # or manually: docker-compose up -d
   ```

3. **Access the GraphQL Playground**

   ```bash
   make playground
   # or visit: http://localhost:8080/playground
   ```

### Local Development Setup

1. **Install dependencies**

   ```bash
   make setup
   ```

2. **Start PostgreSQL** (if not using Docker)

   ```bash
   # Using Docker for database only
   docker run --name postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=graphql_service -p 5432:5432 -d postgres:16-alpine
   ```

3. **Run migrations**

   ```bash
   make migrate-up
   ```

4. **Start the development server**

   ```bash
   make dev
   ```

## 🏗 Architecture

This project implements Clean Architecture with the following layers:

```
┌─────────────────────────────────────────────────────────────┐
│                    Interfaces Layer                         │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │  GraphQL        │  │  HTTP Server    │                  │
│  │  Resolvers      │  │  & Middleware   │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                   Application Layer                         │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │  Use Cases      │  │  DTOs &         │                  │
│  │  & Services     │  │  Orchestration  │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                     Domain Layer                            │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │  Entities &     │  │  Repository     │                  │
│  │  Value Objects  │  │  Interfaces     │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                 Infrastructure Layer                        │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │  Database       │  │  Configuration  │                  │
│  │  Repositories   │  │  & Logging      │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
```

### Key Principles

- **Dependency Inversion**: Inner layers define interfaces, outer layers implement them
- **Framework Independence**: Business logic is isolated from frameworks
- **Testability**: Easy mocking and testing through dependency injection
- **Single Responsibility**: Each layer has a clear, focused purpose

## 📚 API Documentation

### GraphQL Endpoint

- **URL**: `http://localhost:8080/query`
- **Playground**: `http://localhost:8080/playground`
- **Health Check**: `http://localhost:8080/health`

### Example Queries

**Get a user by ID:**

```graphql
query GetUser {
  user(id: "550e8400-e29b-41d4-a716-446655440001") {
    id
    email
    name
    createdAt
    updatedAt
  }
}
```

**List users with pagination:**

```graphql
query GetUsers {
  users(first: 10) {
    edges {
      node {
        id
        email
        name
      }
      cursor
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}
```

**Create a new user:**

```graphql
mutation CreateUser {
  createUser(input: {
    email: "user@example.com"
    name: "John Doe"
  }) {
    user {
      id
      email
      name
    }
    errors {
      message
      field
    }
  }
}
```

For more examples, see the [`examples/graphql/`](examples/graphql/) directory.

## 🛠 Development

### Available Make Commands

```bash
# Development
make dev              # Run with development config
make run              # Run the application
make build            # Build binary
make clean            # Clean build artifacts

# Code Generation
make generate         # Generate all code (GraphQL + mocks)
make generate-graphql # Generate GraphQL code only
make generate-mocks   # Generate mocks only

# Testing
make test             # Run all tests
make test-integration # Run integration tests
make coverage         # Generate coverage report

# Database
make migrate-up       # Run migrations
make migrate-down     # Rollback migrations
make migrate-create NAME=migration_name  # Create new migration

# Docker
make docker-run       # Start with Docker Compose
make docker-stop      # Stop Docker services
make docker-clean     # Clean Docker resources

# Code Quality
make lint             # Run linters
make format           # Format code
make check            # Run all checks

# Utilities
make setup            # Setup development environment
make playground       # Open GraphQL Playground
```

### Project Structure

```
├── api/graphql/              # GraphQL schema files
├── cmd/server/               # Application entry point
├── internal/
│   ├── application/          # Use cases and services
│   ├── domain/               # Core business logic
│   ├── infrastructure/       # External adapters
│   └── interfaces/           # Delivery layer
├── migrations/               # Database migrations
├── configs/                  # Configuration files
├── examples/                 # API usage examples
├── docs/                     # Documentation
└── scripts/                  # Utility scripts
```

### Adding New Features

1. **Define GraphQL Schema**: Add types to `api/graphql/*.graphqls`
2. **Generate Code**: Run `make generate-graphql`
3. **Implement Domain Logic**: Add entities and repositories in `internal/domain/`
4. **Create Use Cases**: Add services in `internal/application/`
5. **Implement Resolvers**: Update resolvers in `internal/interfaces/graphql/resolver/`
6. **Add Tests**: Write unit and integration tests
7. **Update Documentation**: Update relevant docs

## 🐳 Docker Deployment

### Development

```bash
# Start all services (includes database)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Production

```bash
# Build and start production services
make docker-prod-build

# Or manually
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
```

### Environment Variables for Production

```bash
# Required for production
export POSTGRES_PASSWORD=your_secure_password
export DATABASE_URL=postgres://user:pass@host:5432/dbname?sslmode=require

# Optional
export POSTGRES_DB=graphql_service
export POSTGRES_USER=postgres
```

## 🧪 Testing

### Running Tests

```bash
# All tests
make test

# With coverage
make coverage

# Integration tests only
make test-integration

# Specific package
go test ./internal/domain/user/...
```

### Test Structure

- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test component interactions with real database
- **Mocks**: Generated mocks for all interfaces using gomock

### Writing Tests

```go
//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// Add this directive to interface files for automatic mock generation
```

## ⚙️ Configuration

Configuration is managed through YAML files with environment variable overrides:

### Configuration Files

- `configs/config.yaml` - Base configuration
- `configs/config.development.yaml` - Development overrides
- `configs/config.docker.yaml` - Docker environment
- `configs/config.production.yaml` - Production settings

### Environment Variables

```bash
# Application
CONFIG_FILE=path/to/config.yaml
GIN_MODE=release

# Database
DATABASE_URL=postgres://user:pass@host:5432/dbname

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes following our [coding standards](docs/conventions.md)
4. Add tests for new functionality
5. Run the test suite: `make check`
6. Commit your changes: `git commit -m 'Add amazing feature'`
7. Push to the branch: `git push origin feature/amazing-feature`
8. Open a Pull Request

### Code Standards

- Follow Go conventions and idioms
- Maintain Clean Architecture boundaries
- Write comprehensive tests
- Update documentation for new features
- Use conventional commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [gqlgen](https://gqlgen.com/) for GraphQL code generation
- [Gin](https://gin-gonic.com/) for the HTTP framework
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) by Robert C. Martin
- [Domain-Driven Design](https://domainlanguage.com/ddd/) by Eric Evans

## 📞 Support

- 📖 [Documentation](docs/)
- 🐛 [Issue Tracker](https://github.com/captain-corgi/go-graphql-example/issues)
- 💬 [Discussions](https://github.com/captain-corgi/go-graphql-example/discussions)

---

**Happy coding!** 🚀
