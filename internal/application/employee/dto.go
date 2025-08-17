package employee

import (
	"time"
)

// GetEmployeeRequest represents a request to get an employee by ID
type GetEmployeeRequest struct {
	ID string
}

// GetEmployeeResponse represents a response containing an employee
type GetEmployeeResponse struct {
	Employee *EmployeeDTO
}

// ListEmployeesRequest represents a request to list employees with pagination
type ListEmployeesRequest struct {
	Limit  int
	Cursor string
}

// ListEmployeesResponse represents a response containing a list of employees
type ListEmployeesResponse struct {
	Employees []*EmployeeDTO
	NextCursor string
}

// ListEmployeesByDepartmentRequest represents a request to list employees by department
type ListEmployeesByDepartmentRequest struct {
	Department string
	Limit      int
	Cursor     string
}

// ListEmployeesByDepartmentResponse represents a response containing a list of employees by department
type ListEmployeesByDepartmentResponse struct {
	Employees  []*EmployeeDTO
	NextCursor string
}

// ListEmployeesByStatusRequest represents a request to list employees by status
type ListEmployeesByStatusRequest struct {
	Status string
	Limit  int
	Cursor string
}

// ListEmployeesByStatusResponse represents a response containing a list of employees by status
type ListEmployeesByStatusResponse struct {
	Employees  []*EmployeeDTO
	NextCursor string
}

// CreateEmployeeRequest represents a request to create a new employee
type CreateEmployeeRequest struct {
	UserID       string
	EmployeeCode string
	Department   string
	Position     string
	HireDate     time.Time
	Salary       float64
	Status       string
}

// CreateEmployeeResponse represents a response containing the created employee
type CreateEmployeeResponse struct {
	Employee *EmployeeDTO
	Errors   []ErrorDTO
}

// UpdateEmployeeRequest represents a request to update an employee
type UpdateEmployeeRequest struct {
	ID           string
	EmployeeCode *string
	Department   *string
	Position     *string
	HireDate     *time.Time
	Salary       *float64
	Status       *string
}

// UpdateEmployeeResponse represents a response containing the updated employee
type UpdateEmployeeResponse struct {
	Employee *EmployeeDTO
	Errors   []ErrorDTO
}

// DeleteEmployeeRequest represents a request to delete an employee
type DeleteEmployeeRequest struct {
	ID string
}

// DeleteEmployeeResponse represents a response indicating the result of deletion
type DeleteEmployeeResponse struct {
	Success bool
	Errors  []ErrorDTO
}

// EmployeeDTO represents an employee data transfer object
type EmployeeDTO struct {
	ID           string
	UserID       string
	User         *UserDTO
	EmployeeCode string
	Department   string
	Position     string
	HireDate     time.Time
	Salary       float64
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserDTO represents a user data transfer object (for employee relationship)
type UserDTO struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ErrorDTO represents an error data transfer object
type ErrorDTO struct {
	Message string
	Field   string
	Code    string
}