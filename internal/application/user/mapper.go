package user

import (
	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
)

// mapDomainUserToDTO converts a domain User to a UserDTO
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

// mapDomainErrorToDTO converts a domain error to an ErrorDTO
func mapDomainErrorToDTO(err error) ErrorDTO {
	if err == nil {
		return ErrorDTO{}
	}

	// Check if it's a domain error
	if domainErr, ok := err.(errors.DomainError); ok {
		return ErrorDTO{
			Message: domainErr.Message,
			Field:   domainErr.Field,
			Code:    domainErr.Code,
		}
	}

	// For non-domain errors, return a generic error
	return ErrorDTO{
		Message: "An unexpected error occurred",
		Code:    "INTERNAL_ERROR",
	}
}
