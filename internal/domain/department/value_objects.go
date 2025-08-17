package department

import (
	"regexp"
	"strings"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/google/uuid"
)

// DepartmentID represents a department identifier
type DepartmentID struct {
	value string
}

// NewDepartmentID creates a new DepartmentID from a string
func NewDepartmentID(value string) (DepartmentID, error) {
	if value == "" {
		return DepartmentID{}, errors.ErrInvalidDepartmentID
	}

	// Validate UUID format
	if _, err := uuid.Parse(value); err != nil {
		return DepartmentID{}, errors.ErrInvalidDepartmentID
	}

	return DepartmentID{value: value}, nil
}

// GenerateDepartmentID generates a new DepartmentID
func GenerateDepartmentID() DepartmentID {
	return DepartmentID{value: uuid.New().String()}
}

// String returns the string representation of the DepartmentID
func (d DepartmentID) String() string {
	return d.value
}

// Equals checks if two DepartmentIDs are equal
func (d DepartmentID) Equals(other DepartmentID) bool {
	return d.value == other.value
}

// Name represents a department name
type Name struct {
	value string
}

// NewName creates a new Name from a string
func NewName(value string) (Name, error) {
	if value == "" {
		return Name{}, errors.ErrInvalidDepartmentName
	}

	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return Name{}, errors.ErrInvalidDepartmentName
	}

	if len(trimmed) > 100 {
		return Name{}, errors.DomainError{
			Code:    "DEPARTMENT_NAME_TOO_LONG",
			Message: "Department name cannot exceed 100 characters",
			Field:   "name",
		}
	}

	// Check for valid characters (letters, numbers, spaces, hyphens, underscores)
	validNameRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\-_]+$`)
	if !validNameRegex.MatchString(trimmed) {
		return Name{}, errors.DomainError{
			Code:    "INVALID_DEPARTMENT_NAME_CHARACTERS",
			Message: "Department name can only contain letters, numbers, spaces, hyphens, and underscores",
			Field:   "name",
		}
	}

	return Name{value: trimmed}, nil
}

// String returns the string representation of the Name
func (n Name) String() string {
	return n.value
}

// Equals checks if two Names are equal
func (n Name) Equals(other Name) bool {
	return n.value == other.value
}

// Description represents a department description
type Description struct {
	value string
}

// NewDescription creates a new Description from a string
func NewDescription(value string) (Description, error) {
	if value == "" {
		return Description{}, errors.ErrInvalidDepartmentDescription
	}

	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return Description{}, errors.ErrInvalidDepartmentDescription
	}

	if len(trimmed) > 500 {
		return Description{}, errors.DomainError{
			Code:    "DEPARTMENT_DESCRIPTION_TOO_LONG",
			Message: "Department description cannot exceed 500 characters",
			Field:   "description",
		}
	}

	return Description{value: trimmed}, nil
}

// String returns the string representation of the Description
func (d Description) String() string {
	return d.value
}

// Equals checks if two Descriptions are equal
func (d Description) Equals(other Description) bool {
	return d.value == other.value
}

// EmployeeID represents an employee identifier (reused from employee domain)
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
		return EmployeeID{}, errors.ErrInvalidEmployeeID
	}

	return EmployeeID{value: value}, nil
}

// String returns the string representation of the EmployeeID
func (e EmployeeID) String() string {
	return e.value
}

// Equals checks if two EmployeeIDs are equal
func (e EmployeeID) Equals(other EmployeeID) bool {
	return e.value == other.value
}

// ValidateDepartmentID validates a department ID string
func ValidateDepartmentID(id string) error {
	if id == "" {
		return errors.ErrInvalidDepartmentID
	}

	if _, err := uuid.Parse(id); err != nil {
		return errors.ErrInvalidDepartmentID
	}

	return nil
}

// ValidateDepartmentName validates a department name string
func ValidateDepartmentName(name string) error {
	if name == "" {
		return errors.ErrInvalidDepartmentName
	}

	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return errors.ErrInvalidDepartmentName
	}

	if len(trimmed) > 100 {
		return errors.DomainError{
			Code:    "DEPARTMENT_NAME_TOO_LONG",
			Message: "Department name cannot exceed 100 characters",
			Field:   "name",
		}
	}

	validNameRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\-_]+$`)
	if !validNameRegex.MatchString(trimmed) {
		return errors.DomainError{
			Code:    "INVALID_DEPARTMENT_NAME_CHARACTERS",
			Message: "Department name can only contain letters, numbers, spaces, hyphens, and underscores",
			Field:   "name",
		}
	}

	return nil
}

// ValidateDepartmentDescription validates a department description string
func ValidateDepartmentDescription(description string) error {
	if description == "" {
		return errors.ErrInvalidDepartmentDescription
	}

	trimmed := strings.TrimSpace(description)
	if trimmed == "" {
		return errors.ErrInvalidDepartmentDescription
	}

	if len(trimmed) > 500 {
		return errors.DomainError{
			Code:    "DEPARTMENT_DESCRIPTION_TOO_LONG",
			Message: "Department description cannot exceed 500 characters",
			Field:   "description",
		}
	}

	return nil
}