package resolver

import (
	"log/slog"

	"github.com/captain-corgi/go-graphql-example/internal/application/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver holds the dependencies for GraphQL resolvers
type Resolver struct {
	userService user.Service
	logger      *slog.Logger
}

// NewResolver creates a new resolver with the given dependencies
func NewResolver(userService user.Service, logger *slog.Logger) *Resolver {
	return &Resolver{
		userService: userService,
		logger:      logger,
	}
}
