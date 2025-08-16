package logging

import (
	"context"
	"log/slog"
	"time"
)

// DomainLogger provides logging utilities for the domain layer
type DomainLogger struct {
	*Logger
}

// NewDomainLogger creates a logger for domain layer
func NewDomainLogger(logger *Logger) *DomainLogger {
	return &DomainLogger{
		Logger: logger.WithComponent("domain"),
	}
}

// LogEntityCreated logs when a domain entity is created
func (l *DomainLogger) LogEntityCreated(ctx context.Context, entityType, entityID string) {
	l.WithRequestID(ctx).Info("Domain entity created",
		"entity_type", entityType,
		"entity_id", entityID,
	)
}

// LogEntityUpdated logs when a domain entity is updated
func (l *DomainLogger) LogEntityUpdated(ctx context.Context, entityType, entityID string) {
	l.WithRequestID(ctx).Info("Domain entity updated",
		"entity_type", entityType,
		"entity_id", entityID,
	)
}

// LogEntityDeleted logs when a domain entity is deleted
func (l *DomainLogger) LogEntityDeleted(ctx context.Context, entityType, entityID string) {
	l.WithRequestID(ctx).Info("Domain entity deleted",
		"entity_type", entityType,
		"entity_id", entityID,
	)
}

// LogValidationError logs domain validation errors
func (l *DomainLogger) LogValidationError(ctx context.Context, entityType string, err error) {
	l.WithRequestID(ctx).WithError(err).Warn("Domain validation failed",
		"entity_type", entityType,
	)
}

// ApplicationLogger provides logging utilities for the application layer
type ApplicationLogger struct {
	*Logger
}

// NewApplicationLogger creates a logger for application layer
func NewApplicationLogger(logger *Logger) *ApplicationLogger {
	return &ApplicationLogger{
		Logger: logger.WithComponent("application"),
	}
}

// LogUseCaseStarted logs when a use case starts
func (l *ApplicationLogger) LogUseCaseStarted(ctx context.Context, useCase string, params map[string]interface{}) {
	l.WithRequestID(ctx).WithFields(params).Info("Use case started",
		"use_case", useCase,
	)
}

// LogUseCaseCompleted logs when a use case completes successfully
func (l *ApplicationLogger) LogUseCaseCompleted(ctx context.Context, useCase string, duration time.Duration) {
	l.WithRequestID(ctx).Info("Use case completed",
		"use_case", useCase,
		"duration_ms", duration.Milliseconds(),
	)
}

// LogUseCaseFailed logs when a use case fails
func (l *ApplicationLogger) LogUseCaseFailed(ctx context.Context, useCase string, err error, duration time.Duration) {
	l.WithRequestID(ctx).WithError(err).Error("Use case failed",
		"use_case", useCase,
		"duration_ms", duration.Milliseconds(),
	)
}

// InfrastructureLogger provides logging utilities for the infrastructure layer
type InfrastructureLogger struct {
	*Logger
}

// NewInfrastructureLogger creates a logger for infrastructure layer
func NewInfrastructureLogger(logger *Logger) *InfrastructureLogger {
	return &InfrastructureLogger{
		Logger: logger.WithComponent("infrastructure"),
	}
}

// LogDatabaseQuery logs database queries
func (l *InfrastructureLogger) LogDatabaseQuery(ctx context.Context, query string, duration time.Duration) {
	l.WithRequestID(ctx).Debug("Database query executed",
		"query", query,
		"duration_ms", duration.Milliseconds(),
	)
}

// LogDatabaseError logs database errors
func (l *InfrastructureLogger) LogDatabaseError(ctx context.Context, operation string, err error) {
	l.WithRequestID(ctx).WithError(err).Error("Database operation failed",
		"operation", operation,
	)
}

// LogExternalServiceCall logs external service calls
func (l *InfrastructureLogger) LogExternalServiceCall(ctx context.Context, service, endpoint string, duration time.Duration) {
	l.WithRequestID(ctx).Info("External service called",
		"service", service,
		"endpoint", endpoint,
		"duration_ms", duration.Milliseconds(),
	)
}

// LogExternalServiceError logs external service errors
func (l *InfrastructureLogger) LogExternalServiceError(ctx context.Context, service string, err error) {
	l.WithRequestID(ctx).WithError(err).Error("External service call failed",
		"service", service,
	)
}

// InterfaceLogger provides logging utilities for the interface layer
type InterfaceLogger struct {
	*Logger
}

// NewInterfaceLogger creates a logger for interface layer
func NewInterfaceLogger(logger *Logger) *InterfaceLogger {
	return &InterfaceLogger{
		Logger: logger.WithComponent("interface"),
	}
}

// LogHTTPRequest logs HTTP requests
func (l *InterfaceLogger) LogHTTPRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration) {
	level := l.getLogLevelForStatusCode(statusCode)
	l.WithRequestID(ctx).Log(ctx, level, "HTTP request processed",
		"method", method,
		"path", path,
		"status_code", statusCode,
		"duration_ms", duration.Milliseconds(),
	)
}

// LogGraphQLOperation logs GraphQL operations
func (l *InterfaceLogger) LogGraphQLOperation(ctx context.Context, operationType, operationName string, duration time.Duration) {
	l.WithRequestID(ctx).Info("GraphQL operation processed",
		"operation_type", operationType,
		"operation_name", operationName,
		"duration_ms", duration.Milliseconds(),
	)
}

// LogGraphQLError logs GraphQL errors
func (l *InterfaceLogger) LogGraphQLError(ctx context.Context, operationName string, err error) {
	l.WithRequestID(ctx).WithError(err).Error("GraphQL operation failed",
		"operation_name", operationName,
	)
}

// getLogLevelForStatusCode returns appropriate log level based on HTTP status code
func (l *InterfaceLogger) getLogLevelForStatusCode(statusCode int) slog.Level {
	switch {
	case statusCode >= 500:
		return slog.LevelError
	case statusCode >= 400:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
