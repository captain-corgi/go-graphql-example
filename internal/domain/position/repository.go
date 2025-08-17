package position

import (
	"context"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// Repository defines the interface for position data access
type Repository interface {
	// FindByID retrieves a position by its ID
	FindByID(ctx context.Context, id PositionID) (*Position, error)

	// FindByTitle retrieves a position by its title
	FindByTitle(ctx context.Context, title Title) (*Position, error)

	// FindAll retrieves positions with pagination support
	FindAll(ctx context.Context, limit int, cursor string) ([]*Position, string, error)

	// FindByDepartment retrieves positions by department ID
	FindByDepartment(ctx context.Context, departmentID DepartmentID, limit int, cursor string) ([]*Position, string, error)

	// Save persists a position (create or update)
	Save(ctx context.Context, position *Position) error

	// Delete removes a position by ID
	Delete(ctx context.Context, id PositionID) error

	// Exists checks if a position exists by ID
	Exists(ctx context.Context, id PositionID) (bool, error)

	// ExistsByTitle checks if a position exists by title
	ExistsByTitle(ctx context.Context, title Title, excludeID *PositionID) (bool, error)
}