package department

import (
	"github.com/captain-corgi/go-graphql-example/internal/domain/department"
	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
)

// mapDomainDepartmentToDTO converts a domain department to a DTO
func mapDomainDepartmentToDTO(domainDept *department.Department) *DepartmentDTO {
	if domainDept == nil {
		return nil
	}

	dto := &DepartmentDTO{
		ID:          domainDept.ID().String(),
		Name:        domainDept.Name().String(),
		Description: domainDept.Description().String(),
		CreatedAt:   domainDept.CreatedAt(),
		UpdatedAt:   domainDept.UpdatedAt(),
	}

	if domainDept.ManagerID() != nil {
		managerID := domainDept.ManagerID().String()
		dto.ManagerID = &managerID
	}

	return dto
}

// mapDomainDepartmentsToDTOs converts a slice of domain departments to DTOs
func mapDomainDepartmentsToDTOs(domainDepts []*department.Department) []*DepartmentDTO {
	if domainDepts == nil {
		return nil
	}

	dtos := make([]*DepartmentDTO, len(domainDepts))
	for i, dept := range domainDepts {
		dtos[i] = mapDomainDepartmentToDTO(dept)
	}

	return dtos
}

// mapDomainErrorToDTO converts a domain error to a DTO
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

	// For other errors, return a generic error
	return ErrorDTO{
		Message: err.Error(),
		Code:    "INTERNAL_ERROR",
	}
}

// mapDomainErrorsToDTOs converts a slice of domain errors to DTOs
func mapDomainErrorsToDTOs(errs []error) []ErrorDTO {
	if errs == nil {
		return nil
	}

	dtos := make([]ErrorDTO, len(errs))
	for i, err := range errs {
		dtos[i] = mapDomainErrorToDTO(err)
	}

	return dtos
}