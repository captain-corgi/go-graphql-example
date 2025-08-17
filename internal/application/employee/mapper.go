package employee

import (
	"github.com/captain-corgi/go-graphql-example/internal/domain/employee"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
)

// MapEmployeeToDTO converts a domain employee to a DTO
func MapEmployeeToDTO(emp *employee.Employee, user *user.User) *EmployeeDTO {
	if emp == nil {
		return nil
	}

	var userDTO *UserDTO
	if user != nil {
		userDTO = &UserDTO{
			ID:        user.ID().String(),
			Email:     user.Email().String(),
			Name:      user.Name().String(),
			CreatedAt: user.CreatedAt(),
			UpdatedAt: user.UpdatedAt(),
		}
	}

	return &EmployeeDTO{
		ID:           emp.ID().String(),
		UserID:       emp.UserID().String(),
		User:         userDTO,
		EmployeeCode: emp.EmployeeCode().String(),
		Department:   emp.Department().String(),
		Position:     emp.Position().String(),
		HireDate:     emp.HireDate(),
		Salary:       emp.Salary().Value(),
		Status:       emp.Status().String(),
		CreatedAt:    emp.CreatedAt(),
		UpdatedAt:    emp.UpdatedAt(),
	}
}

// MapEmployeesToDTOs converts a slice of domain employees to DTOs
func MapEmployeesToDTOs(employees []*employee.Employee, users map[string]*user.User) []*EmployeeDTO {
	if employees == nil {
		return nil
	}

	dtos := make([]*EmployeeDTO, len(employees))
	for i, emp := range employees {
		var user *user.User
		if users != nil {
			user = users[emp.UserID().String()]
		}
		dtos[i] = MapEmployeeToDTO(emp, user)
	}

	return dtos
}

// MapErrorToDTO converts a domain error to a DTO
func MapErrorToDTO(err error) ErrorDTO {
	if err == nil {
		return ErrorDTO{}
	}

	// Try to cast to domain error
	if domainErr, ok := err.(interface {
		Error() string
		Code() string
		Field() string
	}); ok {
		return ErrorDTO{
			Message: domainErr.Error(),
			Code:    domainErr.Code(),
			Field:   domainErr.Field(),
		}
	}

	// Fallback for generic errors
	return ErrorDTO{
		Message: err.Error(),
		Code:    "UNKNOWN_ERROR",
	}
}

// MapErrorsToDTOs converts a slice of domain errors to DTOs
func MapErrorsToDTOs(errors []error) []ErrorDTO {
	if errors == nil {
		return nil
	}

	dtos := make([]ErrorDTO, len(errors))
	for i, err := range errors {
		dtos[i] = MapErrorToDTO(err)
	}

	return dtos
}