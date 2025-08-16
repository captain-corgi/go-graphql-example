# Project Structure and Organization

## Directory Layout

This project follows Clean Architecture with clear separation of concerns:

```
/
├── .github/                  # GitHub templates and workflows
├── api/graphql/              # GraphQL schema files (*.graphqls)
├── cmd/server/               # Application entry point (main.go)
├── internal/
│   ├── application/          # Use cases and application services
│   ├── domain/               # Core business logic and entities
│   ├── infrastructure/       # External adapters (DB, config, logging)
│   └── interfaces/           # Delivery layer (GraphQL resolvers, HTTP)
├── migrations/               # Database migration files
├── configs/                  # Configuration files by environment
├── docs/                     # Project documentation
├── examples/                 # API usage examples
├── scripts/                  # Utility scripts and database initialization
├── pkg/                      # Shared utilities (use sparingly)
├── Dockerfile                # Multi-stage Docker build
├── docker-compose.yml        # Development Docker Compose
├── docker-compose.override.yml  # Development overrides
├── docker-compose.prod.yml   # Production Docker Compose
├── Makefile                  # Build automation and development commands
├── .dockerignore             # Docker build exclusions
├── README.md                 # Project documentation
├── CONTRIBUTING.md           # Contribution guidelines
├── CODE_OF_CONDUCT.md        # Community guidelines
└── LICENSE                   # MIT License
```

## Layer Responsibilities

### Domain Layer (`internal/domain/`)

- **Pure business logic** - no external dependencies
- Entities, Value Objects, Aggregates
- Repository interfaces (ports)
- Domain services for complex business operations
- **Import rule**: Cannot import other internal packages

### Application Layer (`internal/application/`)

- Use cases and application services
- DTOs for data transfer
- Orchestrates domain operations
- **Import rule**: May import `domain` only

### Infrastructure Layer (`internal/infrastructure/`)

- Database implementations (`persistence/sql/`)
- Configuration loading (`config/`)
- Logging setup (`logging/`)
- External service clients
- **Import rule**: Implements `domain` interfaces, may import `domain` and selectively `application`

### Interfaces Layer (`internal/interfaces/`)

- GraphQL resolvers (`graphql/resolver/`)
- HTTP server and middleware (`http/`)
- Transport-specific concerns
- **Import rule**: May import `application` and `domain`

## File Organization Patterns

### GraphQL Schema (`api/graphql/`)

- Split by concern: `query.graphqls`, `mutation.graphqls`, `scalars.graphqls`
- Domain-specific files: `user.graphqls`, `post.graphqls`
- Keep schema files focused and cohesive

### Generated Code

- GraphQL generated code: `internal/interfaces/graphql/generated/`
- Models: `internal/interfaces/graphql/model/`
- **Never edit generated files directly**

### Testing Structure

- Test files alongside source: `*_test.go`
- Mocks in `mocks/` subdirectories with package name `mocks`
- Integration tests use build tags: `//go:build integration`

### Mock Generation

Each interface file should include:

```go
//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks
```

## Naming Conventions

### Packages

- Short, descriptive names: `user`, `config`, `persistence`
- Avoid stuttering: `user.User` not `user.UserEntity`

### Files

- Group related functionality: `user.go`, `user_test.go`
- Repository implementations: `user_repository.go`
- Service implementations: `user_service.go`

### Interfaces

- Named by capability: `UserRepository`, `UserService`
- Keep interfaces small and focused

## Import Guidelines

### Dependency Direction (Clean Architecture)

```
interfaces → application → domain
     ↓            ↑
infrastructure ----┘
```

### Composition Root

- `cmd/server/main.go` is the only place that imports all layers
- Handles dependency injection and application wiring
- Keep main.go focused on composition, not business logic

## Configuration Organization

### Environment Files

- `configs/config.yaml` - base configuration
- `configs/config.development.yaml` - local development
- `configs/config.docker.yaml` - Docker environment
- `configs/config.production.yaml` - production deployment
- Use `CONFIG_FILE` environment variable for custom paths

### Structure

- Group related settings: `server`, `database`, `logging`
- Provide sensible defaults
- Validate configuration at startup
- Support environment variable substitution in production

## Docker and Deployment

### Docker Files

- `Dockerfile` - Multi-stage build for production
- `docker-compose.yml` - Base services definition
- `docker-compose.override.yml` - Development overrides (auto-loaded)
- `docker-compose.prod.yml` - Production configuration
- `.dockerignore` - Exclude unnecessary files from build context

### Build Automation

- `Makefile` - Comprehensive development workflow automation
- Includes targets for building, testing, Docker operations, and database management
- Color-coded output for better developer experience
- Help target with command descriptions
