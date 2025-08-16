package user

import "time"

// Request DTOs

// GetUserRequest represents a request to get a user by ID
type GetUserRequest struct {
	ID string `json:"id" validate:"required"`
}

// ListUsersRequest represents a request to list users with pagination
type ListUsersRequest struct {
	First int    `json:"first" validate:"min=1,max=100"`
	After string `json:"after"`
}

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=1,max=255"`
}

// UpdateUserRequest represents a request to update an existing user
type UpdateUserRequest struct {
	ID    string  `json:"id" validate:"required"`
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
	Name  *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
}

// DeleteUserRequest represents a request to delete a user
type DeleteUserRequest struct {
	ID string `json:"id" validate:"required"`
}

// Response DTOs

// UserDTO represents a user in the application layer
type UserDTO struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ErrorDTO represents an error in the application layer
type ErrorDTO struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Code    string `json:"code"`
}

// PageInfoDTO represents pagination information
type PageInfoDTO struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor"`
	EndCursor       *string `json:"endCursor"`
}

// UserEdgeDTO represents a user edge in a connection
type UserEdgeDTO struct {
	Node   *UserDTO `json:"node"`
	Cursor string   `json:"cursor"`
}

// UserConnectionDTO represents a paginated list of users
type UserConnectionDTO struct {
	Edges    []*UserEdgeDTO `json:"edges"`
	PageInfo *PageInfoDTO   `json:"pageInfo"`
}

// Response DTOs

// GetUserResponse represents the response for getting a user
type GetUserResponse struct {
	User   *UserDTO   `json:"user"`
	Errors []ErrorDTO `json:"errors,omitempty"`
}

// ListUsersResponse represents the response for listing users
type ListUsersResponse struct {
	Users  *UserConnectionDTO `json:"users"`
	Errors []ErrorDTO         `json:"errors,omitempty"`
}

// CreateUserResponse represents the response for creating a user
type CreateUserResponse struct {
	User   *UserDTO   `json:"user"`
	Errors []ErrorDTO `json:"errors,omitempty"`
}

// UpdateUserResponse represents the response for updating a user
type UpdateUserResponse struct {
	User   *UserDTO   `json:"user"`
	Errors []ErrorDTO `json:"errors,omitempty"`
}

// DeleteUserResponse represents the response for deleting a user
type DeleteUserResponse struct {
	Success bool       `json:"success"`
	Errors  []ErrorDTO `json:"errors,omitempty"`
}
