package user

import (
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
)

// User represents the user domain entity
type User struct {
	id           UserID
	email        Email
	name         Name
	passwordHash HashedPassword
	isActive     bool
	lastLoginAt  *time.Time
	createdAt    time.Time
	updatedAt    time.Time
}

// NewUser creates a new User entity with validation
func NewUser(email, name string) (*User, error) {
	emailVO, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	nameVO, err := NewName(name)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &User{
		id:        GenerateUserID(),
		email:     emailVO,
		name:      nameVO,
		isActive:  true,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// NewUserWithPassword creates a new User entity with password
func NewUserWithPassword(email, name, passwordHash string) (*User, error) {
	emailVO, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	nameVO, err := NewName(name)
	if err != nil {
		return nil, err
	}

	passwordHashVO, err := NewHashedPassword(passwordHash)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &User{
		id:           GenerateUserID(),
		email:        emailVO,
		name:         nameVO,
		passwordHash: passwordHashVO,
		isActive:     true,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// NewUserWithID creates a User entity with a specific ID (for reconstruction from persistence)
func NewUserWithID(id, email, name string, createdAt, updatedAt time.Time) (*User, error) {
	userID, err := NewUserID(id)
	if err != nil {
		return nil, err
	}

	emailVO, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	nameVO, err := NewName(name)
	if err != nil {
		return nil, err
	}

	return &User{
		id:        userID,
		email:     emailVO,
		name:      nameVO,
		isActive:  true,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

// NewUserWithFullDetails creates a User entity with all details (for reconstruction from persistence)
func NewUserWithFullDetails(id, email, name, passwordHash string, isActive bool, lastLoginAt *time.Time, createdAt, updatedAt time.Time) (*User, error) {
	userID, err := NewUserID(id)
	if err != nil {
		return nil, err
	}

	emailVO, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	nameVO, err := NewName(name)
	if err != nil {
		return nil, err
	}

	passwordHashVO, err := NewHashedPassword(passwordHash)
	if err != nil {
		return nil, err
	}

	return &User{
		id:           userID,
		email:        emailVO,
		name:         nameVO,
		passwordHash: passwordHashVO,
		isActive:     isActive,
		lastLoginAt:  lastLoginAt,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}, nil
}

// ID returns the user's ID
func (u *User) ID() UserID {
	return u.id
}

// Email returns the user's email
func (u *User) Email() Email {
	return u.email
}

// Name returns the user's name
func (u *User) Name() Name {
	return u.name
}

// CreatedAt returns when the user was created
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns when the user was last updated
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// PasswordHash returns the user's password hash
func (u *User) PasswordHash() HashedPassword {
	return u.passwordHash
}

// IsActive returns whether the user is active
func (u *User) IsActive() bool {
	return u.isActive
}

// LastLoginAt returns when the user last logged in
func (u *User) LastLoginAt() *time.Time {
	return u.lastLoginAt
}

// UpdateEmail updates the user's email with validation
func (u *User) UpdateEmail(email string) error {
	emailVO, err := NewEmail(email)
	if err != nil {
		return err
	}

	u.email = emailVO
	u.updatedAt = time.Now()
	return nil
}

// UpdateName updates the user's name with validation
func (u *User) UpdateName(name string) error {
	nameVO, err := NewName(name)
	if err != nil {
		return err
	}

	u.name = nameVO
	u.updatedAt = time.Now()
	return nil
}

// UpdatePassword updates the user's password hash
func (u *User) UpdatePassword(passwordHash string) error {
	passwordHashVO, err := NewHashedPassword(passwordHash)
	if err != nil {
		return err
	}

	u.passwordHash = passwordHashVO
	u.updatedAt = time.Now()
	return nil
}

// RecordLogin records a successful login
func (u *User) RecordLogin() {
	now := time.Now()
	u.lastLoginAt = &now
	u.updatedAt = now
}

// Deactivate deactivates the user account
func (u *User) Deactivate() {
	u.isActive = false
	u.updatedAt = time.Now()
}

// Activate activates the user account
func (u *User) Activate() {
	u.isActive = true
	u.updatedAt = time.Now()
}

// Validate performs comprehensive validation of the user entity
func (u *User) Validate() error {
	if u.id.String() == "" {
		return errors.ErrInvalidUserID
	}

	if u.email.String() == "" {
		return errors.ErrInvalidEmail
	}

	if u.name.String() == "" {
		return errors.ErrInvalidName
	}

	if u.createdAt.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_CREATED_AT",
			Message: "Created at timestamp cannot be zero",
			Field:   "createdAt",
		}
	}

	if u.updatedAt.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_UPDATED_AT",
			Message: "Updated at timestamp cannot be zero",
			Field:   "updatedAt",
		}
	}

	if u.updatedAt.Before(u.createdAt) {
		return errors.DomainError{
			Code:    "INVALID_TIMESTAMPS",
			Message: "Updated at cannot be before created at",
			Field:   "updatedAt",
		}
	}

	return nil
}

// Equals checks if two users are equal based on their ID
func (u *User) Equals(other *User) bool {
	if other == nil {
		return false
	}
	return u.id.Equals(other.id)
}
