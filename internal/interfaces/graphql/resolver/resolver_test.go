package resolver

import (
	"log/slog"
	"os"
	"testing"

	"github.com/captain-corgi/go-graphql-example/internal/application/user/mocks"
	"github.com/golang/mock/gomock"
)

func TestNewResolver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	resolver := NewResolver(mockUserService, logger)

	if resolver == nil {
		t.Fatal("NewResolver returned nil")
	}

	if resolver.logger == nil {
		t.Fatal("Resolver logger is nil")
	}

	if resolver.userService == nil {
		t.Fatal("Resolver userService is nil")
	}
}

func TestResolverImplementsInterfaces(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	resolver := NewResolver(mockUserService, logger)

	// Test that resolver implements the generated interfaces
	if resolver.Query() == nil {
		t.Fatal("Query resolver is nil")
	}

	if resolver.Mutation() == nil {
		t.Fatal("Mutation resolver is nil")
	}
}
