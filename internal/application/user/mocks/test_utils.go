package mocks

import (
	"context"
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/application/user"
	"github.com/golang/mock/gomock"
)

// TestDTOBuilder helps create test DTOs with common configurations
type TestDTOBuilder struct {
	id        string
	email     string
	name      string
	createdAt time.Time
	updatedAt time.Time
}

// NewTestDTOBuilder creates a new test DTO builder with default values
func NewTestDTOBuilder() *TestDTOBuilder {
	now := time.Now()
	return &TestDTOBuilder{
		id:        "test-user-id",
		email:     "test@example.com",
		name:      "Test User",
		createdAt: now,
		updatedAt: now,
	}
}

// WithID sets the user ID
func (b *TestDTOBuilder) WithID(id string) *TestDTOBuilder {
	b.id = id
	return b
}

// WithEmail sets the user email
func (b *TestDTOBuilder) WithEmail(email string) *TestDTOBuilder {
	b.email = email
	return b
}

// WithName sets the user name
func (b *TestDTOBuilder) WithName(name string) *TestDTOBuilder {
	b.name = name
	return b
}

// WithCreatedAt sets the created at timestamp
func (b *TestDTOBuilder) WithCreatedAt(createdAt time.Time) *TestDTOBuilder {
	b.createdAt = createdAt
	return b
}

// WithUpdatedAt sets the updated at timestamp
func (b *TestDTOBuilder) WithUpdatedAt(updatedAt time.Time) *TestDTOBuilder {
	b.updatedAt = updatedAt
	return b
}

// BuildUserDTO creates a UserDTO with the configured values
func (b *TestDTOBuilder) BuildUserDTO() *user.UserDTO {
	return &user.UserDTO{
		ID:        b.id,
		Email:     b.email,
		Name:      b.name,
		CreatedAt: b.createdAt,
		UpdatedAt: b.updatedAt,
	}
}

// BuildUserEdgeDTO creates a UserEdgeDTO with the configured values
func (b *TestDTOBuilder) BuildUserEdgeDTO() *user.UserEdgeDTO {
	return &user.UserEdgeDTO{
		Node:   b.BuildUserDTO(),
		Cursor: b.id,
	}
}

// ServiceTestUtils provides common mock service setup utilities
type ServiceTestUtils struct {
	mockService *MockService
}

// NewServiceTestUtils creates a new service test utilities instance
func NewServiceTestUtils(mockService *MockService) *ServiceTestUtils {
	return &ServiceTestUtils{
		mockService: mockService,
	}
}

// ExpectGetUserSuccess sets up mock to return a user successfully
func (u *ServiceTestUtils) ExpectGetUserSuccess(userID string, returnUser *user.UserDTO) {
	req := user.GetUserRequest{ID: userID}
	resp := &user.GetUserResponse{User: returnUser}

	u.mockService.EXPECT().
		GetUser(gomock.Any(), req).
		Return(resp, nil).
		Times(1)
}

// ExpectGetUserNotFound sets up mock to return user not found error
func (u *ServiceTestUtils) ExpectGetUserNotFound(userID string) {
	req := user.GetUserRequest{ID: userID}
	resp := &user.GetUserResponse{
		Errors: []user.ErrorDTO{
			{
				Code:    "USER_NOT_FOUND",
				Message: "User not found",
			},
		},
	}

	u.mockService.EXPECT().
		GetUser(gomock.Any(), req).
		Return(resp, nil).
		Times(1)
}

// ExpectListUsersSuccess sets up mock to return users list successfully
func (u *ServiceTestUtils) ExpectListUsersSuccess(first int, after string, returnUsers *user.UserConnectionDTO) {
	req := user.ListUsersRequest{First: first, After: after}
	resp := &user.ListUsersResponse{Users: returnUsers}

	u.mockService.EXPECT().
		ListUsers(gomock.Any(), req).
		Return(resp, nil).
		Times(1)
}

// ExpectCreateUserSuccess sets up mock to create user successfully
func (u *ServiceTestUtils) ExpectCreateUserSuccess(email, name string, returnUser *user.UserDTO) {
	req := user.CreateUserRequest{Email: email, Name: name}
	resp := &user.CreateUserResponse{User: returnUser}

	u.mockService.EXPECT().
		CreateUser(gomock.Any(), req).
		Return(resp, nil).
		Times(1)
}

// ExpectCreateUserDuplicateEmail sets up mock to return duplicate email error
func (u *ServiceTestUtils) ExpectCreateUserDuplicateEmail(email, name string) {
	req := user.CreateUserRequest{Email: email, Name: name}
	resp := &user.CreateUserResponse{
		Errors: []user.ErrorDTO{
			{
				Code:    "DUPLICATE_EMAIL",
				Message: "Email already exists",
				Field:   "email",
			},
		},
	}

	u.mockService.EXPECT().
		CreateUser(gomock.Any(), req).
		Return(resp, nil).
		Times(1)
}

// ExpectUpdateUserSuccess sets up mock to update user successfully
func (u *ServiceTestUtils) ExpectUpdateUserSuccess(userID string, email, name *string, returnUser *user.UserDTO) {
	req := user.UpdateUserRequest{ID: userID, Email: email, Name: name}
	resp := &user.UpdateUserResponse{User: returnUser}

	u.mockService.EXPECT().
		UpdateUser(gomock.Any(), req).
		Return(resp, nil).
		Times(1)
}

// ExpectUpdateUserNotFound sets up mock to return user not found error for update
func (u *ServiceTestUtils) ExpectUpdateUserNotFound(userID string, email, name *string) {
	req := user.UpdateUserRequest{ID: userID, Email: email, Name: name}
	resp := &user.UpdateUserResponse{
		Errors: []user.ErrorDTO{
			{
				Code:    "USER_NOT_FOUND",
				Message: "User not found",
			},
		},
	}

	u.mockService.EXPECT().
		UpdateUser(gomock.Any(), req).
		Return(resp, nil).
		Times(1)
}

// ExpectDeleteUserSuccess sets up mock to delete user successfully
func (u *ServiceTestUtils) ExpectDeleteUserSuccess(userID string) {
	req := user.DeleteUserRequest{ID: userID}
	resp := &user.DeleteUserResponse{Success: true}

	u.mockService.EXPECT().
		DeleteUser(gomock.Any(), req).
		Return(resp, nil).
		Times(1)
}

// ExpectDeleteUserNotFound sets up mock to return user not found error for deletion
func (u *ServiceTestUtils) ExpectDeleteUserNotFound(userID string) {
	req := user.DeleteUserRequest{ID: userID}
	resp := &user.DeleteUserResponse{
		Success: false,
		Errors: []user.ErrorDTO{
			{
				Code:    "USER_NOT_FOUND",
				Message: "User not found",
			},
		},
	}

	u.mockService.EXPECT().
		DeleteUser(gomock.Any(), req).
		Return(resp, nil).
		Times(1)
}

// BuildTestUserConnection creates a test user connection with the given users
func BuildTestUserConnection(users []*user.UserDTO, hasNextPage, hasPreviousPage bool) *user.UserConnectionDTO {
	edges := make([]*user.UserEdgeDTO, len(users))
	for i, u := range users {
		edges[i] = &user.UserEdgeDTO{
			Node:   u,
			Cursor: u.ID,
		}
	}

	var startCursor, endCursor *string
	if len(edges) > 0 {
		start := edges[0].Cursor
		end := edges[len(edges)-1].Cursor
		startCursor = &start
		endCursor = &end
	}

	return &user.UserConnectionDTO{
		Edges: edges,
		PageInfo: &user.PageInfoDTO{
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
			StartCursor:     startCursor,
			EndCursor:       endCursor,
		},
	}
}

// CommonTestContext provides a common test context with timeout
func CommonTestContext() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_ = cancel // Store cancel function to avoid linter warning, but don't call it since we return the context
	return ctx
}
