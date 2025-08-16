# Infrastructure Layer Integration Tests

This directory contains integration tests for the infrastructure layer components including database operations, configuration loading, and migration handling.

## Overview

The integration tests verify that infrastructure components work correctly with real external dependencies:

- **Configuration Integration Tests**: Test configuration loading from files and environment variables
- **Database Integration Tests**: Test database connections, migrations, and repository operations
- **User Repository Integration Tests**: Test complete CRUD operations against a real database

## Running Integration Tests

### Configuration Tests

Configuration integration tests don't require external dependencies and can be run directly:

```bash
# Run all configuration tests
go test -v ./internal/infrastructure/config

# Run only integration tests
go test -v ./internal/infrastructure/config -run TestConfigIntegration
```

### Database Tests

Database integration tests require a PostgreSQL database. Set the `TEST_DATABASE_URL` environment variable:

```bash
# Set test database URL
export TEST_DATABASE_URL="postgres://user:password@localhost:5432/test_db?sslmode=disable"

# Run database integration tests
go test -v ./internal/infrastructure/database -run TestDatabaseIntegration
go test -v ./internal/infrastructure/database -run TestMigrationIntegration

# Run user repository integration tests
go test -v ./internal/infrastructure/persistence/sql -run TestUserRepositoryIntegration
```

### Running All Integration Tests

```bash
# Set test database URL
export TEST_DATABASE_URL="postgres://user:password@localhost:5432/test_db?sslmode=disable"

# Run all infrastructure integration tests
go test -v ./internal/infrastructure/... -run Integration
```

## Test Database Setup

### Using Docker

The easiest way to set up a test database is using Docker:

```bash
# Start PostgreSQL container
docker run --name postgres-test \
  -e POSTGRES_USER=testuser \
  -e POSTGRES_PASSWORD=testpass \
  -e POSTGRES_DB=testdb \
  -p 5432:5432 \
  -d postgres:15

# Set test database URL
export TEST_DATABASE_URL="postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"

# Run tests
go test -v ./internal/infrastructure/... -run Integration

# Clean up
docker stop postgres-test
docker rm postgres-test
```

### Using Local PostgreSQL

If you have PostgreSQL installed locally:

```bash
# Create test database
createdb test_graphql_service

# Set test database URL
export TEST_DATABASE_URL="postgres://localhost:5432/test_graphql_service?sslmode=disable"

# Run tests
go test -v ./internal/infrastructure/... -run Integration

# Clean up
dropdb test_graphql_service
```

## Test Structure

### Configuration Integration Tests

Located in `internal/infrastructure/config/integration_test.go`:

- **TestLoadConfigFromFile**: Tests loading configuration from YAML files
- **TestLoadConfigWithEnvironmentOverrides**: Tests environment variable overrides
- **TestLoadConfigWithPartialFile**: Tests partial configuration files with defaults
- **TestLoadConfigWithInvalidFile**: Tests error handling for invalid YAML
- **TestLoadConfigWithValidationErrors**: Tests configuration validation
- **TestLoadConfigNoFile**: Tests loading with defaults only
- **TestLoadConfigComplexEnvironmentVariables**: Tests all environment variables

### Database Integration Tests

Located in `internal/infrastructure/database/integration_test.go`:

#### Database Connection Tests

- **TestDatabaseConnection**: Tests basic connection functionality
- **TestDatabaseConnectionWithInvalidURL**: Tests error handling for invalid URLs
- **TestDatabaseConnectionPoolConfiguration**: Tests connection pool settings
- **TestDatabaseHealthCheck**: Tests health check functionality
- **TestDatabaseBasicOperations**: Tests basic SQL operations
- **TestDatabaseTransactionHandling**: Tests transaction support

#### Migration Tests

- **TestRunMigrations**: Tests running database migrations
- **TestRunMigrationsNoChange**: Tests running migrations when no new migrations exist
- **TestRollbackMigration**: Tests rolling back migrations
- **TestGetMigrationVersion**: Tests getting current migration version
- **TestMigrationWithInvalidPath**: Tests error handling for invalid migration paths

### User Repository Integration Tests

Located in `internal/infrastructure/persistence/sql/user_repository_test.go`:

#### Basic CRUD Operations

- **TestCreateUser**: Tests user creation
- **TestCreateUserDuplicateEmail**: Tests duplicate email constraint
- **TestFindByID**: Tests finding users by ID
- **TestFindByIDNotFound**: Tests handling of non-existent users
- **TestFindByEmail**: Tests finding users by email
- **TestUpdateUser**: Tests user updates
- **TestDeleteUser**: Tests user deletion
- **TestExistsByEmail**: Tests checking user existence
- **TestFindAll**: Tests pagination
- **TestCount**: Tests user counting

#### Advanced Tests

- **TestConcurrentOperations**: Tests concurrent repository operations
- **TestRepositoryWithDatabaseFailure**: Tests behavior during database failures
- **TestRepositoryTransactionBehavior**: Tests transaction behavior
- **TestRepositoryPerformance**: Tests performance with larger datasets
- **TestRepositoryDataIntegrity**: Tests data integrity constraints
- **TestRepositoryEdgeCases**: Tests edge cases and boundary conditions

## Test Utilities

### Database Test Utilities

Located in `internal/infrastructure/database/testutils.go`:

- **TestDBManager**: Manages test database lifecycle
- **TestDBSetup**: Helper function for easy test database setup
- **TestDBSetupWithConfig**: Helper function with custom configuration
- **TruncateAllTables**: Utility to clean up test data
- **ExecuteSQL**: Utility to execute arbitrary SQL in tests

#### Usage Example

```go
func TestMyRepository(t *testing.T) {
    // Setup test database with migrations
    db, cleanup := database.TestDBSetup(t, "../../../../migrations")
    defer cleanup()
    
    // Use db for testing...
}
```

## Test Data Management

### Isolation

Each test suite uses its own database or cleans up data between tests to ensure isolation:

- Configuration tests use temporary directories
- Database tests create unique test databases
- Repository tests truncate tables between test cases

### Cleanup

All integration tests include proper cleanup:

- Temporary files and directories are removed
- Test databases are dropped
- Database connections are closed

## Performance Considerations

### Test Database

- Tests create isolated test databases to avoid conflicts
- Connection pools are configured with smaller limits for tests
- Tests include performance assertions to catch regressions

### Parallel Execution

- Configuration tests can run in parallel
- Database tests use separate databases to avoid conflicts
- Repository tests clean up data between test cases

## Troubleshooting

### Common Issues

1. **Database Connection Failures**
   - Ensure PostgreSQL is running
   - Check the `TEST_DATABASE_URL` environment variable
   - Verify database permissions

2. **Migration Failures**
   - Ensure migration files exist in the correct path
   - Check migration file syntax
   - Verify database schema permissions

3. **Test Timeouts**
   - Check database performance
   - Increase timeout values if needed
   - Verify network connectivity

### Debug Mode

Enable debug logging for more detailed test output:

```bash
# Run tests with verbose output
go test -v ./internal/infrastructure/... -run Integration

# Enable debug logging in tests (modify test code to use debug level)
```

## Requirements Coverage

These integration tests cover the following requirements from the specification:

- **Requirement 3.1**: Database connection establishment and configuration
- **Requirement 3.2**: Repository operations with prepared statements and transactions
- **Requirement 3.3**: Database migration handling and error management
- **Requirement 5.3**: Comprehensive testing of infrastructure components

The tests ensure that the infrastructure layer works correctly with real external dependencies and handles various scenarios including error conditions, edge cases, and performance requirements.
