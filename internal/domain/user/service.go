package user

import (
	"context"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// DomainService defines domain-specific business logic that doesn't belong to a single entity
type DomainService interface {
	// ValidateUniqueEmail ensures email uniqueness across the domain
	ValidateUniqueEmail(ctx context.Context, email Email, excludeUserID *UserID) error

	// CanDeleteUser checks if a user can be safely deleted
	CanDeleteUser(ctx context.Context, userID UserID) error
}

// domainService implements the DomainService interface
type domainService struct {
	userRepo Repository
}

// NewDomainService creates a new domain service instance
func NewDomainService(userRepo Repository) DomainService {
	return &domainService{
		userRepo: userRepo,
	}
}

// ValidateUniqueEmail ensures email uniqueness across the domain
func (s *domainService) ValidateUniqueEmail(ctx context.Context, email Email, excludeUserID *UserID) error {
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// If user not found, email is unique
		if err == errors.ErrUserNotFound {
			return nil
		}
		return err
	}

	// If we're excluding a specific user ID (for updates), check if it's the same user
	if excludeUserID != nil && existingUser.ID().Equals(*excludeUserID) {
		return nil
	}

	return errors.ErrDuplicateEmail
}

// CanDeleteUser checks if a user can be safely deleted
func (s *domainService) CanDeleteUser(ctx context.Context, userID UserID) error {
	// Check if user exists
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Add any business rules for deletion here
	// For example: check if user has active orders, posts, etc.
	// For now, we allow all deletions

	return nil
}
