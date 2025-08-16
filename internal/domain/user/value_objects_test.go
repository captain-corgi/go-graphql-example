package user

import (
	"strings"
	"testing"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/google/uuid"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		want    string
		wantErr error
	}{
		{
			name:    "valid email",
			email:   "test@example.com",
			want:    "test@example.com",
			wantErr: nil,
		},
		{
			name:    "valid email with subdomain",
			email:   "user@mail.example.com",
			want:    "user@mail.example.com",
			wantErr: nil,
		},
		{
			name:    "valid email with numbers",
			email:   "user123@example123.com",
			want:    "user123@example123.com",
			wantErr: nil,
		},
		{
			name:    "valid email with special characters",
			email:   "user.name+tag@example.com",
			want:    "user.name+tag@example.com",
			wantErr: nil,
		},
		{
			name:    "valid email with hyphen in domain",
			email:   "user@my-domain.com",
			want:    "user@my-domain.com",
			wantErr: nil,
		},
		{
			name:    "email with uppercase",
			email:   "Test@Example.COM",
			want:    "test@example.com",
			wantErr: nil,
		},
		{
			name:    "email with spaces",
			email:   "  test@example.com  ",
			want:    "test@example.com",
			wantErr: nil,
		},
		{
			name:    "email with mixed case and spaces",
			email:   "  Test.User@Example.COM  ",
			want:    "test.user@example.com",
			wantErr: nil,
		},
		{
			name:    "invalid email - no @",
			email:   "testexample.com",
			want:    "",
			wantErr: errors.ErrInvalidEmail,
		},
		{
			name:    "invalid email - no domain",
			email:   "test@",
			want:    "",
			wantErr: errors.ErrInvalidEmail,
		},
		{
			name:    "invalid email - no local part",
			email:   "@example.com",
			want:    "",
			wantErr: errors.ErrInvalidEmail,
		},
		{
			name:    "invalid email - multiple @",
			email:   "test@@example.com",
			want:    "",
			wantErr: errors.ErrInvalidEmail,
		},
		{
			name:    "invalid email - no TLD",
			email:   "test@example",
			want:    "",
			wantErr: errors.ErrInvalidEmail,
		},
		{
			name:    "invalid email - TLD too short",
			email:   "test@example.c",
			want:    "",
			wantErr: errors.ErrInvalidEmail,
		},
		{
			name:    "empty email",
			email:   "",
			want:    "",
			wantErr: errors.ErrInvalidEmail,
		},
		{
			name:    "whitespace only email",
			email:   "   ",
			want:    "",
			wantErr: errors.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEmail(tt.email)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("NewEmail() expected error %v, got nil", tt.wantErr)
					return
				}
				if err != tt.wantErr {
					t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewEmail() unexpected error = %v", err)
				return
			}

			if got.String() != tt.want {
				t.Errorf("NewEmail() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}

func TestEmail_Equals(t *testing.T) {
	email1, err := NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to create email1: %v", err)
	}

	email2, err := NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to create email2: %v", err)
	}

	email3, err := NewEmail("different@example.com")
	if err != nil {
		t.Fatalf("Failed to create email3: %v", err)
	}

	// Test case normalization
	email4, err := NewEmail("TEST@EXAMPLE.COM")
	if err != nil {
		t.Fatalf("Failed to create email4: %v", err)
	}

	tests := []struct {
		name     string
		email1   Email
		email2   Email
		expected bool
	}{
		{
			name:     "same email values",
			email1:   email1,
			email2:   email2,
			expected: true,
		},
		{
			name:     "different email values",
			email1:   email1,
			email2:   email3,
			expected: false,
		},
		{
			name:     "case normalized emails are equal",
			email1:   email1,
			email2:   email4,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.email1.Equals(tt.email2)
			if result != tt.expected {
				t.Errorf("Email.Equals() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewName(t *testing.T) {
	tests := []struct {
		name    string
		nameStr string
		want    string
		wantErr bool
		errCode string
	}{
		{
			name:    "valid name",
			nameStr: "John Doe",
			want:    "John Doe",
			wantErr: false,
		},
		{
			name:    "single name",
			nameStr: "John",
			want:    "John",
			wantErr: false,
		},
		{
			name:    "name with middle initial",
			nameStr: "John F. Doe",
			want:    "John F. Doe",
			wantErr: false,
		},
		{
			name:    "name with special characters",
			nameStr: "José María O'Connor-Smith",
			want:    "José María O'Connor-Smith",
			wantErr: false,
		},
		{
			name:    "name with numbers",
			nameStr: "John Doe Jr. III",
			want:    "John Doe Jr. III",
			wantErr: false,
		},
		{
			name:    "name with leading/trailing spaces",
			nameStr: "  John Doe  ",
			want:    "John Doe",
			wantErr: false,
		},
		{
			name:    "name with multiple internal spaces",
			nameStr: "John    Doe",
			want:    "John    Doe",
			wantErr: false,
		},
		{
			name:    "maximum length name (100 chars)",
			nameStr: strings.Repeat("a", 100),
			want:    strings.Repeat("a", 100),
			wantErr: false,
		},
		{
			name:    "empty name",
			nameStr: "",
			want:    "",
			wantErr: true,
			errCode: "INVALID_NAME",
		},
		{
			name:    "name with only spaces",
			nameStr: "   ",
			want:    "",
			wantErr: true,
			errCode: "INVALID_NAME",
		},
		{
			name:    "name with only tabs",
			nameStr: "\t\t\t",
			want:    "",
			wantErr: true,
			errCode: "INVALID_NAME",
		},
		{
			name:    "name with mixed whitespace",
			nameStr: " \t \n ",
			want:    "",
			wantErr: true,
			errCode: "INVALID_NAME",
		},
		{
			name:    "name too long (101 chars)",
			nameStr: strings.Repeat("a", 101),
			want:    "",
			wantErr: true,
			errCode: "NAME_TOO_LONG",
		},
		{
			name:    "very long name",
			nameStr: "This is a very long name that exceeds the maximum allowed length of 100 characters for a user name field",
			want:    "",
			wantErr: true,
			errCode: "NAME_TOO_LONG",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewName(tt.nameStr)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewName() expected error, got nil")
					return
				}

				if domainErr, ok := err.(errors.DomainError); ok {
					if domainErr.Code != tt.errCode {
						t.Errorf("NewName() error code = %v, want %v", domainErr.Code, tt.errCode)
					}
				} else {
					t.Errorf("NewName() expected DomainError, got %T", err)
				}
				return
			}

			if err != nil {
				t.Errorf("NewName() unexpected error = %v", err)
				return
			}

			if got.String() != tt.want {
				t.Errorf("NewName() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}

func TestName_Equals(t *testing.T) {
	name1, err := NewName("John Doe")
	if err != nil {
		t.Fatalf("Failed to create name1: %v", err)
	}

	name2, err := NewName("John Doe")
	if err != nil {
		t.Fatalf("Failed to create name2: %v", err)
	}

	name3, err := NewName("Jane Smith")
	if err != nil {
		t.Fatalf("Failed to create name3: %v", err)
	}

	// Test trimming behavior
	name4, err := NewName("  John Doe  ")
	if err != nil {
		t.Fatalf("Failed to create name4: %v", err)
	}

	tests := []struct {
		name     string
		name1    Name
		name2    Name
		expected bool
	}{
		{
			name:     "same name values",
			name1:    name1,
			name2:    name2,
			expected: true,
		},
		{
			name:     "different name values",
			name1:    name1,
			name2:    name3,
			expected: false,
		},
		{
			name:     "trimmed names are equal",
			name1:    name1,
			name2:    name4,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.name1.Equals(tt.name2)
			if result != tt.expected {
				t.Errorf("Name.Equals() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewUserID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
		errType error
	}{
		{
			name:    "valid UUID v4",
			id:      "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:    "valid UUID v1",
			id:      "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			wantErr: false,
		},
		{
			name:    "valid UUID v5",
			id:      "6ba7b811-9dad-11d1-80b4-00c04fd430c8",
			wantErr: false,
		},
		{
			name:    "valid UUID with uppercase",
			id:      "550E8400-E29B-41D4-A716-446655440000",
			wantErr: false,
		},
		{
			name:    "invalid UUID format - too short",
			id:      "550e8400-e29b-41d4-a716",
			wantErr: true,
			errType: errors.ErrInvalidUserID,
		},
		{
			name:    "invalid UUID format - too long",
			id:      "550e8400-e29b-41d4-a716-446655440000-extra",
			wantErr: true,
			errType: errors.ErrInvalidUserID,
		},
		{
			name:    "invalid UUID format - wrong characters",
			id:      "550g8400-e29b-41d4-a716-446655440000",
			wantErr: true,
			errType: errors.ErrInvalidUserID,
		},
		{
			name:    "valid UUID without hyphens",
			id:      "550e8400e29b41d4a716446655440000",
			wantErr: false,
		},
		{
			name:    "invalid UUID format - wrong hyphen positions",
			id:      "550e84-00e29b-41d4a716-446655440000",
			wantErr: true,
			errType: errors.ErrInvalidUserID,
		},
		{
			name:    "empty ID",
			id:      "",
			wantErr: true,
			errType: errors.ErrInvalidUserID,
		},
		{
			name:    "whitespace only ID",
			id:      "   ",
			wantErr: true,
			errType: errors.ErrInvalidUserID,
		},
		{
			name:    "nil UUID string",
			id:      "00000000-0000-0000-0000-000000000000",
			wantErr: false, // Nil UUID is technically valid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUserID(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewUserID() expected error, got nil")
					return
				}

				if tt.errType != nil && err != tt.errType {
					t.Errorf("NewUserID() error = %v, want %v", err, tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("NewUserID() unexpected error = %v", err)
				return
			}

			if got.String() != tt.id {
				t.Errorf("NewUserID() = %v, want %v", got.String(), tt.id)
			}
		})
	}
}

func TestGenerateUserID(t *testing.T) {
	// Test basic functionality
	id1 := GenerateUserID()
	id2 := GenerateUserID()

	if id1.String() == "" {
		t.Error("GenerateUserID() should not return empty string")
	}

	if id1.Equals(id2) {
		t.Error("GenerateUserID() should generate unique IDs")
	}

	// Test that generated IDs are valid UUIDs
	if _, err := uuid.Parse(id1.String()); err != nil {
		t.Errorf("GenerateUserID() should generate valid UUID, got %v", id1.String())
	}

	if _, err := uuid.Parse(id2.String()); err != nil {
		t.Errorf("GenerateUserID() should generate valid UUID, got %v", id2.String())
	}

	// Test uniqueness over multiple generations
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := GenerateUserID()
		idStr := id.String()

		if ids[idStr] {
			t.Errorf("GenerateUserID() generated duplicate ID: %v", idStr)
		}
		ids[idStr] = true
	}
}

func TestUserID_Equals(t *testing.T) {
	id1, err := NewUserID("550e8400-e29b-41d4-a716-446655440000")
	if err != nil {
		t.Fatalf("Failed to create id1: %v", err)
	}

	id2, err := NewUserID("550e8400-e29b-41d4-a716-446655440000")
	if err != nil {
		t.Fatalf("Failed to create id2: %v", err)
	}

	id3, err := NewUserID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	if err != nil {
		t.Fatalf("Failed to create id3: %v", err)
	}

	// Test case sensitivity
	id4, err := NewUserID("550E8400-E29B-41D4-A716-446655440000")
	if err != nil {
		t.Fatalf("Failed to create id4: %v", err)
	}

	tests := []struct {
		name     string
		id1      UserID
		id2      UserID
		expected bool
	}{
		{
			name:     "same ID values",
			id1:      id1,
			id2:      id2,
			expected: true,
		},
		{
			name:     "different ID values",
			id1:      id1,
			id2:      id3,
			expected: false,
		},
		{
			name:     "case sensitive comparison",
			id1:      id1,
			id2:      id4,
			expected: false, // UUIDs are case-sensitive in our implementation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.id1.Equals(tt.id2)
			if result != tt.expected {
				t.Errorf("UserID.Equals() = %v, want %v", result, tt.expected)
			}
		})
	}
}
