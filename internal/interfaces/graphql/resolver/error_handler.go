package resolver

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	domainErrors "github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// handleGraphQLError handles errors and converts them to appropriate GraphQL errors
// It also logs errors with proper context for debugging
func (r *Resolver) handleGraphQLError(ctx context.Context, err error, operation string) error {
	if err == nil {
		return nil
	}

	// Extract request ID from context for tracing
	requestID := getRequestIDFromContext(ctx)

	// Log the error with context
	r.logger.ErrorContext(ctx, "GraphQL operation error",
		"operation", operation,
		"error", err.Error(),
		"request_id", requestID,
	)

	// Check if it's a domain error
	var domainErr domainErrors.DomainError
	if errors.As(err, &domainErr) {
		return &gqlerror.Error{
			Message: domainErr.Message,
			Extensions: map[string]interface{}{
				"code":       domainErr.Code,
				"field":      domainErr.Field,
				"request_id": requestID,
			},
		}
	}

	// For unknown errors, return a generic error to avoid exposing internal details
	return &gqlerror.Error{
		Message: "An internal error occurred",
		Extensions: map[string]interface{}{
			"code":       "INTERNAL_ERROR",
			"request_id": requestID,
		},
	}
}

// validateInput performs basic input validation and sanitization
func (r *Resolver) validateInput(ctx context.Context, operation string, validator func() error) error {
	if err := validator(); err != nil {
		r.logger.WarnContext(ctx, "Input validation failed",
			"operation", operation,
			"error", err.Error(),
		)
		return r.handleGraphQLError(ctx, err, operation)
	}
	return nil
}

// getRequestIDFromContext extracts the request ID from the GraphQL context
func getRequestIDFromContext(ctx context.Context) string {
	// Try to get request ID from GraphQL context (safely handle missing context)
	defer func() {
		if r := recover(); r != nil {
			// GraphQL context not available, continue with regular context
			// Log the panic for debugging purposes
			_ = r
		}
	}()

	if reqCtx := graphql.GetOperationContext(ctx); reqCtx != nil {
		if requestID, ok := reqCtx.Variables["requestId"]; ok {
			if id, ok := requestID.(string); ok {
				return id
			}
		}
	}

	// Try to get request ID from regular context
	if requestID := ctx.Value("request_id"); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}

	return "unknown"
}

// logOperation logs the start of a GraphQL operation with context
func (r *Resolver) logOperation(ctx context.Context, operation string, args map[string]interface{}) {
	requestID := getRequestIDFromContext(ctx)

	logArgs := []interface{}{
		"operation", operation,
		"request_id", requestID,
	}

	// Add operation-specific arguments to log
	for key, value := range args {
		logArgs = append(logArgs, key, value)
	}

	r.logger.InfoContext(ctx, "GraphQL operation started", logArgs...)
}

// logOperationSuccess logs successful completion of a GraphQL operation
func (r *Resolver) logOperationSuccess(ctx context.Context, operation string, result interface{}) {
	requestID := getRequestIDFromContext(ctx)

	r.logger.InfoContext(ctx, "GraphQL operation completed successfully",
		"operation", operation,
		"request_id", requestID,
		"has_result", result != nil,
	)
}
