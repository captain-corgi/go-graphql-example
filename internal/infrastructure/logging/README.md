# Logging Infrastructure

This package provides structured logging infrastructure for the GraphQL service using Go's standard `log/slog` package.

## Features

- **Structured Logging**: JSON and text format support
- **Request Tracing**: Request ID and correlation ID support
- **Layer-Specific Loggers**: Specialized loggers for each architecture layer
- **Configurable Levels**: Debug, Info, Warn, Error levels
- **Context Propagation**: Request context flows through all layers

## Usage

### Basic Logger Setup

```go
import "github.com/captain-corgi/go-graphql-example/internal/infrastructure/logging"

// Create logger factory
factory := logging.NewLoggerFactory(config.LoggingConfig{
    Level:  "info",
    Format: "json",
})

// Get base logger
logger := factory.GetLogger()
logger.Info("Application started")
```

### Layer-Specific Loggers

```go
// Domain layer
domainLogger := factory.GetDomainLogger()
domainLogger.LogEntityCreated(ctx, "User", "user-123")

// Application layer
appLogger := factory.GetApplicationLogger()
appLogger.LogUseCaseStarted(ctx, "CreateUser", map[string]interface{}{
    "email": "user@example.com",
})

// Infrastructure layer
infraLogger := factory.GetInfrastructureLogger()
infraLogger.LogDatabaseQuery(ctx, "SELECT * FROM users", time.Millisecond*50)

// Interface layer
interfaceLogger := factory.GetInterfaceLogger()
interfaceLogger.LogHTTPRequest(ctx, "POST", "/query", 200, time.Millisecond*100)
```

### Request Tracing

```go
// Add request ID to context
ctx = logging.WithRequestID(ctx, "req-12345")

// Add correlation ID to context
ctx = logging.WithCorrelationID(ctx, "corr-67890")

// Logger will automatically include these IDs
logger.WithRequestID(ctx).Info("Processing request")
```

### Custom Fields

```go
// Add custom fields to logger
logger.WithFields(map[string]interface{}{
    "user_id": "123",
    "action":  "create_post",
}).Info("User action performed")

// Chain multiple context additions
logger.WithRequestID(ctx).
    WithCorrelationID("corr-123").
    WithComponent("user-service").
    Info("Service operation completed")
```

## Configuration

The logging system is configured through the `LoggingConfig` struct:

```yaml
logging:
  level: "info"    # debug, info, warn, error
  format: "json"   # json, text
```

Environment variables can override configuration:

- `GRAPHQL_SERVICE_LOGGING_LEVEL`
- `GRAPHQL_SERVICE_LOGGING_FORMAT`

## Log Levels

- **Debug**: Detailed information for debugging
- **Info**: General information about application flow
- **Warn**: Warning messages for potentially harmful situations
- **Error**: Error messages for failures that don't stop the application

## Best Practices

1. **Use Layer-Specific Loggers**: Each layer has specialized logging methods
2. **Include Context**: Always pass context for request tracing
3. **Log at Appropriate Levels**: Use debug for detailed info, error for failures
4. **Include Relevant Fields**: Add structured data to make logs searchable
5. **Don't Log Sensitive Data**: Avoid logging passwords, tokens, or PII

## Request Flow Example

```
[INFO] interface: HTTP request received method=POST path=/query request_id=req-abc123
[INFO] application: Use case started use_case=CreateUser request_id=req-abc123
[DEBUG] infrastructure: Database query executed query="INSERT INTO users..." request_id=req-abc123
[INFO] domain: Domain entity created entity_type=User entity_id=user-456 request_id=req-abc123
[INFO] application: Use case completed use_case=CreateUser duration_ms=45 request_id=req-abc123
[INFO] interface: HTTP request processed status_code=200 duration_ms=50 request_id=req-abc123
```
