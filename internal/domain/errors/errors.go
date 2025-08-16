package errors

import "fmt"

// DomainError represents a domain-specific error with code, message, and optional field
type DomainError struct {
	Code    string
	Message string
	Field   string
}

func (e DomainError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

// User domain errors
var (
	ErrUserNotFound      = DomainError{Code: "USER_NOT_FOUND", Message: "User not found"}
	ErrInvalidEmail      = DomainError{Code: "INVALID_EMAIL", Message: "Invalid email format", Field: "email"}
	ErrDuplicateEmail    = DomainError{Code: "DUPLICATE_EMAIL", Message: "Email already exists", Field: "email"}
	ErrInvalidName       = DomainError{Code: "INVALID_NAME", Message: "Name cannot be empty", Field: "name"}
	ErrInvalidUserID     = DomainError{Code: "INVALID_USER_ID", Message: "Invalid user ID format", Field: "id"}
	ErrUserAlreadyExists = DomainError{Code: "USER_ALREADY_EXISTS", Message: "User already exists"}
)

// Repository errors
var (
	ErrRepositoryConnection = DomainError{Code: "REPOSITORY_CONNECTION", Message: "Repository connection failed"}
	ErrRepositoryOperation  = DomainError{Code: "REPOSITORY_OPERATION", Message: "Repository operation failed"}
)
