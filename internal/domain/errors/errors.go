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

// Employee domain errors
var (
	ErrEmployeeNotFound      = DomainError{Code: "EMPLOYEE_NOT_FOUND", Message: "Employee not found"}
	ErrInvalidEmployeeCode   = DomainError{Code: "INVALID_EMPLOYEE_CODE", Message: "Employee code cannot be empty", Field: "employeeCode"}
	ErrDuplicateEmployeeCode = DomainError{Code: "DUPLICATE_EMPLOYEE_CODE", Message: "Employee code already exists", Field: "employeeCode"}
	ErrInvalidDepartment     = DomainError{Code: "INVALID_DEPARTMENT", Message: "Department cannot be empty", Field: "department"}
	ErrInvalidPosition       = DomainError{Code: "INVALID_POSITION", Message: "Position cannot be empty", Field: "position"}
	ErrInvalidSalary         = DomainError{Code: "INVALID_SALARY", Message: "Salary must be greater than zero", Field: "salary"}
	ErrInvalidStatus         = DomainError{Code: "INVALID_STATUS", Message: "Status cannot be empty", Field: "status"}
	ErrInvalidEmployeeID     = DomainError{Code: "INVALID_EMPLOYEE_ID", Message: "Invalid employee ID format", Field: "id"}
	ErrEmployeeAlreadyExists = DomainError{Code: "EMPLOYEE_ALREADY_EXISTS", Message: "Employee already exists"}
)

// Repository errors
var (
	ErrRepositoryConnection = DomainError{Code: "REPOSITORY_CONNECTION", Message: "Repository connection failed"}
	ErrRepositoryOperation  = DomainError{Code: "REPOSITORY_OPERATION", Message: "Repository operation failed"}
)
