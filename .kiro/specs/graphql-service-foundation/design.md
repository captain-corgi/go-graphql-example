# Design Document

## Overview

This design document outlines the implementation of a foundational GraphQL service using Go, following Clean Architecture principles with Domain-Driven Design. The service will be built using gqlgen for GraphQL code generation and Gin for HTTP routing, providing a robust, testable, and maintainable foundation that can be extended with additional features.

The implementation follows a schema-first approach where GraphQL schemas are defined first, then code is generated, and finally business logic is implemented. This ensures type safety and consistency between the API contract and implementation.

## Architecture

### Layer Structure

The application follows Clean Architecture with four distinct layers:

```
┌─────────────────────────────────────────────────────────────┐
│                    Interfaces Layer                        │
│  - GraphQL Resolvers (generated + custom)                  │
│  - HTTP Handlers (Gin routes)                              │
│  - Middleware (CORS, Auth, Logging)                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Application Layer                        │
│  - Use Cases / Application Services                        │
│  - Input/Output DTOs                                       │
│  - Application-specific business rules                     │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     Domain Layer                           │
│  - Entities (User, Post)                                   │
│  - Value Objects                                           │
│  - Repository Interfaces                                   │
│  - Domain Services                                         │
└─────────────────────────────────────────────────────────────┘
                              ▲
                              │
┌─────────────────────────────────────────────────────────────┐
│                 Infrastructure Layer                       │
│  - Database Implementations (PostgreSQL)                   │
│  - Configuration Management (Viper)                        │
│  - Logging (slog)                                          │
│  - External Service Clients                                │
└─────────────────────────────────────────────────────────────┘
```

### Dependency Flow

- **Interfaces** → **Application** → **Domain**
- **Infrastructure** → **Domain** (implements interfaces)
- **cmd/server** → All layers (composition root)

## Components and Interfaces

### GraphQL Schema Design

The service will implement a User management domain as an example, with the following GraphQL schema:

```graphql
# api/graphql/query.graphqls
type Query {
  user(id: ID!): User
  users(first: Int, after: String): UserConnection!
}

# api/graphql/mutation.graphqls  
type Mutation {
  createUser(input: CreateUserInput!): CreateUserPayload!
  updateUser(id: ID!, input: UpdateUserInput!): UpdateUserPayload!
  deleteUser(id: ID!): DeleteUserPayload!
}

# api/graphql/user.graphqls
type User {
  id: ID!
  email: String!
  name: String!
  createdAt: Time!
  updatedAt: Time!
}

type UserConnection {
  edges: [UserEdge!]!
  pageInfo: PageInfo!
}

type UserEdge {
  node: User!
  cursor: String!
}

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}

input CreateUserInput {
  email: String!
  name: String!
}

input UpdateUserInput {
  email: String
  name: String
}

type CreateUserPayload {
  user: User!
  errors: [Error!]
}

type UpdateUserPayload {
  user: User!
  errors: [Error!]
}

type DeleteUserPayload {
  success: Boolean!
  errors: [Error!]
}

type Error {
  message: String!
  field: String
  code: String
}

scalar Time
```

### Domain Layer Components

#### Entities

```go
// internal/domain/user/user.go
type User struct {
    id        UserID
    email     Email
    name      Name
    createdAt time.Time
    updatedAt time.Time
}

type UserID struct {
    value string
}

type Email struct {
    value string
}

type Name struct {
    value string
}
```

#### Repository Interface

```go
// internal/domain/user/repository.go
//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type Repository interface {
    FindByID(ctx context.Context, id UserID) (*User, error)
    FindAll(ctx context.Context, limit int, cursor string) ([]*User, error)
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id UserID) error
}
```

### Application Layer Components

#### Use Cases

```go
// internal/application/user/service.go
//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type Service interface {
    GetUser(ctx context.Context, req GetUserRequest) (*GetUserResponse, error)
    ListUsers(ctx context.Context, req ListUsersRequest) (*ListUsersResponse, error)
    CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error)
    UpdateUser(ctx context.Context, req UpdateUserRequest) (*UpdateUserResponse, error)
    DeleteUser(ctx context.Context, req DeleteUserRequest) (*DeleteUserResponse, error)
}

type service struct {
    userRepo domain.Repository
    logger   *slog.Logger
}
```

#### DTOs

```go
// internal/application/user/dto.go
type GetUserRequest struct {
    ID string
}

type GetUserResponse struct {
    User *UserDTO
}

type CreateUserRequest struct {
    Email string
    Name  string
}

type CreateUserResponse struct {
    User   *UserDTO
    Errors []ErrorDTO
}

type UserDTO struct {
    ID        string
    Email     string
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type ErrorDTO struct {
    Message string
    Field   string
    Code    string
}
```

### Infrastructure Layer Components

#### Database Repository

```go
// internal/infrastructure/persistence/sql/user_repository.go
type userRepository struct {
    db *sql.DB
}

func (r *userRepository) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
    query := `SELECT id, email, name, created_at, updated_at FROM users WHERE id = $1`
    // Implementation with prepared statements
}
```

#### Configuration

```go
// internal/infrastructure/config/config.go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Logging  LoggingConfig
}

type ServerConfig struct {
    Port         string
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
}

type DatabaseConfig struct {
    URL             string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
}
```

### Interfaces Layer Components

#### GraphQL Resolvers

```go
// internal/interfaces/graphql/resolver/resolver.go
type Resolver struct {
    userService application.Service
    logger      *slog.Logger
}

// internal/interfaces/graphql/resolver/user.resolvers.go (generated + custom)
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
    req := application.GetUserRequest{ID: id}
    resp, err := r.userService.GetUser(ctx, req)
    if err != nil {
        return nil, err
    }
    return mapUserDTOToGraphQL(resp.User), nil
}
```

#### HTTP Server Setup

```go
// internal/interfaces/http/server.go
func NewServer(cfg *config.Config, resolver *graphql.Resolver) *gin.Engine {
    r := gin.Default()
    
    // Middleware
    r.Use(middleware.RequestID())
    r.Use(middleware.Logger())
    r.Use(middleware.CORS())
    
    // GraphQL endpoints
    h := handler.New(generated.NewExecutableSchema(generated.Config{
        Resolvers: resolver,
    }))
    
    r.POST("/query", gin.WrapH(h))
    r.GET("/playground", gin.WrapH(playground.Handler("GraphQL", "/query")))
    
    return r
}
```

## Data Models

### Database Schema

```sql
-- migrations/001_create_users_table.sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);
```

### Domain Model Mapping

- **GraphQL Types** ↔ **Application DTOs** ↔ **Domain Entities** ↔ **Database Records**
- Each layer has its own representation optimized for its concerns
- Mapping functions handle conversion between layers

## Error Handling

### Error Types

```go
// internal/domain/errors/errors.go
type DomainError struct {
    Code    string
    Message string
    Field   string
}

func (e DomainError) Error() string {
    return e.Message
}

var (
    ErrUserNotFound     = DomainError{Code: "USER_NOT_FOUND", Message: "User not found"}
    ErrInvalidEmail     = DomainError{Code: "INVALID_EMAIL", Message: "Invalid email format", Field: "email"}
    ErrDuplicateEmail   = DomainError{Code: "DUPLICATE_EMAIL", Message: "Email already exists", Field: "email"}
)
```

### Error Handling Strategy

1. **Domain Layer**: Returns typed domain errors
2. **Application Layer**: Catches domain errors, logs them, converts to DTOs
3. **Interface Layer**: Maps application errors to GraphQL errors
4. **Infrastructure Layer**: Wraps external errors with context

### GraphQL Error Response

```go
// internal/interfaces/graphql/errors/handler.go
func HandleError(err error) *gqlerror.Error {
    var domainErr domain.DomainError
    if errors.As(err, &domainErr) {
        return &gqlerror.Error{
            Message: domainErr.Message,
            Extensions: map[string]interface{}{
                "code":  domainErr.Code,
                "field": domainErr.Field,
            },
        }
    }
    return gqlerror.Errorf("Internal server error")
}
```

## Testing Strategy

### Testing Pyramid

1. **Unit Tests**: Domain entities, value objects, use cases
2. **Integration Tests**: Repository implementations, HTTP handlers
3. **End-to-End Tests**: Full GraphQL operations

### Mock Generation

```go
// Each interface file includes:
//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks
```

### Test Structure

```go
// internal/application/user/service_test.go
func TestService_CreateUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockRepository(ctrl)
    service := NewService(mockRepo, slog.Default())
    
    tests := []struct {
        name    string
        request CreateUserRequest
        setup   func()
        want    *CreateUserResponse
        wantErr bool
    }{
        // Test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setup()
            got, err := service.CreateUser(context.Background(), tt.request)
            // Assertions
        })
    }
}
```

### Integration Test Example

```go
// internal/infrastructure/persistence/sql/user_repository_test.go
func TestUserRepository_Integration(t *testing.T) {
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := NewUserRepository(db)
    
    // Test repository operations against real database
}
```

## Configuration Management

### Configuration Structure

```yaml
# config/development.yaml
server:
  port: "8080"
  read_timeout: "30s"
  write_timeout: "30s"

database:
  url: "postgres://user:password@localhost:5432/graphql_service?sslmode=disable"
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: "5m"

logging:
  level: "debug"
  format: "json"
```

### Environment Variable Override

```go
// internal/infrastructure/config/loader.go
func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")
    viper.AddConfigPath(".")
    
    // Environment variable overrides
    viper.SetEnvPrefix("GRAPHQL_SERVICE")
    viper.AutomaticEnv()
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }
    
    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }
    
    return &cfg, nil
}
```

## Logging and Observability

### Structured Logging

```go
// internal/infrastructure/logging/logger.go
func NewLogger(cfg LoggingConfig) *slog.Logger {
    var handler slog.Handler
    
    opts := &slog.HandlerOptions{
        Level: parseLevel(cfg.Level),
    }
    
    if cfg.Format == "json" {
        handler = slog.NewJSONHandler(os.Stdout, opts)
    } else {
        handler = slog.NewTextHandler(os.Stdout, opts)
    }
    
    return slog.New(handler)
}
```

### Request Logging Middleware

```go
// internal/interfaces/http/middleware/logger.go
func Logger() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        return fmt.Sprintf(`{"time":"%s","method":"%s","path":"%s","status":%d,"latency":"%s","ip":"%s","user_agent":"%s"}%s`,
            param.TimeStamp.Format(time.RFC3339),
            param.Method,
            param.Path,
            param.StatusCode,
            param.Latency,
            param.ClientIP,
            param.Request.UserAgent(),
            "\n",
        )
    })
}
```

### Request ID Middleware

```go
// internal/interfaces/http/middleware/request_id.go
func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID()
        }
        
        c.Set("request_id", requestID)
        c.Header("X-Request-ID", requestID)
        
        ctx := context.WithValue(c.Request.Context(), "request_id", requestID)
        c.Request = c.Request.WithContext(ctx)
        
        c.Next()
    }
}
```

This design provides a solid foundation for a GraphQL service that can be extended with additional features while maintaining clean architecture principles and ensuring testability, maintainability, and scalability.
