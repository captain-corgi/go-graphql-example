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

// Department domain errors
var (
	ErrDepartmentNotFound        = DomainError{Code: "DEPARTMENT_NOT_FOUND", Message: "Department not found"}
	ErrInvalidDepartmentID       = DomainError{Code: "INVALID_DEPARTMENT_ID", Message: "Invalid department ID format", Field: "id"}
	ErrInvalidDepartmentName     = DomainError{Code: "INVALID_DEPARTMENT_NAME", Message: "Department name cannot be empty", Field: "name"}
	ErrDuplicateDepartmentName   = DomainError{Code: "DUPLICATE_DEPARTMENT_NAME", Message: "Department name already exists", Field: "name"}
	ErrInvalidDepartmentDescription = DomainError{Code: "INVALID_DEPARTMENT_DESCRIPTION", Message: "Department description cannot be empty", Field: "description"}
	ErrDepartmentAlreadyExists   = DomainError{Code: "DEPARTMENT_ALREADY_EXISTS", Message: "Department already exists"}
)

// Position domain errors
var (
	ErrPositionNotFound        = DomainError{Code: "POSITION_NOT_FOUND", Message: "Position not found"}
	ErrInvalidPositionID       = DomainError{Code: "INVALID_POSITION_ID", Message: "Invalid position ID format", Field: "id"}
	ErrInvalidPositionTitle    = DomainError{Code: "INVALID_POSITION_TITLE", Message: "Position title cannot be empty", Field: "title"}
	ErrDuplicatePositionTitle  = DomainError{Code: "DUPLICATE_POSITION_TITLE", Message: "Position title already exists", Field: "title"}
	ErrInvalidPositionDescription = DomainError{Code: "INVALID_POSITION_DESCRIPTION", Message: "Position description cannot be empty", Field: "description"}
	ErrInvalidPositionRequirements = DomainError{Code: "INVALID_POSITION_REQUIREMENTS", Message: "Position requirements cannot be empty", Field: "requirements"}
	ErrPositionAlreadyExists   = DomainError{Code: "POSITION_ALREADY_EXISTS", Message: "Position already exists"}
)

// Leave domain errors
var (
	ErrLeaveNotFound        = DomainError{Code: "LEAVE_NOT_FOUND", Message: "Leave request not found"}
	ErrInvalidLeaveID       = DomainError{Code: "INVALID_LEAVE_ID", Message: "Invalid leave ID format", Field: "id"}
	ErrInvalidLeaveType     = DomainError{Code: "INVALID_LEAVE_TYPE", Message: "Invalid leave type", Field: "leaveType"}
	ErrInvalidLeaveStatus   = DomainError{Code: "INVALID_LEAVE_STATUS", Message: "Invalid leave status", Field: "status"}
	ErrInvalidLeaveReason   = DomainError{Code: "INVALID_LEAVE_REASON", Message: "Leave reason cannot be empty", Field: "reason"}
	ErrLeaveAlreadyExists   = DomainError{Code: "LEAVE_ALREADY_EXISTS", Message: "Leave request already exists"}
)

// Repository errors
var (
	ErrRepositoryConnection = DomainError{Code: "REPOSITORY_CONNECTION", Message: "Repository connection failed"}
	ErrRepositoryOperation  = DomainError{Code: "REPOSITORY_OPERATION", Message: "Repository operation failed"}
)
