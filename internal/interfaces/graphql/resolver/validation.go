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
