# Configuration Management

This package provides configuration management for the GraphQL service using Viper. It supports loading configuration from YAML files and environment variables with proper validation.

## Features

- **YAML Configuration Files**: Load configuration from `config.yaml` files
- **Environment Variable Overrides**: Override any configuration value using environment variables
- **Validation**: Comprehensive validation of all configuration values
- **Defaults**: Sensible defaults for all configuration options
- **Type Safety**: Strongly typed configuration structs

## Configuration Structure

The configuration is organized into three main sections:

### Server Configuration

- `port`: HTTP server port (default: "8080")
- `read_timeout`: HTTP read timeout (default: "30s")
- `write_timeout`: HTTP write timeout (default: "30s")
- `idle_timeout`: HTTP idle timeout (default: "120s")

### Database Configuration

- `url`: Database connection URL (default: postgres://user:password@localhost:5432/graphql_service?sslmode=disable)
- `max_open_conns`: Maximum number of open connections (default: 25)
- `max_idle_conns`: Maximum number of idle connections (default: 5)
- `conn_max_lifetime`: Maximum connection lifetime (default: "5m")
- `conn_max_idle_time`: Maximum connection idle time (default: "5m")

### Logging Configuration

- `level`: Log level - debug, info, warn, error (default: "info")
- `format`: Log format - json, text (default: "json")

## Usage

### Basic Usage

```go
package main

import (
    "log"
    "github.com/captain-corgi/go-graphql-example/internal/infrastructure/config"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Use configuration
    fmt.Printf("Server port: %s\n", cfg.Server.Port)
    fmt.Printf("Database URL: %s\n", cfg.Database.URL)
}
```

### Using MustLoad (panics on error)

```go
func main() {
    // Load configuration and panic if it fails
    cfg := config.MustLoad()
    
    // Use configuration
    startServer(cfg)
}
```

## Configuration Files

### File Locations

The configuration loader searches for `config.yaml` in the following locations:

1. `./configs/config.yaml`
2. `./config/config.yaml`
3. `./config.yaml`

### Example Configuration File

```yaml
server:
  port: "8080"
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "120s"

database:
  url: "postgres://user:password@localhost:5432/graphql_service?sslmode=disable"
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: "5m"
  conn_max_idle_time: "5m"

logging:
  level: "info"
  format: "json"
```

## Environment Variables

All configuration values can be overridden using environment variables with the prefix `GRAPHQL_SERVICE_`. Nested configuration keys use underscores.

### Examples

```bash
# Override server port
export GRAPHQL_SERVICE_SERVER_PORT="9090"

# Override database URL
export GRAPHQL_SERVICE_DATABASE_URL="postgres://localhost/mydb"

# Override log level
export GRAPHQL_SERVICE_LOGGING_LEVEL="debug"

# Override log format
export GRAPHQL_SERVICE_LOGGING_FORMAT="text"

# Override database connection settings
export GRAPHQL_SERVICE_DATABASE_MAX_OPEN_CONNS="50"
export GRAPHQL_SERVICE_DATABASE_MAX_IDLE_CONNS="10"
```

## Validation

The configuration system includes comprehensive validation:

- **Required Fields**: Ensures required fields are not empty
- **Positive Values**: Validates that numeric values are positive where required
- **Enum Values**: Validates that string values are from allowed sets
- **Logical Constraints**: Ensures logical relationships (e.g., max_idle_conns ≤ max_open_conns)

### Validation Errors

If validation fails, the loader returns detailed error messages:

```
config validation failed: server config validation failed: server port is required
config validation failed: logging config validation failed: invalid log level: invalid (must be one of: debug, info, warn, error)
```

## Testing

The package includes comprehensive tests covering:

- Configuration validation
- Environment variable overrides
- Default value handling
- Error conditions

Run tests with:

```bash
go test ./internal/infrastructure/config/...
```

## Environment-Specific Configuration

You can create environment-specific configuration files:

- `config.development.yaml`
- `config.production.yaml`
- `config.test.yaml`

Load specific configurations by setting the config name:

```go
// This would be done in the loader if needed
viper.SetConfigName("config.development")
```
