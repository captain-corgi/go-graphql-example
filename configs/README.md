# Configuration Files

This directory contains configuration files for different environments. The application uses Viper to load configuration with the following precedence:

1. Environment variables (highest priority)
2. Configuration file
3. Default values (lowest priority)

## Environment Files

- `config.yaml` - Default configuration
- `config.development.yaml` - Development environment settings
- `config.test.yaml` - Test environment settings  
- `config.staging.yaml` - Staging environment settings
- `config.production.yaml` - Production environment settings

## Environment Variable Overrides

All configuration values can be overridden using environment variables with the prefix `GRAPHQL_SERVICE_`. Nested values use underscores to separate levels.

Examples:

- `GRAPHQL_SERVICE_SERVER_PORT=9000` overrides `server.port`
- `GRAPHQL_SERVICE_DATABASE_URL=postgres://...` overrides `database.url`
- `GRAPHQL_SERVICE_LOGGING_LEVEL=debug` overrides `logging.level`

## Configuration Sections

### Server

- `port`: HTTP server port (default: "8080")
- `read_timeout`: Maximum duration for reading requests
- `write_timeout`: Maximum duration for writing responses  
- `idle_timeout`: Maximum duration for idle connections

### Database

- `url`: PostgreSQL connection string
- `max_open_conns`: Maximum number of open connections
- `max_idle_conns`: Maximum number of idle connections
- `conn_max_lifetime`: Maximum lifetime of connections
- `conn_max_idle_time`: Maximum idle time for connections

### Logging

- `level`: Log level (debug, info, warn, error)
- `format`: Log format (text, json)

## Usage

The application automatically loads the appropriate configuration file based on the environment. To specify a different environment, set the `GO_ENV` environment variable:

```bash
# Load development config
GO_ENV=development ./server

# Load production config  
GO_ENV=production ./server

# Load test config
GO_ENV=test ./server
```

If no environment is specified, the application loads `config.yaml` as the default.
