package position

import (
	"context"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// DomainService defines domain-specific business logic that doesn't belong to a single entity
type DomainService interface {
	// ValidateUniqueTitle ensures position title uniqueness across the domain
	ValidateUniqueTitle(ctx context.Context, title Title, excludePositionID *PositionID) error

	// CanDeletePosition checks if a position can be safely deleted
	CanDeletePosition(ctx context.Context, positionID PositionID) error

	// ValidateDepartmentExists checks if the assigned department exists and is valid
	ValidateDepartmentExists(ctx context.Context, departmentID *DepartmentID) error
}

// domainService implements the DomainService interface
type domainService struct {
	positionRepo Repository
}

// NewDomainService creates a new domain service instance
func NewDomainService(positionRepo Repository) DomainService {
	return &domainService{
		positionRepo: positionRepo,
	}
}

// ValidateUniqueTitle ensures position title uniqueness across the domain
func (s *domainService) ValidateUniqueTitle(ctx context.Context, title Title, excludePositionID *PositionID) error {
	existingPosition, err := s.positionRepo.FindByTitle(ctx, title)
	if err != nil {
		// If position not found, title is unique
		if err == errors.ErrPositionNotFound {
			return nil
		}
		return err
	}

	// If we're excluding a specific position ID (for updates), check if it's the same position
	if excludePositionID != nil && existingPosition.ID().Equals(*excludePositionID) {
		return nil
	}

	return errors.ErrDuplicatePositionTitle
}

// CanDeletePosition checks if a position can be safely deleted
func (s *domainService) CanDeletePosition(ctx context.Context, positionID PositionID) error {
	// Check if position exists
	_, err := s.positionRepo.FindByID(ctx, positionID)
	if err != nil {
		return err
	}

	// Add any business rules for deletion here
	// For example: check if position has active employees, etc.
	// For now, we allow all deletions

	return nil
}

// ValidateDepartmentExists checks if the assigned department exists and is valid
func (s *domainService) ValidateDepartmentExists(ctx context.Context, departmentID *DepartmentID) error {
	if departmentID == nil {
		// No department assigned is valid
		return nil
	}

	// In a real implementation, you would check against the department repository
	// For now, we'll just validate the format
	if departmentID.String() == "" {
		return errors.ErrInvalidDepartmentID
	}

	return nil
}