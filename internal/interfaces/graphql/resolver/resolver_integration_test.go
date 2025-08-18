package resolver

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	authMocks "github.com/captain-corgi/go-graphql-example/internal/application/auth/mocks"
	"github.com/captain-corgi/go-graphql-example/internal/application/user"
	"github.com/captain-corgi/go-graphql-example/internal/application/user/mocks"
	"github.com/captain-corgi/go-graphql-example/internal/interfaces/graphql/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestResolverIntegration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mockAuthService := authMocks.NewMockService(ctrl)
	resolver := NewResolver(mockUserService, mockAuthService, logger)

	ctx := context.Background()

	t.Run("User Query - Success", func(t *testing.T) {
		expectedUser := &user.UserDTO{
			ID:        "user-123",
			Email:     "john@example.com",
			Name:      "John Doe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserService.EXPECT().
			GetUser(gomock.Any(), user.GetUserRequest{ID: "user-123"}).
			Return(&user.GetUserResponse{User: expectedUser}, nil)

		result, err := resolver.Query().User(ctx, "user-123")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedUser.ID, result.ID)
		assert.Equal(t, expectedUser.Email, result.Email)
		assert.Equal(t, expectedUser.Name, result.Name)
	})

	t.Run("User Query - Not Found", func(t *testing.T) {
		mockUserService.EXPECT().
			GetUser(gomock.Any(), user.GetUserRequest{ID: "nonexistent"}).
			Return(&user.GetUserResponse{
				Errors: []user.ErrorDTO{{
					Message: "User not found",
					Code:    "USER_NOT_FOUND",
				}},
			}, nil)

		result, err := resolver.Query().User(ctx, "nonexistent")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "User not found")
	})

	t.Run("User Query - Invalid Input", func(t *testing.T) {
		result, err := resolver.Query().User(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Invalid user ID")
	})

	t.Run("Users Query - Success", func(t *testing.T) {
		users := []*user.UserEdgeDTO{
			{
				Node: &user.UserDTO{
					ID:        "user-1",
					Email:     "user1@example.com",
					Name:      "User One",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Cursor: "user-1",
			},
			{
				Node: &user.UserDTO{
					ID:        "user-2",
					Email:     "user2@example.com",
					Name:      "User Two",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Cursor: "user-2",
			},
		}

		connection := &user.UserConnectionDTO{
			Edges: users,
			PageInfo: &user.PageInfoDTO{
				HasNextPage:     false,
				HasPreviousPage: false,
				StartCursor:     &users[0].Cursor,
				EndCursor:       &users[1].Cursor,
			},
		}

		mockUserService.EXPECT().
			ListUsers(gomock.Any(), user.ListUsersRequest{First: 10, After: ""}).
			Return(&user.ListUsersResponse{Users: connection}, nil)

		result, err := resolver.Query().Users(ctx, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Edges, 2)
		assert.Equal(t, "user-1", result.Edges[0].Node.ID)
		assert.Equal(t, "user-2", result.Edges[1].Node.ID)
		assert.False(t, result.PageInfo.HasNextPage)
	})

	t.Run("CreateUser Mutation - Success", func(t *testing.T) {
		input := model.CreateUserInput{
			Email: "new@example.com",
			Name:  "New User",
		}

		expectedUser := &user.UserDTO{
			ID:        "new-user-id",
			Email:     "new@example.com",
			Name:      "New User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserService.EXPECT().
			CreateUser(gomock.Any(), user.CreateUserRequest{
				Email: "new@example.com",
				Name:  "New User",
			}).
			Return(&user.CreateUserResponse{User: expectedUser}, nil)

		result, err := resolver.Mutation().CreateUser(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.User)
		assert.Equal(t, expectedUser.ID, result.User.ID)
		assert.Equal(t, expectedUser.Email, result.User.Email)
		assert.Nil(t, result.Errors)
	})

	t.Run("CreateUser Mutation - Validation Error", func(t *testing.T) {
		input := model.CreateUserInput{
			Email: "", // Invalid empty email
			Name:  "Test User",
		}

		result, err := resolver.Mutation().CreateUser(ctx, input)

		assert.NoError(t, err) // GraphQL mutations return errors in the payload
		assert.NotNil(t, result)
		assert.Nil(t, result.User)
		assert.NotNil(t, result.Errors)
		assert.Len(t, result.Errors, 1)
		assert.Contains(t, result.Errors[0].Message, "Invalid email")
	})

	t.Run("CreateUser Mutation - Duplicate Email", func(t *testing.T) {
		input := model.CreateUserInput{
			Email: "existing@example.com",
			Name:  "Test User",
		}

		mockUserService.EXPECT().
			CreateUser(gomock.Any(), user.CreateUserRequest{
				Email: "existing@example.com",
				Name:  "Test User",
			}).
			Return(&user.CreateUserResponse{
				Errors: []user.ErrorDTO{{
					Message: "Email already exists",
					Code:    "DUPLICATE_EMAIL",
					Field:   "email",
				}},
			}, nil)

		result, err := resolver.Mutation().CreateUser(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Nil(t, result.User)
		assert.NotNil(t, result.Errors)
		assert.Len(t, result.Errors, 1)
		assert.Equal(t, "Email already exists", result.Errors[0].Message)
		assert.Equal(t, "DUPLICATE_EMAIL", *result.Errors[0].Code)
		assert.Equal(t, "email", *result.Errors[0].Field)
	})

	t.Run("UpdateUser Mutation - Success", func(t *testing.T) {
		newEmail := "updated@example.com"
		newName := "Updated Name"
		input := model.UpdateUserInput{
			Email: &newEmail,
			Name:  &newName,
		}

		expectedUser := &user.UserDTO{
			ID:        "user-123",
			Email:     "updated@example.com",
			Name:      "Updated Name",
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now(),
		}

		mockUserService.EXPECT().
			UpdateUser(gomock.Any(), user.UpdateUserRequest{
				ID:    "user-123",
				Email: &newEmail,
				Name:  &newName,
			}).
			Return(&user.UpdateUserResponse{User: expectedUser}, nil)

		result, err := resolver.Mutation().UpdateUser(ctx, "user-123", input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.User)
		assert.Equal(t, expectedUser.ID, result.User.ID)
		assert.Equal(t, expectedUser.Email, result.User.Email)
		assert.Equal(t, expectedUser.Name, result.User.Name)
		assert.Nil(t, result.Errors)
	})

	t.Run("UpdateUser Mutation - Invalid Input", func(t *testing.T) {
		input := model.UpdateUserInput{} // No fields to update

		result, err := resolver.Mutation().UpdateUser(ctx, "user-123", input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Nil(t, result.User)
		assert.NotNil(t, result.Errors)
		assert.Len(t, result.Errors, 1)
		assert.Contains(t, result.Errors[0].Message, "At least one field must be provided")
	})

	t.Run("DeleteUser Mutation - Success", func(t *testing.T) {
		mockUserService.EXPECT().
			DeleteUser(gomock.Any(), user.DeleteUserRequest{ID: "user-123"}).
			Return(&user.DeleteUserResponse{Success: true}, nil)

		result, err := resolver.Mutation().DeleteUser(ctx, "user-123")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Nil(t, result.Errors)
	})

	t.Run("DeleteUser Mutation - User Not Found", func(t *testing.T) {
		mockUserService.EXPECT().
			DeleteUser(gomock.Any(), user.DeleteUserRequest{ID: "nonexistent"}).
			Return(&user.DeleteUserResponse{
				Success: false,
				Errors: []user.ErrorDTO{{
					Message: "User not found",
					Code:    "USER_NOT_FOUND",
				}},
			}, nil)

		result, err := resolver.Mutation().DeleteUser(ctx, "nonexistent")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.Success)
		assert.NotNil(t, result.Errors)
		assert.Len(t, result.Errors, 1)
		assert.Equal(t, "User not found", result.Errors[0].Message)
	})

	t.Run("Input Sanitization", func(t *testing.T) {
		// Test that whitespace is properly trimmed
		expectedUser := &user.UserDTO{
			ID:        "user-123",
			Email:     "test@example.com",
			Name:      "Test User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserService.EXPECT().
			GetUser(gomock.Any(), user.GetUserRequest{ID: "user-123"}).
			Return(&user.GetUserResponse{User: expectedUser}, nil)

		// Input with extra whitespace
		result, err := resolver.Query().User(ctx, "  user-123  ")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedUser.ID, result.ID)
	})

	t.Run("Context Propagation", func(t *testing.T) {
		// Test that context is properly passed through to the service layer
		type testKey string
		const key testKey = "test_key"
		ctxWithValue := context.WithValue(ctx, key, "test_value")

		expectedUser := &user.UserDTO{
			ID:        "user-123",
			Email:     "test@example.com",
			Name:      "Test User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserService.EXPECT().
			GetUser(gomock.Any(), user.GetUserRequest{ID: "user-123"}).
			Do(func(ctx context.Context, req user.GetUserRequest) {
				// Verify that the context value is preserved
				assert.Equal(t, "test_value", ctx.Value(key))
			}).
			Return(&user.GetUserResponse{User: expectedUser}, nil)

		result, err := resolver.Query().User(ctxWithValue, "user-123")

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}
