package position

import (
	"regexp"
	"strings"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/google/uuid"
)

// PositionID represents a position identifier
type PositionID struct {
	value string
}

// NewPositionID creates a new PositionID from a string
func NewPositionID(value string) (PositionID, error) {
	if value == "" {
		return PositionID{}, errors.ErrInvalidPositionID
	}

	// Validate UUID format
	if _, err := uuid.Parse(value); err != nil {
		return PositionID{}, errors.ErrInvalidPositionID
	}

	return PositionID{value: value}, nil
}

// GeneratePositionID generates a new PositionID
func GeneratePositionID() PositionID {
	return PositionID{value: uuid.New().String()}
}

// String returns the string representation of the PositionID
func (p PositionID) String() string {
	return p.value
}

// Equals checks if two PositionIDs are equal
func (p PositionID) Equals(other PositionID) bool {
	return p.value == other.value
}

// Title represents a position title
type Title struct {
	value string
}

// NewTitle creates a new Title from a string
func NewTitle(value string) (Title, error) {
	if value == "" {
		return Title{}, errors.ErrInvalidPositionTitle
	}

	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return Title{}, errors.ErrInvalidPositionTitle
	}

	if len(trimmed) > 100 {
		return Title{}, errors.DomainError{
			Code:    "POSITION_TITLE_TOO_LONG",
			Message: "Position title cannot exceed 100 characters",
			Field:   "title",
		}
	}

	// Check for valid characters (letters, numbers, spaces, hyphens, underscores)
	validTitleRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\-_]+$`)
	if !validTitleRegex.MatchString(trimmed) {
		return Title{}, errors.DomainError{
			Code:    "INVALID_POSITION_TITLE_CHARACTERS",
			Message: "Position title can only contain letters, numbers, spaces, hyphens, and underscores",
			Field:   "title",
		}
	}

	return Title{value: trimmed}, nil
}

// String returns the string representation of the Title
func (t Title) String() string {
	return t.value
}

// Equals checks if two Titles are equal
func (t Title) Equals(other Title) bool {
	return t.value == other.value
}

// Description represents a position description
type Description struct {
	value string
}

// NewDescription creates a new Description from a string
func NewDescription(value string) (Description, error) {
	if value == "" {
		return Description{}, errors.ErrInvalidPositionDescription
	}

	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return Description{}, errors.ErrInvalidPositionDescription
	}

	if len(trimmed) > 1000 {
		return Description{}, errors.DomainError{
			Code:    "POSITION_DESCRIPTION_TOO_LONG",
			Message: "Position description cannot exceed 1000 characters",
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

// Requirements represents position requirements
type Requirements struct {
	value string
}

// NewRequirements creates a new Requirements from a string
func NewRequirements(value string) (Requirements, error) {
	if value == "" {
		return Requirements{}, errors.ErrInvalidPositionRequirements
	}

	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return Requirements{}, errors.ErrInvalidPositionRequirements
	}

	if len(trimmed) > 2000 {
		return Requirements{}, errors.DomainError{
			Code:    "POSITION_REQUIREMENTS_TOO_LONG",
			Message: "Position requirements cannot exceed 2000 characters",
			Field:   "requirements",
		}
	}

	return Requirements{value: trimmed}, nil
}

// String returns the string representation of the Requirements
func (r Requirements) String() string {
	return r.value
}

// Equals checks if two Requirements are equal
func (r Requirements) Equals(other Requirements) bool {
	return r.value == other.value
}

// SalaryRange represents a salary range
type SalaryRange struct {
	minSalary float64
	maxSalary float64
}

// NewSalaryRange creates a new SalaryRange
func NewSalaryRange(minSalary, maxSalary float64) (SalaryRange, error) {
	if minSalary < 0 {
		return SalaryRange{}, errors.DomainError{
			Code:    "INVALID_MIN_SALARY",
			Message: "Minimum salary cannot be negative",
			Field:   "minSalary",
		}
	}

	if maxSalary < 0 {
		return SalaryRange{}, errors.DomainError{
			Code:    "INVALID_MAX_SALARY",
			Message: "Maximum salary cannot be negative",
			Field:   "maxSalary",
		}
	}

	if minSalary > maxSalary {
		return SalaryRange{}, errors.DomainError{
			Code:    "INVALID_SALARY_RANGE",
			Message: "Minimum salary cannot be greater than maximum salary",
			Field:   "salaryRange",
		}
	}

	return SalaryRange{
		minSalary: minSalary,
		maxSalary: maxSalary,
	}, nil
}

// MinSalary returns the minimum salary
func (s SalaryRange) MinSalary() float64 {
	return s.minSalary
}

// MaxSalary returns the maximum salary
func (s SalaryRange) MaxSalary() float64 {
	return s.maxSalary
}

// Equals checks if two SalaryRanges are equal
func (s SalaryRange) Equals(other SalaryRange) bool {
	return s.minSalary == other.minSalary && s.maxSalary == other.maxSalary
}

// DepartmentID represents a department identifier (reused from department domain)
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

// String returns the string representation of the DepartmentID
func (d DepartmentID) String() string {
	return d.value
}

// Equals checks if two DepartmentIDs are equal
func (d DepartmentID) Equals(other DepartmentID) bool {
	return d.value == other.value
}

// ValidatePositionID validates a position ID string
func ValidatePositionID(id string) error {
	if id == "" {
		return errors.ErrInvalidPositionID
	}

	if _, err := uuid.Parse(id); err != nil {
		return errors.ErrInvalidPositionID
	}

	return nil
}

// ValidatePositionTitle validates a position title string
func ValidatePositionTitle(title string) error {
	if title == "" {
		return errors.ErrInvalidPositionTitle
	}

	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return errors.ErrInvalidPositionTitle
	}

	if len(trimmed) > 100 {
		return errors.DomainError{
			Code:    "POSITION_TITLE_TOO_LONG",
			Message: "Position title cannot exceed 100 characters",
			Field:   "title",
		}
	}

	validTitleRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\-_]+$`)
	if !validTitleRegex.MatchString(trimmed) {
		return errors.DomainError{
			Code:    "INVALID_POSITION_TITLE_CHARACTERS",
			Message: "Position title can only contain letters, numbers, spaces, hyphens, and underscores",
			Field:   "title",
		}
	}

	return nil
}

// ValidatePositionDescription validates a position description string
func ValidatePositionDescription(description string) error {
	if description == "" {
		return errors.ErrInvalidPositionDescription
	}

	trimmed := strings.TrimSpace(description)
	if trimmed == "" {
		return errors.ErrInvalidPositionDescription
	}

	if len(trimmed) > 1000 {
		return errors.DomainError{
			Code:    "POSITION_DESCRIPTION_TOO_LONG",
			Message: "Position description cannot exceed 1000 characters",
			Field:   "description",
		}
	}

	return nil
}

// ValidatePositionRequirements validates position requirements string
func ValidatePositionRequirements(requirements string) error {
	if requirements == "" {
		return errors.ErrInvalidPositionRequirements
	}

	trimmed := strings.TrimSpace(requirements)
	if trimmed == "" {
		return errors.ErrInvalidPositionRequirements
	}

	if len(trimmed) > 2000 {
		return errors.DomainError{
			Code:    "POSITION_REQUIREMENTS_TOO_LONG",
			Message: "Position requirements cannot exceed 2000 characters",
			Field:   "requirements",
		}
	}

	return nil
}