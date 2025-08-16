package user

import (
	"strings"
	"testing"
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		userName string
		wantErr  error
	}{
		{
			name:     "valid user",
			email:    "test@example.com",
			userName: "John Doe",
			wantErr:  nil,
		},
		{
			name:     "valid user with special characters in name",
			email:    "test@example.com",
			userName: "José María O'Connor-Smith",
			wantErr:  nil,
		},
		{
			name:     "valid user with numbers in name",
			email:    "test@example.com",
			userName: "John Doe Jr. III",
			wantErr:  nil,
		},
		{
			name:     "invalid email",
			email:    "invalid-email",
			userName: "John Doe",
			wantErr:  errors.ErrInvalidEmail,
		},
		{
			name:     "empty name",
			email:    "test@example.com",
			userName: "",
			wantErr:  errors.ErrInvalidName,
		},
		{
			name:     "empty email",
			email:    "",
			userName: "John Doe",
			wantErr:  errors.ErrInvalidEmail,
		},
		{
			name:     "name too long",
			email:    "test@example.com",
			userName: strings.Repeat("a", 101),
			wantErr:  errors.DomainError{Code: "NAME_TOO_LONG", Message: "Name cannot exceed 100 characters", Field: "name"},
		},
		{
			name:     "whitespace only name",
			email:    "test@example.com",
			userName: "   ",
			wantErr:  errors.ErrInvalidName,
		},
		{
			name:     "email with uppercase gets normalized",
			email:    "Test@Example.COM",
			userName: "John Doe",
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.email, tt.userName)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("NewUser() expected error %v, got nil", tt.wantErr)
					return
				}
				// For domain errors, check the code
				if domainErr, ok := tt.wantErr.(errors.DomainError); ok {
					if actualErr, ok := err.(errors.DomainError); ok {
						if actualErr.Code != domainErr.Code {
							t.Errorf("NewUser() error code = %v, wantErr code %v", actualErr.Code, domainErr.Code)
						}
					} else {
						t.Errorf("NewUser() expected DomainError, got %T", err)
					}
				} else if err != tt.wantErr {
					t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewUser() unexpected error = %v", err)
				return
			}

			if user == nil {
				t.Error("NewUser() returned nil user")
				return
			}

			// Verify user properties
			expectedEmail := strings.ToLower(strings.TrimSpace(tt.email))
			if user.Email().String() != expectedEmail {
				t.Errorf("NewUser() email = %v, want %v", user.Email().String(), expectedEmail)
			}

			expectedName := strings.TrimSpace(tt.userName)
			if user.Name().String() != expectedName {
				t.Errorf("NewUser() name = %v, want %v", user.Name().String(), expectedName)
			}

			if user.ID().String() == "" {
				t.Error("NewUser() ID should not be empty")
			}

			// Verify ID is a valid UUID
			if _, err := uuid.Parse(user.ID().String()); err != nil {
				t.Errorf("NewUser() ID should be a valid UUID, got %v", user.ID().String())
			}

			if user.CreatedAt().IsZero() {
				t.Error("NewUser() CreatedAt should not be zero")
			}

			if user.UpdatedAt().IsZero() {
				t.Error("NewUser() UpdatedAt should not be zero")
			}

			// CreatedAt and UpdatedAt should be equal for new users
			if !user.CreatedAt().Equal(user.UpdatedAt()) {
				t.Error("NewUser() CreatedAt and UpdatedAt should be equal for new users")
			}

			// Verify the user is valid
			if err := user.Validate(); err != nil {
				t.Errorf("NewUser() created invalid user: %v", err)
			}
		})
	}
}

func TestNewUserWithID(t *testing.T) {
	validID := uuid.New().String()
	createdAt := time.Now().Add(-time.Hour)
	updatedAt := time.Now()

	tests := []struct {
		name      string
		id        string
		email     string
		userName  string
		createdAt time.Time
		updatedAt time.Time
		wantErr   bool
	}{
		{
			name:      "valid user with ID",
			id:        validID,
			email:     "test@example.com",
			userName:  "John Doe",
			createdAt: createdAt,
			updatedAt: updatedAt,
			wantErr:   false,
		},
		{
			name:      "invalid ID format",
			id:        "invalid-id",
			email:     "test@example.com",
			userName:  "John Doe",
			createdAt: createdAt,
			updatedAt: updatedAt,
			wantErr:   true,
		},
		{
			name:      "empty ID",
			id:        "",
			email:     "test@example.com",
			userName:  "John Doe",
			createdAt: createdAt,
			updatedAt: updatedAt,
			wantErr:   true,
		},
		{
			name:      "invalid email",
			id:        validID,
			email:     "invalid-email",
			userName:  "John Doe",
			createdAt: createdAt,
			updatedAt: updatedAt,
			wantErr:   true,
		},
		{
			name:      "invalid name",
			id:        validID,
			email:     "test@example.com",
			userName:  "",
			createdAt: createdAt,
			updatedAt: updatedAt,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUserWithID(tt.id, tt.email, tt.userName, tt.createdAt, tt.updatedAt)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewUserWithID() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("NewUserWithID() unexpected error = %v", err)
				return
			}

			if user == nil {
				t.Error("NewUserWithID() returned nil user")
				return
			}

			if user.ID().String() != tt.id {
				t.Errorf("NewUserWithID() ID = %v, want %v", user.ID().String(), tt.id)
			}

			if user.CreatedAt() != tt.createdAt {
				t.Errorf("NewUserWithID() CreatedAt = %v, want %v", user.CreatedAt(), tt.createdAt)
			}

			if user.UpdatedAt() != tt.updatedAt {
				t.Errorf("NewUserWithID() UpdatedAt = %v, want %v", user.UpdatedAt(), tt.updatedAt)
			}
		})
	}
}

func TestUser_UpdateEmail(t *testing.T) {
	tests := []struct {
		name     string
		newEmail string
		wantErr  bool
	}{
		{
			name:     "valid email update",
			newEmail: "newemail@example.com",
			wantErr:  false,
		},
		{
			name:     "email with uppercase gets normalized",
			newEmail: "UPPER@EXAMPLE.COM",
			wantErr:  false,
		},
		{
			name:     "email with spaces gets trimmed",
			newEmail: "  spaced@example.com  ",
			wantErr:  false,
		},
		{
			name:     "invalid email format",
			newEmail: "invalid-email",
			wantErr:  true,
		},
		{
			name:     "empty email",
			newEmail: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh user for each test
			testUser, err := NewUser("original@example.com", "John Doe")
			if err != nil {
				t.Fatalf("Failed to create test user: %v", err)
			}

			originalUpdatedAt := testUser.UpdatedAt()
			time.Sleep(1 * time.Millisecond) // Ensure time difference

			err = testUser.UpdateEmail(tt.newEmail)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateEmail() expected error, got nil")
				}
				// Verify email wasn't changed on error
				if testUser.Email().String() != "original@example.com" {
					t.Errorf("UpdateEmail() should not change email on error")
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateEmail() unexpected error = %v", err)
				return
			}

			expectedEmail := strings.ToLower(strings.TrimSpace(tt.newEmail))
			if testUser.Email().String() != expectedEmail {
				t.Errorf("UpdateEmail() email = %v, want %v", testUser.Email().String(), expectedEmail)
			}

			if !testUser.UpdatedAt().After(originalUpdatedAt) {
				t.Error("UpdateEmail() should update the UpdatedAt timestamp")
			}
		})
	}
}

func TestUser_UpdateName(t *testing.T) {
	tests := []struct {
		name    string
		newName string
		wantErr bool
		errCode string
	}{
		{
			name:    "valid name update",
			newName: "Jane Smith",
			wantErr: false,
		},
		{
			name:    "name with special characters",
			newName: "José María O'Connor-Smith",
			wantErr: false,
		},
		{
			name:    "name with spaces gets trimmed",
			newName: "  Trimmed Name  ",
			wantErr: false,
		},
		{
			name:    "empty name",
			newName: "",
			wantErr: true,
			errCode: "INVALID_NAME",
		},
		{
			name:    "whitespace only name",
			newName: "   ",
			wantErr: true,
			errCode: "INVALID_NAME",
		},
		{
			name:    "name too long",
			newName: strings.Repeat("a", 101),
			wantErr: true,
			errCode: "NAME_TOO_LONG",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh user for each test
			testUser, err := NewUser("test@example.com", "Original Name")
			if err != nil {
				t.Fatalf("Failed to create test user: %v", err)
			}

			originalUpdatedAt := testUser.UpdatedAt()
			time.Sleep(1 * time.Millisecond) // Ensure time difference

			err = testUser.UpdateName(tt.newName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateName() expected error, got nil")
					return
				}

				if domainErr, ok := err.(errors.DomainError); ok {
					if domainErr.Code != tt.errCode {
						t.Errorf("UpdateName() error code = %v, want %v", domainErr.Code, tt.errCode)
					}
				} else {
					t.Errorf("UpdateName() expected DomainError, got %T", err)
				}

				// Verify name wasn't changed on error
				if testUser.Name().String() != "Original Name" {
					t.Errorf("UpdateName() should not change name on error")
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateName() unexpected error = %v", err)
				return
			}

			expectedName := strings.TrimSpace(tt.newName)
			if testUser.Name().String() != expectedName {
				t.Errorf("UpdateName() name = %v, want %v", testUser.Name().String(), expectedName)
			}

			if !testUser.UpdatedAt().After(originalUpdatedAt) {
				t.Error("UpdateName() should update the UpdatedAt timestamp")
			}
		})
	}
}

func TestUser_Validate(t *testing.T) {
	validUser, err := NewUser("test@example.com", "John Doe")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create users with various invalid states for testing
	userWithEmptyID := &User{
		id:        UserID{value: ""},
		email:     validUser.email,
		name:      validUser.name,
		createdAt: validUser.createdAt,
		updatedAt: validUser.updatedAt,
	}

	userWithEmptyEmail := &User{
		id:        validUser.id,
		email:     Email{value: ""},
		name:      validUser.name,
		createdAt: validUser.createdAt,
		updatedAt: validUser.updatedAt,
	}

	userWithEmptyName := &User{
		id:        validUser.id,
		email:     validUser.email,
		name:      Name{value: ""},
		createdAt: validUser.createdAt,
		updatedAt: validUser.updatedAt,
	}

	userWithZeroCreatedAt := &User{
		id:        validUser.id,
		email:     validUser.email,
		name:      validUser.name,
		createdAt: time.Time{},
		updatedAt: validUser.updatedAt,
	}

	userWithZeroUpdatedAt := &User{
		id:        validUser.id,
		email:     validUser.email,
		name:      validUser.name,
		createdAt: validUser.createdAt,
		updatedAt: time.Time{},
	}

	userWithInvalidTimestamps := &User{
		id:        validUser.id,
		email:     validUser.email,
		name:      validUser.name,
		createdAt: time.Now(),
		updatedAt: time.Now().Add(-time.Hour), // UpdatedAt before CreatedAt
	}

	tests := []struct {
		name    string
		user    *User
		wantErr bool
		errCode string
	}{
		{
			name:    "valid user",
			user:    validUser,
			wantErr: false,
		},
		{
			name:    "invalid user ID",
			user:    userWithEmptyID,
			wantErr: true,
			errCode: "INVALID_USER_ID",
		},
		{
			name:    "invalid email",
			user:    userWithEmptyEmail,
			wantErr: true,
			errCode: "INVALID_EMAIL",
		},
		{
			name:    "invalid name",
			user:    userWithEmptyName,
			wantErr: true,
			errCode: "INVALID_NAME",
		},
		{
			name:    "zero created at",
			user:    userWithZeroCreatedAt,
			wantErr: true,
			errCode: "INVALID_CREATED_AT",
		},
		{
			name:    "zero updated at",
			user:    userWithZeroUpdatedAt,
			wantErr: true,
			errCode: "INVALID_UPDATED_AT",
		},
		{
			name:    "invalid timestamps - updated before created",
			user:    userWithInvalidTimestamps,
			wantErr: true,
			errCode: "INVALID_TIMESTAMPS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error, got nil")
					return
				}

				if domainErr, ok := err.(errors.DomainError); ok {
					if domainErr.Code != tt.errCode {
						t.Errorf("Validate() error code = %v, want %v", domainErr.Code, tt.errCode)
					}
				} else {
					t.Errorf("Validate() expected DomainError, got %T", err)
				}
				return
			}

			if err != nil {
				t.Errorf("Validate() unexpected error = %v", err)
			}
		})
	}
}

func TestUser_Equals(t *testing.T) {
	user1, err := NewUser("test1@example.com", "User One")
	if err != nil {
		t.Fatalf("Failed to create user1: %v", err)
	}

	user2, err := NewUser("test2@example.com", "User Two")
	if err != nil {
		t.Fatalf("Failed to create user2: %v", err)
	}

	// Create user3 with same ID as user1
	user3, err := NewUserWithID(user1.ID().String(), "test3@example.com", "User Three", time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to create user3: %v", err)
	}

	tests := []struct {
		name     string
		user     *User
		other    *User
		expected bool
	}{
		{
			name:     "same user instance",
			user:     user1,
			other:    user1,
			expected: true,
		},
		{
			name:     "different users",
			user:     user1,
			other:    user2,
			expected: false,
		},
		{
			name:     "users with same ID",
			user:     user1,
			other:    user3,
			expected: true,
		},
		{
			name:     "user compared to nil",
			user:     user1,
			other:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.Equals(tt.other)
			if result != tt.expected {
				t.Errorf("User.Equals() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUser_Getters(t *testing.T) {
	email := "test@example.com"
	name := "John Doe"
	user, err := NewUser(email, name)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test ID getter
	if user.ID().String() == "" {
		t.Error("ID() should return non-empty ID")
	}

	// Test Email getter
	if user.Email().String() != email {
		t.Errorf("Email() = %v, want %v", user.Email().String(), email)
	}

	// Test Name getter
	if user.Name().String() != name {
		t.Errorf("Name() = %v, want %v", user.Name().String(), name)
	}

	// Test CreatedAt getter
	if user.CreatedAt().IsZero() {
		t.Error("CreatedAt() should return non-zero time")
	}

	// Test UpdatedAt getter
	if user.UpdatedAt().IsZero() {
		t.Error("UpdatedAt() should return non-zero time")
	}

	// Test that getters return copies/immutable values
	originalID := user.ID()
	returnedID := user.ID()
	if !originalID.Equals(returnedID) {
		t.Error("ID() should return consistent values")
	}
}
