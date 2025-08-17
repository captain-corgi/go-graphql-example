package department

import (
	"context"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// Repository defines the interface for department data access
type Repository interface {
	// FindByID retrieves a department by its ID
	FindByID(ctx context.Context, id DepartmentID) (*Department, error)

	// FindByName retrieves a department by its name
	FindByName(ctx context.Context, name Name) (*Department, error)

	// FindAll retrieves departments with pagination support
	FindAll(ctx context.Context, limit int, cursor string) ([]*Department, string, error)

	// FindByManager retrieves departments by manager ID
	FindByManager(ctx context.Context, managerID EmployeeID, limit int, cursor string) ([]*Department, string, error)

	// Save persists a department (create or update)
	Save(ctx context.Context, department *Department) error

	// Delete removes a department by ID
	Delete(ctx context.Context, id DepartmentID) error

	// Exists checks if a department exists by ID
	Exists(ctx context.Context, id DepartmentID) (bool, error)

	// ExistsByName checks if a department exists by name
	ExistsByName(ctx context.Context, name Name, excludeID *DepartmentID) (bool, error)
}