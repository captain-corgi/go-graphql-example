package mocks

import (
	"context"
	"testing"

	"github.com/captain-corgi/go-graphql-example/internal/application/user"
	"github.com/golang/mock/gomock"
)

// TestMockGeneration verifies that mock generation is working correctly
func TestMockGeneration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Application service mock generation", func(t *testing.T) {
		mockService := NewMockService(ctrl)

		// Test that we can set up expectations
		req := user.GetUserRequest{ID: "test-id"}
		expectedResp := &user.GetUserResponse{
			User: &user.UserDTO{
				ID:    "test-id",
				Email: "test@example.com",
				Name:  "Test User",
			},
		}

		mockService.EXPECT().
			GetUser(gomock.Any(), req).
			Return(expectedResp, nil).
			Times(1)

		// Test that the mock works
		result, err := mockService.GetUser(context.Background(), req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result == nil {
			t.Error("Expected response, got nil")
			return
		}
		if result.User == nil {
			t.Error("Expected user in response, got nil")
			return
		}
		if result.User.ID != "test-id" {
			t.Errorf("Expected user ID 'test-id', got %s", result.User.ID)
		}
	})

	t.Run("Test utilities work correctly", func(t *testing.T) {
		mockService := NewMockService(ctrl)
		utils := NewServiceTestUtils(mockService)

		// Test that utility methods work
		testUser := NewTestDTOBuilder().
			WithID("test-id").
			WithEmail("test@example.com").
			WithName("Test User").
			BuildUserDTO()

		utils.ExpectGetUserSuccess("test-id", testUser)

		// Verify the expectation
		req := user.GetUserRequest{ID: "test-id"}
		resp, err := mockService.GetUser(context.Background(), req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if resp.User.ID != "test-id" {
			t.Errorf("Expected user ID 'test-id', got %s", resp.User.ID)
		}
	})

	t.Run("Test DTO builder works correctly", func(t *testing.T) {
		builder := NewTestDTOBuilder()
		userDTO := builder.
			WithID("custom-id").
			WithEmail("custom@example.com").
			WithName("Custom User").
			BuildUserDTO()

		if userDTO.ID != "custom-id" {
			t.Errorf("Expected ID 'custom-id', got %s", userDTO.ID)
		}
		if userDTO.Email != "custom@example.com" {
			t.Errorf("Expected email 'custom@example.com', got %s", userDTO.Email)
		}
		if userDTO.Name != "Custom User" {
			t.Errorf("Expected name 'Custom User', got %s", userDTO.Name)
		}
	})

	t.Run("Test connection builder works correctly", func(t *testing.T) {
		users := []*user.UserDTO{
			{ID: "1", Email: "user1@example.com", Name: "User 1"},
			{ID: "2", Email: "user2@example.com", Name: "User 2"},
		}

		connection := BuildTestUserConnection(users, true, false)

		if len(connection.Edges) != 2 {
			t.Errorf("Expected 2 edges, got %d", len(connection.Edges))
		}
		if !connection.PageInfo.HasNextPage {
			t.Error("Expected HasNextPage to be true")
		}
		if connection.PageInfo.HasPreviousPage {
			t.Error("Expected HasPreviousPage to be false")
		}
		if connection.PageInfo.StartCursor == nil || *connection.PageInfo.StartCursor != "1" {
			t.Error("Expected StartCursor to be '1'")
		}
		if connection.PageInfo.EndCursor == nil || *connection.PageInfo.EndCursor != "2" {
			t.Error("Expected EndCursor to be '2'")
		}
	})
}
