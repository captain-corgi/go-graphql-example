package resolver

import (
	"context"
	"log/slog"

	"github.com/captain-corgi/go-graphql-example/internal/application/auth"
	"github.com/captain-corgi/go-graphql-example/internal/application/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver holds the dependencies for GraphQL resolvers
type Resolver struct {
	userService user.Service
	authService auth.Service
	logger      *slog.Logger
}

// NewResolver creates a new resolver with the given dependencies
func NewResolver(userService user.Service, authService auth.Service, logger *slog.Logger) *Resolver {
	return &Resolver{
		userService: userService,
		authService: authService,
		logger:      logger,
	}
}

// getClientInfo extracts device and IP information from the request context
func (r *Resolver) getClientInfo(ctx context.Context) (*string, *string) {
	// TODO: Extract device info from User-Agent header if available
	// TODO: Extract IP address from X-Forwarded-For or RemoteAddr if available
	// For now, return nil values
	return nil, nil
}
