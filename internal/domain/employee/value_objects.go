package employee

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/google/uuid"
)

// EmployeeID represents a unique identifier for an employee
type EmployeeID struct {
	value string
}

// NewEmployeeID creates a new EmployeeID from a string
func NewEmployeeID(value string) (EmployeeID, error) {
	if value == "" {
		return EmployeeID{}, errors.ErrInvalidEmployeeID
	}

	// Validate UUID format
	if _, err := uuid.Parse(value); err != nil {
		return EmployeeID{}, errors.DomainError{
			Code:    "INVALID_EMPLOYEE_ID_FORMAT",
			Message: "Employee ID must be a valid UUID",
			Field:   "id",
		}
	}

	return EmployeeID{value: value}, nil
}

// GenerateEmployeeID creates a new random EmployeeID
func GenerateEmployeeID() EmployeeID {
	return EmployeeID{value: uuid.New().String()}
}

// String returns the string representation of the EmployeeID
func (e EmployeeID) String() string {
	return e.value
}

// Equals checks if two EmployeeIDs are equal
func (e EmployeeID) Equals(other EmployeeID) bool {
	return e.value == other.value
}

// EmployeeCode represents an employee's unique code
type EmployeeCode struct {
	value string
}

// NewEmployeeCode creates a new EmployeeCode from a string
func NewEmployeeCode(value string) (EmployeeCode, error) {
	if value == "" {
		return EmployeeCode{}, errors.ErrInvalidEmployeeCode
	}

	// Remove leading/trailing whitespace
	value = strings.TrimSpace(value)

	// Validate length (3-10 characters)
	if len(value) < 3 || len(value) > 10 {
		return EmployeeCode{}, errors.DomainError{
			Code:    "INVALID_EMPLOYEE_CODE_LENGTH",
			Message: "Employee code must be between 3 and 10 characters",
			Field:   "employeeCode",
		}
	}

	// Validate format (alphanumeric and hyphens only)
	matched, err := regexp.MatchString(`^[A-Za-z0-9-]+$`, value)
	if err != nil || !matched {
		return EmployeeCode{}, errors.DomainError{
			Code:    "INVALID_EMPLOYEE_CODE_FORMAT",
			Message: "Employee code must contain only letters, numbers, and hyphens",
			Field:   "employeeCode",
		}
	}

	return EmployeeCode{value: strings.ToUpper(value)}, nil
}

// String returns the string representation of the EmployeeCode
func (e EmployeeCode) String() string {
	return e.value
}

// Equals checks if two EmployeeCodes are equal
func (e EmployeeCode) Equals(other EmployeeCode) bool {
	return e.value == other.value
}

// Department represents an employee's department
type Department struct {
	value string
}

// NewDepartment creates a new Department from a string
func NewDepartment(value string) (Department, error) {
	if value == "" {
		return Department{}, errors.ErrInvalidDepartment
	}

	// Remove leading/trailing whitespace
	value = strings.TrimSpace(value)

	// Validate length (2-50 characters)
	if len(value) < 2 || len(value) > 50 {
		return Department{}, errors.DomainError{
			Code:    "INVALID_DEPARTMENT_LENGTH",
			Message: "Department must be between 2 and 50 characters",
			Field:   "department",
		}
	}

	// Validate format (letters, spaces, and hyphens only)
	matched, err := regexp.MatchString(`^[A-Za-z\s-]+$`, value)
	if err != nil || !matched {
		return Department{}, errors.DomainError{
			Code:    "INVALID_DEPARTMENT_FORMAT",
			Message: "Department must contain only letters, spaces, and hyphens",
			Field:   "department",
		}
	}

	return Department{value: value}, nil
}

// String returns the string representation of the Department
func (d Department) String() string {
	return d.value
}

// Equals checks if two Departments are equal
func (d Department) Equals(other Department) bool {
	return d.value == other.value
}

// Position represents an employee's position/title
type Position struct {
	value string
}

// NewPosition creates a new Position from a string
func NewPosition(value string) (Position, error) {
	if value == "" {
		return Position{}, errors.ErrInvalidPosition
	}

	// Remove leading/trailing whitespace
	value = strings.TrimSpace(value)

	// Validate length (2-50 characters)
	if len(value) < 2 || len(value) > 50 {
		return Position{}, errors.DomainError{
			Code:    "INVALID_POSITION_LENGTH",
			Message: "Position must be between 2 and 50 characters",
			Field:   "position",
		}
	}

	// Validate format (letters, spaces, and hyphens only)
	matched, err := regexp.MatchString(`^[A-Za-z\s-]+$`, value)
	if err != nil || !matched {
		return Position{}, errors.DomainError{
			Code:    "INVALID_POSITION_FORMAT",
			Message: "Position must contain only letters, spaces, and hyphens",
			Field:   "position",
		}
	}

	return Position{value: value}, nil
}

// String returns the string representation of the Position
func (p Position) String() string {
	return p.value
}

// Equals checks if two Positions are equal
func (p Position) Equals(other Position) bool {
	return p.value == other.value
}

// Salary represents an employee's salary
type Salary struct {
	value float64
}

// NewSalary creates a new Salary from a float64
func NewSalary(value float64) (Salary, error) {
	if value <= 0 {
		return Salary{}, errors.ErrInvalidSalary
	}

	// Validate maximum salary (1 million)
	if value > 1000000 {
		return Salary{}, errors.DomainError{
			Code:    "INVALID_SALARY_RANGE",
			Message: "Salary cannot exceed 1,000,000",
			Field:   "salary",
		}
	}

	return Salary{value: value}, nil
}

// Value returns the float64 value of the Salary
func (s Salary) Value() float64 {
	return s.value
}

// Equals checks if two Salaries are equal
func (s Salary) Equals(other Salary) bool {
	return s.value == other.value
}

// String returns the string representation of the Salary
func (s Salary) String() string {
	return fmt.Sprintf("%.2f", s.value)
}

// Status represents an employee's employment status
type Status struct {
	value string
}

// Valid status values
const (
	StatusActive   = "ACTIVE"
	StatusInactive = "INACTIVE"
	StatusTerminated = "TERMINATED"
	StatusOnLeave  = "ON_LEAVE"
)

// NewStatus creates a new Status from a string
func NewStatus(value string) (Status, error) {
	if value == "" {
		return Status{}, errors.ErrInvalidStatus
	}

	// Remove leading/trailing whitespace and convert to uppercase
	value = strings.TrimSpace(strings.ToUpper(value))

	// Validate status value
	switch value {
	case StatusActive, StatusInactive, StatusTerminated, StatusOnLeave:
		return Status{value: value}, nil
	default:
		return Status{}, errors.DomainError{
			Code:    "INVALID_STATUS_VALUE",
			Message: fmt.Sprintf("Status must be one of: %s, %s, %s, %s", StatusActive, StatusInactive, StatusTerminated, StatusOnLeave),
			Field:   "status",
		}
	}
}

// String returns the string representation of the Status
func (s Status) String() string {
	return s.value
}

// Equals checks if two Statuses are equal
func (s Status) Equals(other Status) bool {
	return s.value == other.value
}

// IsActive checks if the status is active
func (s Status) IsActive() bool {
	return s.value == StatusActive
}

// IsInactive checks if the status is inactive
func (s Status) IsInactive() bool {
	return s.value == StatusInactive
}

// IsTerminated checks if the status is terminated
func (s Status) IsTerminated() bool {
	return s.value == StatusTerminated
}

// IsOnLeave checks if the status is on leave
func (s Status) IsOnLeave() bool {
	return s.value == StatusOnLeave
}