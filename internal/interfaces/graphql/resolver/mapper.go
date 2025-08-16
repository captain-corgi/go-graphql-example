package resolver

import (
	"time"

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
