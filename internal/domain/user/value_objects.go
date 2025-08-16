package user

import (
	"regexp"
	"strings"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/google/uuid"
)

// UserID represents a unique identifier for a user
type UserID struct {
	value string
}

// NewUserID creates a new UserID from a string
func NewUserID(id string) (UserID, error) {
	if id == "" {
		return UserID{}, errors.ErrInvalidUserID
	}

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return UserID{}, errors.ErrInvalidUserID
	}

	return UserID{value: id}, nil
}

// GenerateUserID creates a new random UserID
func GenerateUserID() UserID {
	return UserID{value: uuid.New().String()}
}

// String returns the string representation of the UserID
func (id UserID) String() string {
	return id.value
}

// Equals checks if two UserIDs are equal
func (id UserID) Equals(other UserID) bool {
	return id.value == other.value
}

// Email represents a validated email address
type Email struct {
	value string
}

// emailRegex is a simple email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a new Email value object with validation
func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return Email{}, errors.ErrInvalidEmail
	}

	if !emailRegex.MatchString(email) {
		return Email{}, errors.ErrInvalidEmail
	}

	return Email{value: email}, nil
}

// String returns the string representation of the Email
func (e Email) String() string {
	return e.value
}

// Equals checks if two Emails are equal
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// Name represents a validated user name
type Name struct {
	value string
}

// NewName creates a new Name value object with validation
func NewName(name string) (Name, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return Name{}, errors.ErrInvalidName
	}

	// Additional validation: name should be between 1 and 100 characters
	if len(name) > 100 {
		return Name{}, errors.DomainError{
			Code:    "NAME_TOO_LONG",
			Message: "Name cannot exceed 100 characters",
			Field:   "name",
		}
	}

	return Name{value: name}, nil
}

// String returns the string representation of the Name
func (n Name) String() string {
	return n.value
}

// Equals checks if two Names are equal
func (n Name) Equals(other Name) bool {
	return n.value == other.value
}
