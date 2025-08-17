package department

import (
	"time"
)

// DepartmentDTO represents a department data transfer object
type DepartmentDTO struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ManagerID   *string    `json:"managerId,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// DepartmentConnectionDTO represents a paginated connection of departments
type DepartmentConnectionDTO struct {
	Departments []*DepartmentDTO `json:"departments"`
	NextCursor  string           `json:"nextCursor"`
}

// GetDepartmentRequest represents a request to get a department by ID
type GetDepartmentRequest struct {
	ID string `json:"id"`
}

// GetDepartmentResponse represents a response for getting a department
type GetDepartmentResponse struct {
	Department *DepartmentDTO `json:"department,omitempty"`
	Errors     []ErrorDTO     `json:"errors,omitempty"`
}

// ListDepartmentsRequest represents a request to list departments
type ListDepartmentsRequest struct {
	Limit  int    `json:"limit"`
	Cursor string `json:"cursor,omitempty"`
}

// ListDepartmentsResponse represents a response for listing departments
type ListDepartmentsResponse struct {
	Departments *DepartmentConnectionDTO `json:"departments,omitempty"`
	Errors      []ErrorDTO               `json:"errors,omitempty"`
}

// ListDepartmentsByManagerRequest represents a request to list departments by manager
type ListDepartmentsByManagerRequest struct {
	ManagerID string `json:"managerId"`
	Limit     int    `json:"limit"`
	Cursor    string `json:"cursor,omitempty"`
}

// ListDepartmentsByManagerResponse represents a response for listing departments by manager
type ListDepartmentsByManagerResponse struct {
	Departments *DepartmentConnectionDTO `json:"departments,omitempty"`
	Errors      []ErrorDTO               `json:"errors,omitempty"`
}

// CreateDepartmentRequest represents a request to create a department
type CreateDepartmentRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ManagerID   *string `json:"managerId,omitempty"`
}

// CreateDepartmentResponse represents a response for creating a department
type CreateDepartmentResponse struct {
	Department *DepartmentDTO `json:"department,omitempty"`
	Errors     []ErrorDTO     `json:"errors,omitempty"`
}

// UpdateDepartmentRequest represents a request to update a department
type UpdateDepartmentRequest struct {
	ID          string  `json:"id"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	ManagerID   *string `json:"managerId,omitempty"`
}

// UpdateDepartmentResponse represents a response for updating a department
type UpdateDepartmentResponse struct {
	Department *DepartmentDTO `json:"department,omitempty"`
	Errors     []ErrorDTO     `json:"errors,omitempty"`
}

// DeleteDepartmentRequest represents a request to delete a department
type DeleteDepartmentRequest struct {
	ID string `json:"id"`
}

// DeleteDepartmentResponse represents a response for deleting a department
type DeleteDepartmentResponse struct {
	Success bool       `json:"success"`
	Errors  []ErrorDTO `json:"errors,omitempty"`
}

// ErrorDTO represents an error data transfer object
type ErrorDTO struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Code    string `json:"code,omitempty"`
}