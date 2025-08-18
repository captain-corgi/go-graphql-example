package auth

import (
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
)

// Request DTOs

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email      string  `json:"email"`
	Name       string  `json:"name"`
	Password   string  `json:"password"`
	DeviceInfo *string `json:"deviceInfo,omitempty"`
	IPAddress  *string `json:"ipAddress,omitempty"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	DeviceInfo *string `json:"deviceInfo,omitempty"`
	IPAddress  *string `json:"ipAddress,omitempty"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string  `json:"refreshToken"`
	DeviceInfo   *string `json:"deviceInfo,omitempty"`
	IPAddress    *string `json:"ipAddress,omitempty"`
}

// LogoutRequest represents a user logout request
type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// Response DTOs

// AuthResponse represents an authentication response
type AuthResponse struct {
	User         *UserDTO   `json:"user,omitempty"`
	AccessToken  string     `json:"accessToken,omitempty"`
	RefreshToken string     `json:"refreshToken,omitempty"`
	ExpiresAt    time.Time  `json:"expiresAt,omitempty"`
	Errors       []ErrorDTO `json:"errors,omitempty"`
}

// LogoutResponse represents a logout response
type LogoutResponse struct {
	Success bool       `json:"success"`
	Errors  []ErrorDTO `json:"errors,omitempty"`
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
}

// UserDTO represents a user data transfer object
type UserDTO struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ErrorDTO represents an error data transfer object
type ErrorDTO struct {
	Message string  `json:"message"`
	Field   *string `json:"field,omitempty"`
	Code    *string `json:"code,omitempty"`
}

// Mapper functions

// mapDomainUserToDTO maps a domain user to a UserDTO
func mapDomainUserToDTO(domainUser *user.User) *UserDTO {
	if domainUser == nil {
		return nil
	}

	return &UserDTO{
		ID:        domainUser.ID().String(),
		Email:     domainUser.Email().String(),
		Name:      domainUser.Name().String(),
		CreatedAt: domainUser.CreatedAt(),
		UpdatedAt: domainUser.UpdatedAt(),
	}
}

// mapDomainErrorToDTO maps a domain error to an ErrorDTO
func mapDomainErrorToDTO(err error) ErrorDTO {
	if err == nil {
		return ErrorDTO{}
	}

	// Handle domain errors with structured information
	if domainErr, ok := err.(errors.DomainError); ok {
		return ErrorDTO{
			Message: domainErr.Message,
			Field:   &domainErr.Field,
			Code:    &domainErr.Code,
		}
	}

	// Handle predefined domain errors
	switch err {
	case errors.ErrUserNotFound:
		code := "USER_NOT_FOUND"
		field := "user"
		return ErrorDTO{
			Message: "User not found",
			Field:   &field,
			Code:    &code,
		}
	case errors.ErrInvalidUserID:
		code := "INVALID_USER_ID"
		field := "userId"
		return ErrorDTO{
			Message: "Invalid user ID",
			Field:   &field,
			Code:    &code,
		}
	case errors.ErrInvalidEmail:
		code := "INVALID_EMAIL"
		field := "email"
		return ErrorDTO{
			Message: "Invalid email address",
			Field:   &field,
			Code:    &code,
		}
	case errors.ErrInvalidName:
		code := "INVALID_NAME"
		field := "name"
		return ErrorDTO{
			Message: "Invalid name",
			Field:   &field,
			Code:    &code,
		}
	case errors.ErrDuplicateEmail:
		code := "DUPLICATE_EMAIL"
		field := "email"
		return ErrorDTO{
			Message: "Email address already exists",
			Field:   &field,
			Code:    &code,
		}
	default:
		// Generic error
		code := "INTERNAL_ERROR"
		return ErrorDTO{
			Message: err.Error(),
			Code:    &code,
		}
	}
}

