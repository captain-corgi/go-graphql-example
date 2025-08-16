# Database Infrastructure

This package provides database infrastructure components for the GraphQL service, including connection management, transaction handling, and database migrations.

## Components

### Connection Management (`connection.go`)

The `DB` struct wraps the standard `sql.DB` with additional functionality:

- **Connection Pooling**: Configurable connection pool settings
- **Health Checks**: Built-in health check functionality
- **Logging**: Structured logging for database operations
- **URL Masking**: Secure logging of database URLs

```go
db, err := database.NewConnection(cfg.Database, logger)
if err != nil {
    return fmt.Errorf("failed to create database connection: %w", err)
}
defer db.Close()
```

### Transaction Management (`transaction.go`)

Provides utilities for handling database transactions:

- **Automatic Rollback**: Transactions are automatically rolled back on errors or panics
- **Logging**: Transaction operations are logged with duration and status
- **Options Support**: Custom transaction options (isolation levels, read-only, etc.)

```go
err := db.WithTransaction(ctx, func(tx *sql.Tx) error {
    // Your transactional operations here
    return nil
})
```

### Migration Management (`migration.go`)

Handles database schema migrations using golang-migrate:

- **Automatic Migration**: Run all pending migrations on startup
- **Rollback Support**: Rollback migrations when needed
- **Version Tracking**: Track current migration version
- **Dirty State Detection**: Detect and handle dirty migration states

```go
migrationManager := database.NewMigrationManager(db, logger)
err := migrationManager.RunMigrations(ctx, "./migrations")
```

### Database Manager (`database.go`)

Provides a unified interface for all database components:

```go
manager, err := database.NewManager(cfg.Database, logger)
if err != nil {
    return err
}

// Initialize with migrations
err = manager.Initialize(ctx, "./migrations")
if err != nil {
    return err
}

defer manager.Close()
```

## Configuration

Database configuration is handled through the `config.DatabaseConfig` struct:

```yaml
database:
  url: "postgres://user:password@localhost:5432/dbname?sslmode=disable"
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: "5m"
  conn_max_idle_time: "1m"
```

## Migrations

Migration files should be placed in the `migrations/` directory with the following naming convention:

- `{version}_{description}.up.sql` - Forward migration
- `{version}_{description}.down.sql` - Rollback migration

Example:

- `001_create_users_table.up.sql`
- `001_create_users_table.down.sql`

## Testing

### Unit Tests

Unit tests are provided for core functionality and can be run without a database:

```bash
go test ./internal/infrastructure/database/...
```

### Integration Tests

Integration tests require a test database and are automatically skipped if `TEST_DATABASE_URL` is not set:

```bash
export TEST_DATABASE_URL="postgres://user:password@localhost:5432/test_db?sslmode=disable"
go test ./internal/infrastructure/persistence/sql/...
```

## Usage Example

```go
package main

import (
    "context"
    "log/slog"
    
    "github.com/captain-corgi/go-graphql-example/internal/infrastructure/config"
    "github.com/captain-corgi/go-graphql-example/internal/infrastructure/database"
    "github.com/captain-corgi/go-graphql-example/internal/infrastructure/persistence/sql"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        panic(err)
    }
    
    // Create logger
    logger := slog.Default()
    
    // Create database manager
    dbManager, err := database.NewManager(cfg.Database, logger)
    if err != nil {
        panic(err)
    }
    defer dbManager.Close()
    
    // Initialize database (run migrations)
    ctx := context.Background()
    err = dbManager.Initialize(ctx, "./migrations")
    if err != nil {
        panic(err)
    }
    
    // Create repository
    userRepo := sql.NewUserRepository(dbManager.DB, logger)
    
    // Use repository...
}
```

## Error Handling

The database infrastructure provides comprehensive error handling:

- **Connection Errors**: Detailed error messages for connection failures
- **Transaction Errors**: Automatic rollback with error context
- **Migration Errors**: Clear error messages for migration failures
- **Constraint Violations**: Domain-specific error mapping for database constraints

## Security Considerations

- Database URLs are masked in logs to prevent credential exposure
- Prepared statements are used to prevent SQL injection
- Connection pooling limits are enforced to prevent resource exhaustion
- Transaction timeouts prevent long-running operations from blocking resources
