package resolver

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/application/user"
	"github.com/captain-corgi/go-graphql-example/internal/application/user/mocks"
	"github.com/captain-corgi/go-graphql-example/internal/interfaces/graphql/generated"
	"github.com/captain-corgi/go-graphql-example/internal/interfaces/graphql/model"
	"github.com/golang/mock/gomock"
)

func TestGraphQLSchemaExecution(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	resolver := NewResolver(mockUserService, logger)

	// Create the executable schema
	schema := generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	})

	if schema == nil {
		t.Fatal("Failed to create executable schema")
	}

	ctx := context.Background()

	t.Run("User Query Success", func(t *testing.T) {
		// Mock successful user retrieval
		expectedUser := &user.UserDTO{
			ID:        "test-id",
			Email:     "test@example.com",
			Name:      "Test User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserService.EXPECT().
			GetUser(gomock.Any(), user.GetUserRequest{ID: "test-id"}).
			Return(&user.GetUserResponse{User: expectedUser}, nil)

		result, err := resolver.Query().User(ctx, "test-id")

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if result == nil {
			t.Fatal("Expected user result, got nil")
		}

		if result.ID != expectedUser.ID {
			t.Fatalf("Expected user ID %s, got %s", expectedUser.ID, result.ID)
		}
	})

	t.Run("CreateUser Mutation Success", func(t *testing.T) {
		// Mock successful user creation
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

		result, err := resolver.Mutation().CreateUser(ctx, model.CreateUserInput{
			Email: "new@example.com",
			Name:  "New User",
		})

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if result == nil {
			t.Fatal("Expected create user result, got nil")
		}

		if result.User == nil {
			t.Fatal("Expected user in result, got nil")
		}

		if result.User.ID != expectedUser.ID {
			t.Fatalf("Expected user ID %s, got %s", expectedUser.ID, result.User.ID)
		}
	})
}
