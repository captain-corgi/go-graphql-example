package leave

import (
	"context"
	"time"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// Repository defines the interface for leave data access
type Repository interface {
	// FindByID retrieves a leave request by its ID
	FindByID(ctx context.Context, id LeaveID) (*Leave, error)

	// FindByEmployee retrieves leave requests by employee ID
	FindByEmployee(ctx context.Context, employeeID EmployeeID, limit int, cursor string) ([]*Leave, string, error)

	// FindByStatus retrieves leave requests by status
	FindByStatus(ctx context.Context, status Status, limit int, cursor string) ([]*Leave, string, error)

	// FindByDateRange retrieves leave requests within a date range
	FindByDateRange(ctx context.Context, startDate, endDate time.Time, limit int, cursor string) ([]*Leave, string, error)

	// FindAll retrieves leave requests with pagination support
	FindAll(ctx context.Context, limit int, cursor string) ([]*Leave, string, error)

	// Save persists a leave request (create or update)
	Save(ctx context.Context, leave *Leave) error

	// Delete removes a leave request by ID
	Delete(ctx context.Context, id LeaveID) error

	// Exists checks if a leave request exists by ID
	Exists(ctx context.Context, id LeaveID) (bool, error)

	// CountByEmployeeAndDateRange counts leave requests for an employee within a date range
	CountByEmployeeAndDateRange(ctx context.Context, employeeID EmployeeID, startDate, endDate time.Time) (int, error)
}