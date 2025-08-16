package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user/mocks"
	"github.com/golang/mock/gomock"
)

var testTime = time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

func TestNewDomainService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := user.NewDomainService(mockRepo)

	if service == nil {
		t.Error("NewDomainService() should not return nil")
	}
}

func TestDomainService_ValidateUniqueEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := user.NewDomainService(mockRepo)

	email, err := user.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}

	userID, err := user.NewUserID("550e8400-e29b-41d4-a716-446655440000")
	if err != nil {
		t.Fatalf("Failed to create userID: %v", err)
	}

	existingUser, err := user.NewUserWithID(userID.String(), email.String(), "Existing User",
		testTime, testTime)
	if err != nil {
		t.Fatalf("Failed to create existing user: %v", err)
	}

	tests := []struct {
		name          string
		email         user.Email
		excludeUserID *user.UserID
		setupMock     func()
		wantErr       error
	}{
		{
			name:          "email is unique - user not found",
			email:         email,
			excludeUserID: nil,
			setupMock: func() {
				mockRepo.EXPECT().
					FindByEmail(gomock.Any(), email).
					Return(nil, errors.ErrUserNotFound)
			},
			wantErr: nil,
		},
		{
			name:          "email is not unique - user exists",
			email:         email,
			excludeUserID: nil,
			setupMock: func() {
				mockRepo.EXPECT().
					FindByEmail(gomock.Any(), email).
					Return(existingUser, nil)
			},
			wantErr: errors.ErrDuplicateEmail,
		},
		{
			name:          "email is unique for update - same user",
			email:         email,
			excludeUserID: &userID,
			setupMock: func() {
				mockRepo.EXPECT().
					FindByEmail(gomock.Any(), email).
					Return(existingUser, nil)
			},
			wantErr: nil,
		},
		{
			name:  "email is not unique for update - different user",
			email: email,
			excludeUserID: func() *user.UserID {
				differentID, _ := user.NewUserID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
				return &differentID
			}(),
			setupMock: func() {
				mockRepo.EXPECT().
					FindByEmail(gomock.Any(), email).
					Return(existingUser, nil)
			},
			wantErr: errors.ErrDuplicateEmail,
		},
		{
			name:          "repository error",
			email:         email,
			excludeUserID: nil,
			setupMock: func() {
				mockRepo.EXPECT().
					FindByEmail(gomock.Any(), email).
					Return(nil, errors.ErrRepositoryConnection)
			},
			wantErr: errors.ErrRepositoryConnection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := service.ValidateUniqueEmail(context.Background(), tt.email, tt.excludeUserID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ValidateUniqueEmail() expected error %v, got nil", tt.wantErr)
					return
				}
				if err != tt.wantErr {
					t.Errorf("ValidateUniqueEmail() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateUniqueEmail() unexpected error = %v", err)
			}
		})
	}
}

func TestDomainService_CanDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := user.NewDomainService(mockRepo)

	userID, err := user.NewUserID("550e8400-e29b-41d4-a716-446655440000")
	if err != nil {
		t.Fatalf("Failed to create userID: %v", err)
	}

	testUser, err := user.NewUserWithID(userID.String(), "test@example.com", "Test User",
		testTime, testTime)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	tests := []struct {
		name      string
		userID    user.UserID
		setupMock func()
		wantErr   error
	}{
		{
			name:   "user can be deleted",
			userID: userID,
			setupMock: func() {
				mockRepo.EXPECT().
					FindByID(gomock.Any(), userID).
					Return(testUser, nil)
			},
			wantErr: nil,
		},
		{
			name:   "user not found",
			userID: userID,
			setupMock: func() {
				mockRepo.EXPECT().
					FindByID(gomock.Any(), userID).
					Return(nil, errors.ErrUserNotFound)
			},
			wantErr: errors.ErrUserNotFound,
		},
		{
			name:   "repository error",
			userID: userID,
			setupMock: func() {
				mockRepo.EXPECT().
					FindByID(gomock.Any(), userID).
					Return(nil, errors.ErrRepositoryConnection)
			},
			wantErr: errors.ErrRepositoryConnection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := service.CanDeleteUser(context.Background(), tt.userID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("CanDeleteUser() expected error %v, got nil", tt.wantErr)
					return
				}
				if err != tt.wantErr {
					t.Errorf("CanDeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("CanDeleteUser() unexpected error = %v", err)
			}
		})
	}
}
