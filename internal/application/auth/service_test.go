package auth

import (
	"context"
	"testing"
	"time"

	"log/slog"
	"os"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	domainAuthMocks "github.com/captain-corgi/go-graphql-example/internal/domain/auth/mocks"
	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
	userMocks "github.com/captain-corgi/go-graphql-example/internal/domain/user/mocks"
	authInfra "github.com/captain-corgi/go-graphql-example/internal/infrastructure/auth"
)

func TestAuthService_Register(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	t.Run("Successful registration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := userMocks.NewMockRepository(ctrl)
		sessionRepo := domainAuthMocks.NewMockSessionRepository(ctrl)
		jwtService := authInfra.NewJWTService("test-secret-key-32-chars-minimum", time.Minute*15, time.Hour*24, "test")
		passwordService := authInfra.NewPasswordService()

		service := NewService(userRepo, sessionRepo, jwtService, passwordService, logger)
		ctx := context.Background()

		// Setup mocks
		userRepo.EXPECT().ExistsByEmail(gomock.Any(), gomock.Any()).Return(false, nil)
		userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		sessionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		req := RegisterRequest{
			Email:    "test@example.com",
			Name:     "Test User",
			Password: "testpassword123",
		}

		resp, err := service.Register(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Empty(t, resp.Errors)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.NotNil(t, resp.User)
		assert.Equal(t, "test@example.com", resp.User.Email)
		assert.Equal(t, "Test User", resp.User.Name)
	})

	t.Run("User already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := userMocks.NewMockRepository(ctrl)
		sessionRepo := domainAuthMocks.NewMockSessionRepository(ctrl)
		jwtService := authInfra.NewJWTService("test-secret-key-32-chars-minimum", time.Minute*15, time.Hour*24, "test")
		passwordService := authInfra.NewPasswordService()

		service := NewService(userRepo, sessionRepo, jwtService, passwordService, logger)

		ctx := context.Background()

		// Setup mocks
		userRepo.EXPECT().ExistsByEmail(gomock.Any(), gomock.Any()).Return(true, nil)

		req := RegisterRequest{
			Email:    "existing@example.com",
			Name:     "Test User",
			Password: "testpassword123",
		}

		resp, err := service.Register(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Errors)
		assert.Empty(t, resp.AccessToken)
	})

	t.Run("Invalid email", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := userMocks.NewMockRepository(ctrl)
		sessionRepo := domainAuthMocks.NewMockSessionRepository(ctrl)
		jwtService := authInfra.NewJWTService("test-secret-key-32-chars-minimum", time.Minute*15, time.Hour*24, "test")
		passwordService := authInfra.NewPasswordService()

		service := NewService(userRepo, sessionRepo, jwtService, passwordService, logger)

		ctx := context.Background()

		req := RegisterRequest{
			Email:    "",
			Name:     "Test User",
			Password: "testpassword123",
		}

		resp, err := service.Register(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Errors)
		assert.Empty(t, resp.AccessToken)
	})

	t.Run("Password too short", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := userMocks.NewMockRepository(ctrl)
		sessionRepo := domainAuthMocks.NewMockSessionRepository(ctrl)
		jwtService := authInfra.NewJWTService("test-secret-key-32-chars-minimum", time.Minute*15, time.Hour*24, "test")
		passwordService := authInfra.NewPasswordService()

		service := NewService(userRepo, sessionRepo, jwtService, passwordService, logger)

		ctx := context.Background()

		// Setup mocks - user doesn't exist, but password validation will fail before user creation
		userRepo.EXPECT().ExistsByEmail(gomock.Any(), gomock.Any()).Return(false, nil)

		req := RegisterRequest{
			Email:    "test@example.com",
			Name:     "Test User",
			Password: "short",
		}

		resp, err := service.Register(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Errors)
		assert.Empty(t, resp.AccessToken)
	})
}

func TestAuthService_Login(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	t.Run("Successful login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := userMocks.NewMockRepository(ctrl)
		sessionRepo := domainAuthMocks.NewMockSessionRepository(ctrl)
		jwtService := authInfra.NewJWTService("test-secret-key-32-chars-minimum", time.Minute*15, time.Hour*24, "test")
		passwordService := authInfra.NewPasswordService()

		service := NewService(userRepo, sessionRepo, jwtService, passwordService, logger)
		ctx := context.Background()

		// Create a test user with hashed password
		password := "testpassword123"
		hashedPassword, err := passwordService.HashPassword(password)
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}
		
		// Test password verification directly
		err = passwordService.VerifyPassword(password, hashedPassword)
		if err != nil {
			t.Fatalf("Password verification failed: %v", err)
		}
		
		testUser, err := user.NewUserWithFullDetails(
			"550e8400-e29b-41d4-a716-446655440000", // Valid UUID
			"test@example.com",
			"Test User",
			hashedPassword,
			true,
			nil,
			time.Now(),
			time.Now(),
		)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
		
		// Verify the user's password hash is correct
		if testUser.PasswordHash().String() != hashedPassword {
			t.Fatalf("Password hash mismatch: got %s, expected %s", testUser.PasswordHash().String(), hashedPassword)
		}

		// Setup mocks
		userRepo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(testUser, nil)
		userRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		sessionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		req := LoginRequest{
			Email:    "test@example.com",
			Password: password,
		}

		resp, err := service.Login(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Empty(t, resp.Errors)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.NotNil(t, resp.User)
		assert.Equal(t, "test@example.com", resp.User.Email)
	})

	t.Run("User not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := userMocks.NewMockRepository(ctrl)
		sessionRepo := domainAuthMocks.NewMockSessionRepository(ctrl)
		jwtService := authInfra.NewJWTService("test-secret-key-32-chars-minimum", time.Minute*15, time.Hour*24, "test")
		passwordService := authInfra.NewPasswordService()

		service := NewService(userRepo, sessionRepo, jwtService, passwordService, logger)

		ctx := context.Background()

		// Setup mocks
		userRepo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(nil, errors.ErrUserNotFound)

		req := LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "testpassword123",
		}

		resp, err := service.Login(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Errors)
		assert.Empty(t, resp.AccessToken)
	})

	t.Run("Invalid password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := userMocks.NewMockRepository(ctrl)
		sessionRepo := domainAuthMocks.NewMockSessionRepository(ctrl)
		jwtService := authInfra.NewJWTService("test-secret-key-32-chars-minimum", time.Minute*15, time.Hour*24, "test")
		passwordService := authInfra.NewPasswordService()

		service := NewService(userRepo, sessionRepo, jwtService, passwordService, logger)

		ctx := context.Background()

		// Create a test user with hashed password
		hashedPassword, _ := passwordService.HashPassword("correctpassword")
		testUser, _ := user.NewUserWithFullDetails(
			"550e8400-e29b-41d4-a716-446655440001", // Valid UUID
			"test@example.com",
			"Test User",
			hashedPassword,
			true,
			nil,
			time.Now(),
			time.Now(),
		)

		// Setup mocks
		userRepo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(testUser, nil)

		req := LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		resp, err := service.Login(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Errors)
		assert.Empty(t, resp.AccessToken)
	})

	t.Run("Inactive user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := userMocks.NewMockRepository(ctrl)
		sessionRepo := domainAuthMocks.NewMockSessionRepository(ctrl)
		jwtService := authInfra.NewJWTService("test-secret-key-32-chars-minimum", time.Minute*15, time.Hour*24, "test")
		passwordService := authInfra.NewPasswordService()

		service := NewService(userRepo, sessionRepo, jwtService, passwordService, logger)

		ctx := context.Background()

		// Create an inactive test user
		hashedPassword, _ := passwordService.HashPassword("testpassword123")
		testUser, _ := user.NewUserWithFullDetails(
			"test-id",
			"test@example.com",
			"Test User",
			hashedPassword,
			false, // inactive
			nil,
			time.Now(),
			time.Now(),
		)

		// Setup mocks
		userRepo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(testUser, nil)

		req := LoginRequest{
			Email:    "test@example.com",
			Password: "testpassword123",
		}

		resp, err := service.Login(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Errors)
		assert.Empty(t, resp.AccessToken)
	})
}
