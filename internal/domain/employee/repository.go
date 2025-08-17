package employee

import (
	"context"

	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// Repository defines the interface for employee persistence operations
type Repository interface {
	// FindByID retrieves an employee by their ID
	FindByID(ctx context.Context, id EmployeeID) (*Employee, error)

	// FindByUserID retrieves an employee by their user ID
	FindByUserID(ctx context.Context, userID user.UserID) (*Employee, error)

	// FindByEmployeeCode retrieves an employee by their employee code
	FindByEmployeeCode(ctx context.Context, employeeCode EmployeeCode) (*Employee, error)

	// FindAll retrieves employees with pagination support
	// limit: maximum number of employees to return
	// cursor: pagination cursor (empty string for first page)
	// Returns employees and next cursor (empty if no more pages)
	FindAll(ctx context.Context, limit int, cursor string) ([]*Employee, string, error)

	// FindByDepartment retrieves employees by department
	FindByDepartment(ctx context.Context, department Department, limit int, cursor string) ([]*Employee, string, error)

	// FindByStatus retrieves employees by status
	FindByStatus(ctx context.Context, status Status, limit int, cursor string) ([]*Employee, string, error)

	// Create persists a new employee
	Create(ctx context.Context, employee *Employee) error

	// Update modifies an existing employee
	Update(ctx context.Context, employee *Employee) error

	// Delete removes an employee by their ID
	Delete(ctx context.Context, id EmployeeID) error

	// ExistsByEmployeeCode checks if an employee with the given employee code exists
	ExistsByEmployeeCode(ctx context.Context, employeeCode EmployeeCode) (bool, error)

	// ExistsByUserID checks if an employee with the given user ID exists
	ExistsByUserID(ctx context.Context, userID user.UserID) (bool, error)

	// Count returns the total number of employees
	Count(ctx context.Context) (int64, error)

	// CountByDepartment returns the number of employees in a department
	CountByDepartment(ctx context.Context, department Department) (int64, error)

	// CountByStatus returns the number of employees with a specific status
	CountByStatus(ctx context.Context, status Status) (int64, error)
}