package resolver

import (
	"strings"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/captain-corgi/go-graphql-example/internal/interfaces/graphql/model"
)

// validateUserID validates a user ID input
func validateUserID(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.ErrInvalidUserID
	}
	return nil
}

// validateCreateUserInput validates CreateUserInput
func validateCreateUserInput(input model.CreateUserInput) error {
	// Validate email
	if strings.TrimSpace(input.Email) == "" {
		return errors.ErrInvalidEmail
	}

	// Basic email format validation (more comprehensive validation is done in domain layer)
	if !strings.Contains(input.Email, "@") {
		return errors.ErrInvalidEmail
	}

	// Validate name
	if strings.TrimSpace(input.Name) == "" {
		return errors.ErrInvalidName
	}

	return nil
}

// validateUpdateUserInput validates UpdateUserInput
func validateUpdateUserInput(input model.UpdateUserInput) error {
	// At least one field must be provided for update
	if input.Email == nil && input.Name == nil {
		return errors.DomainError{
			Code:    "NO_UPDATE_FIELDS",
			Message: "At least one field must be provided for update",
		}
	}

	// Validate email if provided
	if input.Email != nil {
		if strings.TrimSpace(*input.Email) == "" {
			return errors.ErrInvalidEmail
		}
		if !strings.Contains(*input.Email, "@") {
			return errors.ErrInvalidEmail
		}
	}

	// Validate name if provided
	if input.Name != nil && strings.TrimSpace(*input.Name) == "" {
		return errors.ErrInvalidName
	}

	return nil
}

// validatePaginationParams validates pagination parameters
func validatePaginationParams(first *int, after *string) error {
	if first != nil {
		if *first < 0 {
			return errors.DomainError{
				Code:    "INVALID_FIRST",
				Message: "First parameter must be non-negative",
				Field:   "first",
			}
		}
		if *first > 100 {
			return errors.DomainError{
				Code:    "INVALID_FIRST",
				Message: "First parameter cannot exceed 100",
				Field:   "first",
			}
		}
	}

	// After parameter validation (cursor should be non-empty if provided)
	if after != nil && strings.TrimSpace(*after) == "" {
		return errors.DomainError{
			Code:    "INVALID_CURSOR",
			Message: "Cursor cannot be empty",
			Field:   "after",
		}
	}

	return nil
}

// sanitizeString trims whitespace and returns the sanitized string
func sanitizeString(s string) string {
	return strings.TrimSpace(s)
}

// sanitizeStringPointer trims whitespace from a string pointer
func sanitizeStringPointer(s *string) *string {
	if s == nil {
		return nil
	}
	sanitized := sanitizeString(*s)
	return &sanitized
}

// validateEmployeeID validates an employee ID input
func validateEmployeeID(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.ErrInvalidEmployeeID
	}
	return nil
}

// validateCreateEmployeeInput validates CreateEmployeeInput
func validateCreateEmployeeInput(input model.CreateEmployeeInput) error {
	// Validate user ID
	if strings.TrimSpace(input.UserID) == "" {
		return errors.DomainError{
			Code:    "INVALID_USER_ID",
			Message: "User ID cannot be empty",
			Field:   "userId",
		}
	}

	// Validate employee code
	if strings.TrimSpace(input.EmployeeCode) == "" {
		return errors.ErrInvalidEmployeeCode
	}

	// Validate department
	if strings.TrimSpace(input.Department) == "" {
		return errors.ErrInvalidDepartment
	}

	// Validate position
	if strings.TrimSpace(input.Position) == "" {
		return errors.ErrInvalidPosition
	}

	// Validate hire date
	if strings.TrimSpace(input.HireDate) == "" {
		return errors.DomainError{
			Code:    "INVALID_HIRE_DATE",
			Message: "Hire date cannot be empty",
			Field:   "hireDate",
		}
	}

	// Validate salary
	if input.Salary <= 0 {
		return errors.ErrInvalidSalary
	}

	// Validate status
	if strings.TrimSpace(input.Status) == "" {
		return errors.ErrInvalidStatus
	}

	return nil
}

// validateUpdateEmployeeInput validates UpdateEmployeeInput
func validateUpdateEmployeeInput(input model.UpdateEmployeeInput) error {
	// At least one field must be provided for update
	if input.EmployeeCode == nil && input.Department == nil && input.Position == nil && 
	   input.HireDate == nil && input.Salary == nil && input.Status == nil {
		return errors.DomainError{
			Code:    "NO_UPDATE_FIELDS",
			Message: "At least one field must be provided for update",
		}
	}

	// Validate employee code if provided
	if input.EmployeeCode != nil && strings.TrimSpace(*input.EmployeeCode) == "" {
		return errors.ErrInvalidEmployeeCode
	}

	// Validate department if provided
	if input.Department != nil && strings.TrimSpace(*input.Department) == "" {
		return errors.ErrInvalidDepartment
	}

	// Validate position if provided
	if input.Position != nil && strings.TrimSpace(*input.Position) == "" {
		return errors.ErrInvalidPosition
	}

	// Validate hire date if provided
	if input.HireDate != nil && strings.TrimSpace(*input.HireDate) == "" {
		return errors.DomainError{
			Code:    "INVALID_HIRE_DATE",
			Message: "Hire date cannot be empty",
			Field:   "hireDate",
		}
	}

	// Validate salary if provided
	if input.Salary != nil && *input.Salary <= 0 {
		return errors.ErrInvalidSalary
	}

	// Validate status if provided
	if input.Status != nil && strings.TrimSpace(*input.Status) == "" {
		return errors.ErrInvalidStatus
	}

	return nil
}
