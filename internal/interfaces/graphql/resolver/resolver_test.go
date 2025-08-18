package resolver

import (
	"log/slog"
	"os"
	"testing"

	authMocks "github.com/captain-corgi/go-graphql-example/internal/application/auth/mocks"
	userMocks "github.com/captain-corgi/go-graphql-example/internal/application/user/mocks"
	"github.com/golang/mock/gomock"
)

func TestNewResolver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := userMocks.NewMockService(ctrl)
	mockAuthService := authMocks.NewMockService(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	resolver := NewResolver(mockUserService, mockAuthService, logger)

	if resolver == nil {
		t.Fatal("NewResolver returned nil")
	}

	if resolver.logger == nil {
		t.Fatal("Resolver logger is nil")
	}

	if resolver.userService == nil {
		t.Fatal("Resolver userService is nil")
	}

	if resolver.authService == nil {
		t.Fatal("Resolver authService is nil")
	}
}

func TestResolverImplementsInterfaces(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := userMocks.NewMockService(ctrl)
	mockAuthService := authMocks.NewMockService(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	resolver := NewResolver(mockUserService, mockAuthService, logger)

	// Test that resolver implements the generated interfaces
	if resolver.Query() == nil {
		t.Fatal("Query resolver is nil")
	}

	if resolver.Mutation() == nil {
		t.Fatal("Mutation resolver is nil")
	}
}
