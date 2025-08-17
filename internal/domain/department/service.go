package department

import (
	"context"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// DomainService defines domain-specific business logic that doesn't belong to a single entity
type DomainService interface {
	// ValidateUniqueName ensures department name uniqueness across the domain
	ValidateUniqueName(ctx context.Context, name Name, excludeDepartmentID *DepartmentID) error

	// CanDeleteDepartment checks if a department can be safely deleted
	CanDeleteDepartment(ctx context.Context, departmentID DepartmentID) error

	// ValidateManagerExists checks if the assigned manager exists and is valid
	ValidateManagerExists(ctx context.Context, managerID *EmployeeID) error
}

// domainService implements the DomainService interface
type domainService struct {
	deptRepo Repository
}

// NewDomainService creates a new domain service instance
func NewDomainService(deptRepo Repository) DomainService {
	return &domainService{
		deptRepo: deptRepo,
	}
}

// ValidateUniqueName ensures department name uniqueness across the domain
func (s *domainService) ValidateUniqueName(ctx context.Context, name Name, excludeDepartmentID *DepartmentID) error {
	existingDepartment, err := s.deptRepo.FindByName(ctx, name)
	if err != nil {
		// If department not found, name is unique
		if err == errors.ErrDepartmentNotFound {
			return nil
		}
		return err
	}

	// If we're excluding a specific department ID (for updates), check if it's the same department
	if excludeDepartmentID != nil && existingDepartment.ID().Equals(*excludeDepartmentID) {
		return nil
	}

	return errors.ErrDuplicateDepartmentName
}

// CanDeleteDepartment checks if a department can be safely deleted
func (s *domainService) CanDeleteDepartment(ctx context.Context, departmentID DepartmentID) error {
	// Check if department exists
	_, err := s.deptRepo.FindByID(ctx, departmentID)
	if err != nil {
		return err
	}

	// Add any business rules for deletion here
	// For example: check if department has active employees, etc.
	// For now, we allow all deletions

	return nil
}

// ValidateManagerExists checks if the assigned manager exists and is valid
func (s *domainService) ValidateManagerExists(ctx context.Context, managerID *EmployeeID) error {
	if managerID == nil {
		// No manager assigned is valid
		return nil
	}

	// In a real implementation, you would check against the employee repository
	// For now, we'll just validate the format
	if managerID.String() == "" {
		return errors.ErrInvalidEmployeeID
	}

	return nil
}