package auth

import (
	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"golang.org/x/crypto/bcrypt"
)

// PasswordService handles password hashing and verification
type PasswordService struct {
	cost int
}

// NewPasswordService creates a new password service
func NewPasswordService() *PasswordService {
	return &PasswordService{
		cost: bcrypt.DefaultCost,
	}
}

// HashPassword hashes a password using bcrypt
func (p *PasswordService) HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.DomainError{
			Code:    "INVALID_PASSWORD",
			Message: "Password cannot be empty",
			Field:   "password",
		}
	}

	if len(password) < 8 {
		return "", errors.DomainError{
			Code:    "PASSWORD_TOO_SHORT",
			Message: "Password must be at least 8 characters long",
			Field:   "password",
		}
	}

	if len(password) > 128 {
		return "", errors.DomainError{
			Code:    "PASSWORD_TOO_LONG",
			Message: "Password cannot exceed 128 characters",
			Field:   "password",
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		return "", errors.DomainError{
			Code:    "PASSWORD_HASH_FAILED",
			Message: "Failed to hash password",
			Field:   "password",
		}
	}

	return string(hash), nil
}

// VerifyPassword verifies a password against its hash
func (p *PasswordService) VerifyPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.DomainError{
			Code:    "INVALID_CREDENTIALS",
			Message: "Invalid email or password",
			Field:   "password",
		}
	}
	return nil
}

