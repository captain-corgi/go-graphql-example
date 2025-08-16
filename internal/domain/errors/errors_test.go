package errors

import (
	"errors"
	"testing"
)

func TestDomainError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      DomainError
		expected string
	}{
		{
			name: "error with field",
			err: DomainError{
				Code:    "INVALID_EMAIL",
				Message: "Invalid email format",
				Field:   "email",
			},
			expected: "email: Invalid email format",
		},
		{
			name: "error without field",
			err: DomainError{
				Code:    "USER_NOT_FOUND",
				Message: "User not found",
				Field:   "",
			},
			expected: "User not found",
		},
		{
			name: "error with empty field",
			err: DomainError{
				Code:    "GENERAL_ERROR",
				Message: "Something went wrong",
				Field:   "",
			},
			expected: "Something went wrong",
		},
		{
			name: "error with whitespace field",
			err: DomainError{
				Code:    "INVALID_INPUT",
				Message: "Invalid input provided",
				Field:   "   ",
			},
			expected: "   : Invalid input provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.expected {
				t.Errorf("DomainError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDomainError_AsError(t *testing.T) {
	domainErr := DomainError{
		Code:    "TEST_ERROR",
		Message: "Test error message",
		Field:   "testField",
	}

	// Test that DomainError implements error interface
	var err error = domainErr
	if err.Error() != "testField: Test error message" {
		t.Errorf("DomainError as error interface = %v, want %v", err.Error(), "testField: Test error message")
	}

	// Test with errors.Is
	if !errors.Is(domainErr, domainErr) {
		t.Error("DomainError should be equal to itself using errors.Is")
	}

	// Test with errors.As
	var targetErr DomainError
	if !errors.As(domainErr, &targetErr) {
		t.Error("DomainError should be assignable using errors.As")
	}

	if targetErr.Code != domainErr.Code {
		t.Errorf("errors.As result Code = %v, want %v", targetErr.Code, domainErr.Code)
	}
}

func TestPredefinedErrors(t *testing.T) {
	// Test that predefined errors have the expected properties
	tests := []struct {
		name         string
		err          DomainError
		expectedCode string
		hasField     bool
	}{
		{"ErrUserNotFound", ErrUserNotFound, "USER_NOT_FOUND", false},
		{"ErrInvalidEmail", ErrInvalidEmail, "INVALID_EMAIL", true},
		{"ErrDuplicateEmail", ErrDuplicateEmail, "DUPLICATE_EMAIL", true},
		{"ErrInvalidName", ErrInvalidName, "INVALID_NAME", true},
		{"ErrInvalidUserID", ErrInvalidUserID, "INVALID_USER_ID", true},
		{"ErrUserAlreadyExists", ErrUserAlreadyExists, "USER_ALREADY_EXISTS", false},
		{"ErrRepositoryConnection", ErrRepositoryConnection, "REPOSITORY_CONNECTION", false},
		{"ErrRepositoryOperation", ErrRepositoryOperation, "REPOSITORY_OPERATION", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.expectedCode {
				t.Errorf("%s.Code = %v, want %v", tt.name, tt.err.Code, tt.expectedCode)
			}

			if tt.err.Message == "" {
				t.Errorf("%s.Message should not be empty", tt.name)
			}

			if tt.hasField && tt.err.Field == "" {
				t.Errorf("%s.Field should not be empty", tt.name)
			}

			if !tt.hasField && tt.err.Field != "" {
				t.Errorf("%s.Field should be empty, got %v", tt.name, tt.err.Field)
			}

			// Test that error message is properly formatted
			errorMsg := tt.err.Error()
			if errorMsg == "" {
				t.Errorf("%s.Error() should not return empty string", tt.name)
			}
		})
	}
}

func TestPredefinedErrors_Uniqueness(t *testing.T) {
	// Collect all error codes to ensure uniqueness
	errorCodes := map[string]string{
		ErrUserNotFound.Code:         "ErrUserNotFound",
		ErrInvalidEmail.Code:         "ErrInvalidEmail",
		ErrDuplicateEmail.Code:       "ErrDuplicateEmail",
		ErrInvalidName.Code:          "ErrInvalidName",
		ErrInvalidUserID.Code:        "ErrInvalidUserID",
		ErrUserAlreadyExists.Code:    "ErrUserAlreadyExists",
		ErrRepositoryConnection.Code: "ErrRepositoryConnection",
		ErrRepositoryOperation.Code:  "ErrRepositoryOperation",
	}

	// Check that we have the expected number of unique codes
	expectedCount := 8
	if len(errorCodes) != expectedCount {
		t.Errorf("Expected %d unique error codes, got %d", expectedCount, len(errorCodes))

		// Print duplicates for debugging
		codes := []DomainError{
			ErrUserNotFound, ErrInvalidEmail, ErrDuplicateEmail, ErrInvalidName,
			ErrInvalidUserID, ErrUserAlreadyExists, ErrRepositoryConnection, ErrRepositoryOperation,
		}

		codeCount := make(map[string]int)
		for _, err := range codes {
			codeCount[err.Code]++
		}

		for code, count := range codeCount {
			if count > 1 {
				t.Errorf("Duplicate error code found: %s (appears %d times)", code, count)
			}
		}
	}
}

func TestDomainError_Comparison(t *testing.T) {
	err1 := DomainError{Code: "TEST_ERROR", Message: "Test message", Field: "field"}
	err2 := DomainError{Code: "TEST_ERROR", Message: "Test message", Field: "field"}
	err3 := DomainError{Code: "DIFFERENT_ERROR", Message: "Test message", Field: "field"}
	err4 := DomainError{Code: "TEST_ERROR", Message: "Different message", Field: "field"}
	err5 := DomainError{Code: "TEST_ERROR", Message: "Test message", Field: "different_field"}

	tests := []struct {
		name     string
		err1     DomainError
		err2     DomainError
		expected bool
	}{
		{
			name:     "identical errors",
			err1:     err1,
			err2:     err2,
			expected: true,
		},
		{
			name:     "different codes",
			err1:     err1,
			err2:     err3,
			expected: false,
		},
		{
			name:     "different messages",
			err1:     err1,
			err2:     err4,
			expected: false,
		},
		{
			name:     "different fields",
			err1:     err1,
			err2:     err5,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := (tt.err1.Code == tt.err2.Code &&
				tt.err1.Message == tt.err2.Message &&
				tt.err1.Field == tt.err2.Field)

			if result != tt.expected {
				t.Errorf("DomainError comparison = %v, want %v", result, tt.expected)
			}
		})
	}
}
