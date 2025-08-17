package resolver

import (
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/application/employee"
	"github.com/captain-corgi/go-graphql-example/internal/application/user"
	"github.com/captain-corgi/go-graphql-example/internal/interfaces/graphql/model"
)

// mapUserDTOToGraphQL converts a UserDTO to a GraphQL User model
func mapUserDTOToGraphQL(dto *user.UserDTO) *model.User {
	if dto == nil {
		return nil
	}

	return &model.User{
		ID:        dto.ID,
		Email:     dto.Email,
		Name:      dto.Name,
		CreatedAt: dto.CreatedAt.Format(time.RFC3339),
		UpdatedAt: dto.UpdatedAt.Format(time.RFC3339),
	}
}

// mapUserConnectionDTOToGraphQL converts a UserConnectionDTO to a GraphQL UserConnection model
func mapUserConnectionDTOToGraphQL(dto *user.UserConnectionDTO) *model.UserConnection {
	if dto == nil {
		return nil
	}

	edges := make([]*model.UserEdge, len(dto.Edges))
	for i, edge := range dto.Edges {
		edges[i] = &model.UserEdge{
			Node:   mapUserDTOToGraphQL(edge.Node),
			Cursor: edge.Cursor,
		}
	}

	return &model.UserConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     dto.PageInfo.HasNextPage,
			HasPreviousPage: dto.PageInfo.HasPreviousPage,
			StartCursor:     dto.PageInfo.StartCursor,
			EndCursor:       dto.PageInfo.EndCursor,
		},
	}
}

// mapErrorDTOsToGraphQL converts ErrorDTOs to GraphQL Error models
func mapErrorDTOsToGraphQL(dtos []user.ErrorDTO) []*model.Error {
	if len(dtos) == 0 {
		return nil
	}

	errors := make([]*model.Error, len(dtos))
	for i, dto := range dtos {
		errors[i] = mapErrorDTOToGraphQL(dto)
	}

	return errors
}

// mapErrorDTOToGraphQL converts an ErrorDTO to a GraphQL Error model
func mapErrorDTOToGraphQL(dto user.ErrorDTO) *model.Error {
	var field *string
	if dto.Field != "" {
		field = &dto.Field
	}

	var code *string
	if dto.Code != "" {
		code = &dto.Code
	}

	return &model.Error{
		Message: dto.Message,
		Field:   field,
		Code:    code,
	}
}

// mapCreateUserInputToRequest converts GraphQL CreateUserInput to application request
func mapCreateUserInputToRequest(input model.CreateUserInput) user.CreateUserRequest {
	return user.CreateUserRequest{
		Email: input.Email,
		Name:  input.Name,
	}
}

// mapUpdateUserInputToRequest converts GraphQL UpdateUserInput to application request
func mapUpdateUserInputToRequest(id string, input model.UpdateUserInput) user.UpdateUserRequest {
	return user.UpdateUserRequest{
		ID:    id,
		Email: input.Email,
		Name:  input.Name,
	}
}

// mapEmployeeDTOToGraphQL converts an EmployeeDTO to a GraphQL Employee model
func mapEmployeeDTOToGraphQL(dto *employee.EmployeeDTO) *model.Employee {
	if dto == nil {
		return nil
	}

	var user *model.User
	if dto.User != nil {
		user = &model.User{
			ID:        dto.User.ID,
			Email:     dto.User.Email,
			Name:      dto.User.Name,
			CreatedAt: dto.User.CreatedAt.Format(time.RFC3339),
			UpdatedAt: dto.User.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &model.Employee{
		ID:           dto.ID,
		User:         user,
		EmployeeCode: dto.EmployeeCode,
		Department:   dto.Department,
		Position:     dto.Position,
		HireDate:     dto.HireDate.Format("2006-01-02"),
		Salary:       dto.Salary,
		Status:       dto.Status,
		CreatedAt:    dto.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    dto.UpdatedAt.Format(time.RFC3339),
	}
}

// mapEmployeeConnectionDTOToGraphQL converts EmployeeDTOs to a GraphQL EmployeeConnection model
func mapEmployeeConnectionDTOToGraphQL(dtos []*employee.EmployeeDTO, nextCursor string) *model.EmployeeConnection {
	if len(dtos) == 0 {
		return &model.EmployeeConnection{
			Edges: []*model.EmployeeEdge{},
			PageInfo: &model.PageInfo{
				HasNextPage:     false,
				HasPreviousPage: false,
			},
		}
	}

	edges := make([]*model.EmployeeEdge, len(dtos))
	for i, dto := range dtos {
		// Generate cursor for pagination
		cursor := dto.CreatedAt.Format(time.RFC3339) + ":" + dto.ID
		edges[i] = &model.EmployeeEdge{
			Node:   mapEmployeeDTOToGraphQL(dto),
			Cursor: cursor,
		}
	}

	// Determine if there are more pages
	hasNextPage := nextCursor != ""

	return &model.EmployeeConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: false, // We don't support backward pagination for now
			StartCursor:     &edges[0].Cursor,
			EndCursor:       &edges[len(edges)-1].Cursor,
		},
	}
}

// mapEmployeeErrorDTOsToGraphQL converts Employee ErrorDTOs to GraphQL Error models
func mapEmployeeErrorDTOsToGraphQL(dtos []employee.ErrorDTO) []*model.Error {
	if len(dtos) == 0 {
		return nil
	}

	errors := make([]*model.Error, len(dtos))
	for i, dto := range dtos {
		errors[i] = mapEmployeeErrorDTOToGraphQL(dto)
	}

	return errors
}

// mapEmployeeErrorDTOToGraphQL converts an Employee ErrorDTO to a GraphQL Error model
func mapEmployeeErrorDTOToGraphQL(dto employee.ErrorDTO) *model.Error {
	var field *string
	if dto.Field != "" {
		field = &dto.Field
	}

	var code *string
	if dto.Code != "" {
		code = &dto.Code
	}

	return &model.Error{
		Message: dto.Message,
		Field:   field,
		Code:    code,
	}
}
