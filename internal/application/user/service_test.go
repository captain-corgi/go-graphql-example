package user

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user/mocks"
)

func TestService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	logger := slog.Default()
	service := NewService(mockRepo, logger)

	tests := []struct {
		name    string
		request CreateUserRequest
		setup   func()
		want    *CreateUserResponse
		wantErr bool
	}{
		{
			name: "successful user creation",
			request: CreateUserRequest{
				Email: "test@example.com",
				Name:  "Test User",
			},
			setup: func() {
				email, _ := user.NewEmail("test@example.com")
				mockRepo.EXPECT().ExistsByEmail(gomock.Any(), email).Return(false, nil)
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			want: &CreateUserResponse{
				User: &UserDTO{
					Email: "test@example.com",
					Name:  "Test User",
				},
			},
			wantErr: false,
		},
		{
			name: "duplicate email error",
			request: CreateUserRequest{
				Email: "existing@example.com",
				Name:  "Test User",
			},
			setup: func() {
				email, _ := user.NewEmail("existing@example.com")
				mockRepo.EXPECT().ExistsByEmail(gomock.Any(), email).Return(true, nil)
			},
			want: &CreateUserResponse{
				Errors: []ErrorDTO{
					{
						Message: "Email already exists",
						Field:   "email",
						Code:    "DUPLICATE_EMAIL",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid email format",
			request: CreateUserRequest{
				Email: "invalid-email",
				Name:  "Test User",
			},
			setup: func() {
				// No mock setup needed as validation happens before repository call
			},
			want: &CreateUserResponse{
				Errors: []ErrorDTO{
					{
						Message: "Invalid email format",
						Field:   "email",
						Code:    "INVALID_EMAIL",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty name",
			request: CreateUserRequest{
				Email: "test@example.com",
				Name:  "",
			},
			setup: func() {
				// No mock setup needed as validation happens before repository call
			},
			want: &CreateUserResponse{
				Errors: []ErrorDTO{
					{
						Message: "Name cannot be empty",
						Field:   "name",
						Code:    "INVALID_NAME",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := service.CreateUser(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)

			if tt.want.User != nil {
				require.NotNil(t, got.User)
				assert.Equal(t, tt.want.User.Email, got.User.Email)
				assert.Equal(t, tt.want.User.Name, got.User.Name)
				assert.NotEmpty(t, got.User.ID)
				assert.False(t, got.User.CreatedAt.IsZero())
				assert.False(t, got.User.UpdatedAt.IsZero())
			}

			if len(tt.want.Errors) > 0 {
				require.Len(t, got.Errors, len(tt.want.Errors))
				for i, expectedErr := range tt.want.Errors {
					assert.Equal(t, expectedErr.Code, got.Errors[i].Code)
					assert.Equal(t, expectedErr.Message, got.Errors[i].Message)
					assert.Equal(t, expectedErr.Field, got.Errors[i].Field)
				}
			}
		})
	}
}

func TestService_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	logger := slog.Default()
	service := NewService(mockRepo, logger)

	// Create a test user
	testUser, err := user.NewUserWithID(
		"123e4567-e89b-12d3-a456-426614174000",
		"test@example.com",
		"Test User",
		time.Now().Add(-time.Hour),
		time.Now(),
	)
	require.NoError(t, err)

	tests := []struct {
		name    string
		request GetUserRequest
		setup   func()
		want    *GetUserResponse
		wantErr bool
	}{
		{
			name: "successful user retrieval",
			request: GetUserRequest{
				ID: "123e4567-e89b-12d3-a456-426614174000",
			},
			setup: func() {
				userID, _ := user.NewUserID("123e4567-e89b-12d3-a456-426614174000")
				mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(testUser, nil)
			},
			want: &GetUserResponse{
				User: &UserDTO{
					ID:    "123e4567-e89b-12d3-a456-426614174000",
					Email: "test@example.com",
					Name:  "Test User",
				},
			},
			wantErr: false,
		},
		{
			name: "user not found",
			request: GetUserRequest{
				ID: "123e4567-e89b-12d3-a456-426614174000",
			},
			setup: func() {
				userID, _ := user.NewUserID("123e4567-e89b-12d3-a456-426614174000")
				mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(nil, errors.ErrUserNotFound)
			},
			want: &GetUserResponse{
				Errors: []ErrorDTO{
					{
						Message: "User not found",
						Code:    "USER_NOT_FOUND",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid user ID",
			request: GetUserRequest{
				ID: "",
			},
			setup: func() {
				// No mock setup needed as validation happens before repository call
			},
			want: &GetUserResponse{
				Errors: []ErrorDTO{
					{
						Message: "Invalid user ID format",
						Field:   "id",
						Code:    "INVALID_USER_ID",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := service.GetUser(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)

			if tt.want.User != nil {
				require.NotNil(t, got.User)
				assert.Equal(t, tt.want.User.ID, got.User.ID)
				assert.Equal(t, tt.want.User.Email, got.User.Email)
				assert.Equal(t, tt.want.User.Name, got.User.Name)
			}

			if len(tt.want.Errors) > 0 {
				require.Len(t, got.Errors, len(tt.want.Errors))
				for i, expectedErr := range tt.want.Errors {
					assert.Equal(t, expectedErr.Code, got.Errors[i].Code)
					assert.Equal(t, expectedErr.Message, got.Errors[i].Message)
					assert.Equal(t, expectedErr.Field, got.Errors[i].Field)
				}
			}
		})
	}
}

func TestService_ListUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	logger := slog.Default()
	service := NewService(mockRepo, logger)

	// Create test users
	testUser1, err := user.NewUserWithID(
		"123e4567-e89b-12d3-a456-426614174001",
		"user1@example.com",
		"User One",
		time.Now().Add(-2*time.Hour),
		time.Now().Add(-time.Hour),
	)
	require.NoError(t, err)

	testUser2, err := user.NewUserWithID(
		"123e4567-e89b-12d3-a456-426614174002",
		"user2@example.com",
		"User Two",
		time.Now().Add(-time.Hour),
		time.Now(),
	)
	require.NoError(t, err)

	tests := []struct {
		name    string
		request ListUsersRequest
		setup   func()
		want    *ListUsersResponse
		wantErr bool
	}{
		{
			name: "successful user listing with default limit",
			request: ListUsersRequest{
				First: 0, // Should use default of 10
				After: "",
			},
			setup: func() {
				mockRepo.EXPECT().FindAll(gomock.Any(), 11, "").Return([]*user.User{testUser1, testUser2}, "", nil)
			},
			want: &ListUsersResponse{
				Users: &UserConnectionDTO{
					Edges: []*UserEdgeDTO{
						{
							Node: &UserDTO{
								ID:    "123e4567-e89b-12d3-a456-426614174001",
								Email: "user1@example.com",
								Name:  "User One",
							},
							Cursor: "123e4567-e89b-12d3-a456-426614174001",
						},
						{
							Node: &UserDTO{
								ID:    "123e4567-e89b-12d3-a456-426614174002",
								Email: "user2@example.com",
								Name:  "User Two",
							},
							Cursor: "123e4567-e89b-12d3-a456-426614174002",
						},
					},
					PageInfo: &PageInfoDTO{
						HasNextPage:     false,
						HasPreviousPage: false,
						StartCursor:     stringPtr("123e4567-e89b-12d3-a456-426614174001"),
						EndCursor:       stringPtr("123e4567-e89b-12d3-a456-426614174002"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "successful user listing with custom limit",
			request: ListUsersRequest{
				First: 5,
				After: "",
			},
			setup: func() {
				mockRepo.EXPECT().FindAll(gomock.Any(), 6, "").Return([]*user.User{testUser1}, "", nil)
			},
			want: &ListUsersResponse{
				Users: &UserConnectionDTO{
					Edges: []*UserEdgeDTO{
						{
							Node: &UserDTO{
								ID:    "123e4567-e89b-12d3-a456-426614174001",
								Email: "user1@example.com",
								Name:  "User One",
							},
							Cursor: "123e4567-e89b-12d3-a456-426614174001",
						},
					},
					PageInfo: &PageInfoDTO{
						HasNextPage:     false,
						HasPreviousPage: false,
						StartCursor:     stringPtr("123e4567-e89b-12d3-a456-426614174001"),
						EndCursor:       stringPtr("123e4567-e89b-12d3-a456-426614174001"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "successful user listing with pagination",
			request: ListUsersRequest{
				First: 1,
				After: "cursor123",
			},
			setup: func() {
				// Return 2 users to simulate hasNextPage = true
				mockRepo.EXPECT().FindAll(gomock.Any(), 2, "cursor123").Return([]*user.User{testUser1, testUser2}, "", nil)
			},
			want: &ListUsersResponse{
				Users: &UserConnectionDTO{
					Edges: []*UserEdgeDTO{
						{
							Node: &UserDTO{
								ID:    "123e4567-e89b-12d3-a456-426614174001",
								Email: "user1@example.com",
								Name:  "User One",
							},
							Cursor: "123e4567-e89b-12d3-a456-426614174001",
						},
					},
					PageInfo: &PageInfoDTO{
						HasNextPage:     true,
						HasPreviousPage: true,
						StartCursor:     stringPtr("123e4567-e89b-12d3-a456-426614174001"),
						EndCursor:       stringPtr("123e4567-e89b-12d3-a456-426614174001"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty result",
			request: ListUsersRequest{
				First: 10,
				After: "",
			},
			setup: func() {
				mockRepo.EXPECT().FindAll(gomock.Any(), 11, "").Return([]*user.User{}, "", nil)
			},
			want: &ListUsersResponse{
				Users: &UserConnectionDTO{
					Edges: []*UserEdgeDTO{},
					PageInfo: &PageInfoDTO{
						HasNextPage:     false,
						HasPreviousPage: false,
						StartCursor:     nil,
						EndCursor:       nil,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid first parameter - negative",
			request: ListUsersRequest{
				First: -1,
				After: "",
			},
			setup: func() {
				// No mock setup needed as validation happens before repository call
			},
			want: &ListUsersResponse{
				Errors: []ErrorDTO{
					{
						Message: "First parameter must be non-negative",
						Field:   "first",
						Code:    "INVALID_FIRST",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid first parameter - too large",
			request: ListUsersRequest{
				First: 101,
				After: "",
			},
			setup: func() {
				// No mock setup needed as validation happens before repository call
			},
			want: &ListUsersResponse{
				Errors: []ErrorDTO{
					{
						Message: "First parameter cannot exceed 100",
						Field:   "first",
						Code:    "INVALID_FIRST",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "repository error",
			request: ListUsersRequest{
				First: 10,
				After: "",
			},
			setup: func() {
				mockRepo.EXPECT().FindAll(gomock.Any(), 11, "").Return(nil, "", errors.ErrRepositoryOperation)
			},
			want: &ListUsersResponse{
				Errors: []ErrorDTO{
					{
						Message: "Repository operation failed",
						Code:    "REPOSITORY_OPERATION",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := service.ListUsers(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)

			if tt.want.Users != nil {
				require.NotNil(t, got.Users)
				assert.Equal(t, len(tt.want.Users.Edges), len(got.Users.Edges))

				for i, expectedEdge := range tt.want.Users.Edges {
					assert.Equal(t, expectedEdge.Node.ID, got.Users.Edges[i].Node.ID)
					assert.Equal(t, expectedEdge.Node.Email, got.Users.Edges[i].Node.Email)
					assert.Equal(t, expectedEdge.Node.Name, got.Users.Edges[i].Node.Name)
					assert.Equal(t, expectedEdge.Cursor, got.Users.Edges[i].Cursor)
				}

				assert.Equal(t, tt.want.Users.PageInfo.HasNextPage, got.Users.PageInfo.HasNextPage)
				assert.Equal(t, tt.want.Users.PageInfo.HasPreviousPage, got.Users.PageInfo.HasPreviousPage)

				if tt.want.Users.PageInfo.StartCursor != nil {
					require.NotNil(t, got.Users.PageInfo.StartCursor)
					assert.Equal(t, *tt.want.Users.PageInfo.StartCursor, *got.Users.PageInfo.StartCursor)
				} else {
					assert.Nil(t, got.Users.PageInfo.StartCursor)
				}

				if tt.want.Users.PageInfo.EndCursor != nil {
					require.NotNil(t, got.Users.PageInfo.EndCursor)
					assert.Equal(t, *tt.want.Users.PageInfo.EndCursor, *got.Users.PageInfo.EndCursor)
				} else {
					assert.Nil(t, got.Users.PageInfo.EndCursor)
				}
			}

			if len(tt.want.Errors) > 0 {
				require.Len(t, got.Errors, len(tt.want.Errors))
				for i, expectedErr := range tt.want.Errors {
					assert.Equal(t, expectedErr.Code, got.Errors[i].Code)
					assert.Equal(t, expectedErr.Message, got.Errors[i].Message)
					assert.Equal(t, expectedErr.Field, got.Errors[i].Field)
				}
			}
		})
	}
}

func TestService_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	logger := slog.Default()
	service := NewService(mockRepo, logger)

	// Helper function to create a fresh test user for each test
	createTestUser := func() *user.User {
		testUser, err := user.NewUserWithID(
			"123e4567-e89b-12d3-a456-426614174000",
			"test@example.com",
			"Test User",
			time.Now().Add(-time.Hour),
			time.Now(),
		)
		require.NoError(t, err)
		return testUser
	}

	tests := []struct {
		name    string
		request UpdateUserRequest
		setup   func()
		want    *UpdateUserResponse
		wantErr bool
	}{
		{
			name: "successful user update - email only",
			request: UpdateUserRequest{
				ID:    "123e4567-e89b-12d3-a456-426614174000",
				Email: stringPtr("newemail@example.com"),
				Name:  nil,
			},
			setup: func() {
				userID, _ := user.NewUserID("123e4567-e89b-12d3-a456-426614174000")
				newEmail, _ := user.NewEmail("newemail@example.com")

				mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(createTestUser(), nil)
				mockRepo.EXPECT().ExistsByEmail(gomock.Any(), newEmail).Return(false, nil)
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			want: &UpdateUserResponse{
				User: &UserDTO{
					ID:    "123e4567-e89b-12d3-a456-426614174000",
					Email: "newemail@example.com",
					Name:  "Test User",
				},
			},
			wantErr: false,
		},
		{
			name: "successful user update - name only",
			request: UpdateUserRequest{
				ID:    "123e4567-e89b-12d3-a456-426614174000",
				Email: nil,
				Name:  stringPtr("New Name"),
			},
			setup: func() {
				userID, _ := user.NewUserID("123e4567-e89b-12d3-a456-426614174000")

				mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(createTestUser(), nil)
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			want: &UpdateUserResponse{
				User: &UserDTO{
					ID:    "123e4567-e89b-12d3-a456-426614174000",
					Email: "test@example.com",
					Name:  "New Name",
				},
			},
			wantErr: false,
		},
		{
			name: "successful user update - both email and name",
			request: UpdateUserRequest{
				ID:    "123e4567-e89b-12d3-a456-426614174000",
				Email: stringPtr("updated@example.com"),
				Name:  stringPtr("Updated Name"),
			},
			setup: func() {
				userID, _ := user.NewUserID("123e4567-e89b-12d3-a456-426614174000")
				newEmail, _ := user.NewEmail("updated@example.com")

				mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(createTestUser(), nil)
				mockRepo.EXPECT().ExistsByEmail(gomock.Any(), newEmail).Return(false, nil)
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			want: &UpdateUserResponse{
				User: &UserDTO{
					ID:    "123e4567-e89b-12d3-a456-426614174000",
					Email: "updated@example.com",
					Name:  "Updated Name",
				},
			},
			wantErr: false,
		},
		{
			name: "successful user update - same email (no duplicate check)",
			request: UpdateUserRequest{
				ID:    "123e4567-e89b-12d3-a456-426614174000",
				Email: stringPtr("test@example.com"), // Same as current email
				Name:  nil,
			},
			setup: func() {
				userID, _ := user.NewUserID("123e4567-e89b-12d3-a456-426614174000")

				mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(createTestUser(), nil)
				// No ExistsByEmail call expected since email is the same
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			want: &UpdateUserResponse{
				User: &UserDTO{
					ID:    "123e4567-e89b-12d3-a456-426614174000",
					Email: "test@example.com",
					Name:  "Test User",
				},
			},
			wantErr: false,
		},
		{
			name: "user not found",
			request: UpdateUserRequest{
				ID:    "123e4567-e89b-12d3-a456-426614174000",
				Email: stringPtr("newemail@example.com"),
				Name:  nil,
			},
			setup: func() {
				userID, _ := user.NewUserID("123e4567-e89b-12d3-a456-426614174000")
				mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(nil, errors.ErrUserNotFound)
			},
			want: &UpdateUserResponse{
				Errors: []ErrorDTO{
					{
						Message: "User not found",
						Code:    "USER_NOT_FOUND",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "duplicate email error",
			request: UpdateUserRequest{
				ID:    "123e4567-e89b-12d3-a456-426614174000",
				Email: stringPtr("existing@example.com"),
				Name:  nil,
			},
			setup: func() {
				userID, _ := user.NewUserID("123e4567-e89b-12d3-a456-426614174000")
				newEmail, _ := user.NewEmail("existing@example.com")

				mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(createTestUser(), nil)
				mockRepo.EXPECT().ExistsByEmail(gomock.Any(), newEmail).Return(true, nil)
			},
			want: &UpdateUserResponse{
				Errors: []ErrorDTO{
					{
						Message: "Email already exists",
						Field:   "email",
						Code:    "DUPLICATE_EMAIL",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid email format",
			request: UpdateUserRequest{
				ID:    "123e4567-e89b-12d3-a456-426614174000",
				Email: stringPtr("invalid-email"),
				Name:  nil,
			},
			setup: func() {
				userID, _ := user.NewUserID("123e4567-e89b-12d3-a456-426614174000")
				mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(createTestUser(), nil)
			},
			want: &UpdateUserResponse{
				Errors: []ErrorDTO{
					{
						Message: "Invalid email format",
						Field:   "email",
						Code:    "INVALID_EMAIL",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty name",
			request: UpdateUserRequest{
				ID:    "123e4567-e89b-12d3-a456-426614174000",
				Email: nil,
				Name:  stringPtr(""),
			},
			setup: func() {
				// No mock setup needed as validation happens before repository call
			},
			want: &UpdateUserResponse{
				Errors: []ErrorDTO{
					{
						Message: "Name cannot be empty",
						Field:   "name",
						Code:    "INVALID_NAME",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid user ID",
			request: UpdateUserRequest{
				ID:    "",
				Email: stringPtr("test@example.com"),
				Name:  nil,
			},
			setup: func() {
				// No mock setup needed as validation happens before repository call
			},
			want: &UpdateUserResponse{
				Errors: []ErrorDTO{
					{
						Message: "Invalid user ID format",
						Field:   "id",
						Code:    "INVALID_USER_ID",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "repository update error",
			request: UpdateUserRequest{
				ID:    "123e4567-e89b-12d3-a456-426614174000",
				Email: stringPtr("newemail@example.com"),
				Name:  nil,
			},
			setup: func() {
				userID, _ := user.NewUserID("123e4567-e89b-12d3-a456-426614174000")
				newEmail, _ := user.NewEmail("newemail@example.com")

				mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(createTestUser(), nil)
				mockRepo.EXPECT().ExistsByEmail(gomock.Any(), newEmail).Return(false, nil)
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.ErrRepositoryOperation)
			},
			want: &UpdateUserResponse{
				Errors: []ErrorDTO{
					{
						Message: "Repository operation failed",
						Code:    "REPOSITORY_OPERATION",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := service.UpdateUser(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)

			if tt.want.User != nil {
				require.NotNil(t, got.User)
				assert.Equal(t, tt.want.User.ID, got.User.ID)
				assert.Equal(t, tt.want.User.Email, got.User.Email)
				assert.Equal(t, tt.want.User.Name, got.User.Name)
				assert.False(t, got.User.CreatedAt.IsZero())
				assert.False(t, got.User.UpdatedAt.IsZero())
			}

			if len(tt.want.Errors) > 0 {
				require.Len(t, got.Errors, len(tt.want.Errors))
				for i, expectedErr := range tt.want.Errors {
					assert.Equal(t, expectedErr.Code, got.Errors[i].Code)
					assert.Equal(t, expectedErr.Message, got.Errors[i].Message)
					assert.Equal(t, expectedErr.Field, got.Errors[i].Field)
				}
			}
		})
	}
}

func TestService_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	logger := slog.Default()
	service := NewService(mockRepo, logger)

	// Create a test user
	testUser, err := user.NewUserWithID(
		"123e4567-e89b-12d3-a456-426614174000",
		"test@example.com",
		"Test User",
		time.Now().Add(-time.Hour),
		time.Now(),
	)
	require.NoError(t, err)

	tests := []struct {
		name    string
		request DeleteUserRequest
		setup   func()
		want    *DeleteUserResponse
		wantErr bool
	}{
		{
			name: "successful user deletion",
			request: DeleteUserRequest{
				ID: "123e4567-e89b-12d3-a456-426614174000",
			},
			setup: func() {
				userID, _ := user.NewUserID("123e4567-e89b-12d3-a456-426614174000")
				mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(testUser, nil)
				mockRepo.EXPECT().Delete(gomock.Any(), userID).Return(nil)
			},
			want: &DeleteUserResponse{
				Success: true,
			},
			wantErr: false,
		},
		{
			name: "user not found",
			request: DeleteUserRequest{
				ID: "123e4567-e89b-12d3-a456-426614174000",
			},
			setup: func() {
				userID, _ := user.NewUserID("123e4567-e89b-12d3-a456-426614174000")
				mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(nil, errors.ErrUserNotFound)
			},
			want: &DeleteUserResponse{
				Success: false,
				Errors: []ErrorDTO{
					{
						Message: "User not found",
						Code:    "USER_NOT_FOUND",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid user ID",
			request: DeleteUserRequest{
				ID: "",
			},
			setup: func() {
				// No mock setup needed as validation happens before repository call
			},
			want: &DeleteUserResponse{
				Success: false,
				Errors: []ErrorDTO{
					{
						Message: "Invalid user ID format",
						Field:   "id",
						Code:    "INVALID_USER_ID",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "repository delete error",
			request: DeleteUserRequest{
				ID: "123e4567-e89b-12d3-a456-426614174000",
			},
			setup: func() {
				userID, _ := user.NewUserID("123e4567-e89b-12d3-a456-426614174000")
				mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(testUser, nil)
				mockRepo.EXPECT().Delete(gomock.Any(), userID).Return(errors.ErrRepositoryOperation)
			},
			want: &DeleteUserResponse{
				Success: false,
				Errors: []ErrorDTO{
					{
						Message: "Repository operation failed",
						Code:    "REPOSITORY_OPERATION",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := service.DeleteUser(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)

			assert.Equal(t, tt.want.Success, got.Success)

			if len(tt.want.Errors) > 0 {
				require.Len(t, got.Errors, len(tt.want.Errors))
				for i, expectedErr := range tt.want.Errors {
					assert.Equal(t, expectedErr.Code, got.Errors[i].Code)
					assert.Equal(t, expectedErr.Message, got.Errors[i].Message)
					assert.Equal(t, expectedErr.Field, got.Errors[i].Field)
				}
			} else {
				assert.Empty(t, got.Errors)
			}
		})
	}
}

// Test helper functions

func TestMapDomainUserToDTO(t *testing.T) {
	tests := []struct {
		name     string
		input    *user.User
		expected *UserDTO
	}{
		{
			name:     "nil user",
			input:    nil,
			expected: nil,
		},
		{
			name: "valid user",
			input: func() *user.User {
				u, _ := user.NewUserWithID(
					"123e4567-e89b-12d3-a456-426614174000",
					"test@example.com",
					"Test User",
					time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
				)
				return u
			}(),
			expected: &UserDTO{
				ID:        "123e4567-e89b-12d3-a456-426614174000",
				Email:     "test@example.com",
				Name:      "Test User",
				CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapDomainUserToDTO(tt.input)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.Email, result.Email)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tt.expected.UpdatedAt, result.UpdatedAt)
		})
	}
}

func TestMapDomainErrorToDTO(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected ErrorDTO
	}{
		{
			name:     "nil error",
			input:    nil,
			expected: ErrorDTO{},
		},
		{
			name:  "domain error with field",
			input: errors.ErrInvalidEmail,
			expected: ErrorDTO{
				Message: "Invalid email format",
				Field:   "email",
				Code:    "INVALID_EMAIL",
			},
		},
		{
			name:  "domain error without field",
			input: errors.ErrUserNotFound,
			expected: ErrorDTO{
				Message: "User not found",
				Field:   "",
				Code:    "USER_NOT_FOUND",
			},
		},
		{
			name:  "non-domain error",
			input: fmt.Errorf("some generic error"),
			expected: ErrorDTO{
				Message: "An unexpected error occurred",
				Code:    "INTERNAL_ERROR",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapDomainErrorToDTO(tt.input)
			assert.Equal(t, tt.expected.Message, result.Message)
			assert.Equal(t, tt.expected.Field, result.Field)
			assert.Equal(t, tt.expected.Code, result.Code)
		})
	}
}

// Test validation methods

func TestService_ValidateGetUserRequest(t *testing.T) {
	service := &service{}

	tests := []struct {
		name    string
		request GetUserRequest
		wantErr error
	}{
		{
			name:    "valid request",
			request: GetUserRequest{ID: "123e4567-e89b-12d3-a456-426614174000"},
			wantErr: nil,
		},
		{
			name:    "empty ID",
			request: GetUserRequest{ID: ""},
			wantErr: errors.ErrInvalidUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateGetUserRequest(tt.request)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_ValidateListUsersRequest(t *testing.T) {
	service := &service{}

	tests := []struct {
		name    string
		request ListUsersRequest
		wantErr error
	}{
		{
			name:    "valid request",
			request: ListUsersRequest{First: 10, After: ""},
			wantErr: nil,
		},
		{
			name:    "zero first (valid)",
			request: ListUsersRequest{First: 0, After: ""},
			wantErr: nil,
		},
		{
			name:    "negative first",
			request: ListUsersRequest{First: -1, After: ""},
			wantErr: errors.DomainError{Code: "INVALID_FIRST", Message: "First parameter must be non-negative", Field: "first"},
		},
		{
			name:    "first too large",
			request: ListUsersRequest{First: 101, After: ""},
			wantErr: errors.DomainError{Code: "INVALID_FIRST", Message: "First parameter cannot exceed 100", Field: "first"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateListUsersRequest(tt.request)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_ValidateCreateUserRequest(t *testing.T) {
	service := &service{}

	tests := []struct {
		name    string
		request CreateUserRequest
		wantErr error
	}{
		{
			name:    "valid request",
			request: CreateUserRequest{Email: "test@example.com", Name: "Test User"},
			wantErr: nil,
		},
		{
			name:    "empty email",
			request: CreateUserRequest{Email: "", Name: "Test User"},
			wantErr: errors.ErrInvalidEmail,
		},
		{
			name:    "empty name",
			request: CreateUserRequest{Email: "test@example.com", Name: ""},
			wantErr: errors.ErrInvalidName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateCreateUserRequest(tt.request)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_ValidateUpdateUserRequest(t *testing.T) {
	service := &service{}

	tests := []struct {
		name    string
		request UpdateUserRequest
		wantErr error
	}{
		{
			name:    "valid request with email",
			request: UpdateUserRequest{ID: "123", Email: stringPtr("test@example.com"), Name: nil},
			wantErr: nil,
		},
		{
			name:    "valid request with name",
			request: UpdateUserRequest{ID: "123", Email: nil, Name: stringPtr("Test User")},
			wantErr: nil,
		},
		{
			name:    "empty ID",
			request: UpdateUserRequest{ID: "", Email: stringPtr("test@example.com"), Name: nil},
			wantErr: errors.ErrInvalidUserID,
		},
		{
			name:    "empty email pointer",
			request: UpdateUserRequest{ID: "123", Email: stringPtr(""), Name: nil},
			wantErr: errors.ErrInvalidEmail,
		},
		{
			name:    "empty name pointer",
			request: UpdateUserRequest{ID: "123", Email: nil, Name: stringPtr("")},
			wantErr: errors.ErrInvalidName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateUpdateUserRequest(tt.request)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_ValidateDeleteUserRequest(t *testing.T) {
	service := &service{}

	tests := []struct {
		name    string
		request DeleteUserRequest
		wantErr error
	}{
		{
			name:    "valid request",
			request: DeleteUserRequest{ID: "123e4567-e89b-12d3-a456-426614174000"},
			wantErr: nil,
		},
		{
			name:    "empty ID",
			request: DeleteUserRequest{ID: ""},
			wantErr: errors.ErrInvalidUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateDeleteUserRequest(tt.request)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper functions for tests
func stringPtr(s string) *string {
	return &s
}
