package mocks

import (
	"context"
	"testing"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
	"github.com/golang/mock/gomock"
)

// TestMockGeneration verifies that mock generation is working correctly
func TestMockGeneration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Repository mock generation", func(t *testing.T) {
		mockRepo := NewMockRepository(ctrl)

		// Test that we can set up expectations
		userID, _ := user.NewUserID("test-id")
		testUser, _ := user.NewUser("test@example.com", "Test User")

		mockRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(testUser, nil).
			Times(1)

		// Test that the mock works
		result, err := mockRepo.FindByID(context.Background(), userID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result == nil {
			t.Error("Expected user, got nil")
		}
	})

	t.Run("Domain service mock generation", func(t *testing.T) {
		mockDomainService := NewMockDomainService(ctrl)

		// Test that we can set up expectations
		email, _ := user.NewEmail("test@example.com")

		mockDomainService.EXPECT().
			ValidateUniqueEmail(gomock.Any(), email, nil).
			Return(nil).
			Times(1)

		// Test that the mock works
		err := mockDomainService.ValidateUniqueEmail(context.Background(), email, nil)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Test utilities work correctly", func(t *testing.T) {
		mockRepo := NewMockRepository(ctrl)
		utils := NewRepositoryTestUtils(mockRepo)

		// Test that utility methods work
		userID, _ := user.NewUserID("test-id")
		utils.ExpectFindByIDNotFound(userID)

		// Verify the expectation
		_, err := mockRepo.FindByID(context.Background(), userID)
		if err != errors.ErrUserNotFound {
			t.Errorf("Expected ErrUserNotFound, got %v", err)
		}
	})
}
