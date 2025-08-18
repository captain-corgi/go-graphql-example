package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/domain/auth"
	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
	infrastructure_auth "github.com/captain-corgi/go-graphql-example/internal/infrastructure/auth"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// Service defines the interface for authentication application services
type Service interface {
	// Register creates a new user account and returns authentication tokens
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)

	// Login authenticates a user and returns authentication tokens
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)

	// RefreshToken generates new access and refresh tokens using a valid refresh token
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (*AuthResponse, error)

	// Logout invalidates a refresh token
	Logout(ctx context.Context, req LogoutRequest) (*LogoutResponse, error)

	// ValidateAccessToken validates an access token and returns user claims
	ValidateAccessToken(ctx context.Context, token string) (*TokenClaims, error)
}

// service implements the Service interface
type service struct {
	userRepo        user.Repository
	sessionRepo     auth.SessionRepository
	jwtService      *infrastructure_auth.JWTService
	passwordService *infrastructure_auth.PasswordService
	logger          *slog.Logger
}

// NewService creates a new authentication service
func NewService(
	userRepo user.Repository,
	sessionRepo auth.SessionRepository,
	jwtService *infrastructure_auth.JWTService,
	passwordService *infrastructure_auth.PasswordService,
	logger *slog.Logger,
) Service {
	return &service{
		userRepo:        userRepo,
		sessionRepo:     sessionRepo,
		jwtService:      jwtService,
		passwordService: passwordService,
		logger:          logger,
	}
}

// Register creates a new user account and returns authentication tokens
func (s *service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	s.logger.InfoContext(ctx, "Registering new user", "email", req.Email)

	// Validate request
	if err := s.validateRegisterRequest(req); err != nil {
		s.logger.WarnContext(ctx, "Invalid register request", "error", err, "email", req.Email)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Check if user already exists
	email, err := user.NewEmail(req.Email)
	if err != nil {
		s.logger.WarnContext(ctx, "Invalid email format", "error", err, "email", req.Email)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check if user exists", "error", err, "email", req.Email)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	if exists {
		s.logger.WarnContext(ctx, "User already exists", "email", req.Email)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(errors.ErrDuplicateEmail)},
		}, nil
	}

	// Hash password
	passwordHash, err := s.passwordService.HashPassword(req.Password)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to hash password", "error", err)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Create user
	domainUser, err := user.NewUserWithPassword(req.Email, req.Name, passwordHash)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to create domain user", "error", err)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Persist user
	if err := s.userRepo.Create(ctx, domainUser); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create user in repository", "error", err)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Generate tokens
	tokenPair, err := s.jwtService.GenerateTokenPair(domainUser.ID().String(), domainUser.Email().String())
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate tokens", "error", err)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(errors.DomainError{
				Code:    "TOKEN_GENERATION_FAILED",
				Message: "Failed to generate authentication tokens",
				Field:   "tokens",
			})},
		}, nil
	}

	// Create session
	refreshTokenHash := s.hashRefreshToken(tokenPair.RefreshToken)
	session, err := auth.NewSession(
		domainUser.ID(),
		refreshTokenHash,
		time.Now().Add(s.jwtService.GetRefreshTokenTTL()),
		req.DeviceInfo,
		req.IPAddress,
	)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create session", "error", err)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create session in repository", "error", err)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	s.logger.InfoContext(ctx, "Successfully registered user", "userID", domainUser.ID().String(), "email", req.Email)
	return &AuthResponse{
		User:         mapDomainUserToDTO(domainUser),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}

// Login authenticates a user and returns authentication tokens
func (s *service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	s.logger.InfoContext(ctx, "User login attempt", "email", req.Email)

	// Validate request
	if err := s.validateLoginRequest(req); err != nil {
		s.logger.WarnContext(ctx, "Invalid login request", "error", err, "email", req.Email)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Find user by email
	email, err := user.NewEmail(req.Email)
	if err != nil {
		s.logger.WarnContext(ctx, "Invalid email format", "error", err, "email", req.Email)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	domainUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || domainUser == nil {
		s.logger.WarnContext(ctx, "User not found", "error", err, "email", req.Email)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(errors.DomainError{
				Code:    "INVALID_CREDENTIALS",
				Message: "Invalid email or password",
				Field:   "credentials",
			})},
		}, nil
	}

	// Check if user is active
	if !domainUser.IsActive() {
		s.logger.WarnContext(ctx, "Inactive user login attempt", "email", req.Email)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(errors.DomainError{
				Code:    "ACCOUNT_INACTIVE",
				Message: "Account is inactive",
				Field:   "account",
			})},
		}, nil
	}

	// Verify password
	if err := s.passwordService.VerifyPassword(req.Password, domainUser.PasswordHash().String()); err != nil {
		s.logger.WarnContext(ctx, "Invalid password", "email", req.Email)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Record login
	domainUser.RecordLogin()
	if err := s.userRepo.Update(ctx, domainUser); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update user login time", "error", err)
		// Don't fail the login for this
	}

	// Generate tokens
	tokenPair, err := s.jwtService.GenerateTokenPair(domainUser.ID().String(), domainUser.Email().String())
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate tokens", "error", err)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(errors.DomainError{
				Code:    "TOKEN_GENERATION_FAILED",
				Message: "Failed to generate authentication tokens",
				Field:   "tokens",
			})},
		}, nil
	}

	// Create session
	refreshTokenHash := s.hashRefreshToken(tokenPair.RefreshToken)
	session, err := auth.NewSession(
		domainUser.ID(),
		refreshTokenHash,
		time.Now().Add(s.jwtService.GetRefreshTokenTTL()),
		req.DeviceInfo,
		req.IPAddress,
	)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create session", "error", err)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create session in repository", "error", err)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	s.logger.InfoContext(ctx, "Successfully logged in user", "userID", domainUser.ID().String(), "email", req.Email)
	return &AuthResponse{
		User:         mapDomainUserToDTO(domainUser),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}

// RefreshToken generates new access and refresh tokens using a valid refresh token
func (s *service) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*AuthResponse, error) {
	s.logger.InfoContext(ctx, "Refreshing token")

	// Validate request
	if err := s.validateRefreshTokenRequest(req); err != nil {
		s.logger.WarnContext(ctx, "Invalid refresh token request", "error", err)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Find session by refresh token hash
	refreshTokenHash := s.hashRefreshToken(req.RefreshToken)
	refreshTokenHashVO, err := auth.NewRefreshTokenHash(refreshTokenHash)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to create refresh token hash", "error", err)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	session, err := s.sessionRepo.FindByRefreshTokenHash(ctx, refreshTokenHashVO)
	if err != nil {
		s.logger.WarnContext(ctx, "Session not found", "error", err)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(errors.DomainError{
				Code:    "INVALID_REFRESH_TOKEN",
				Message: "Invalid or expired refresh token",
				Field:   "refreshToken",
			})},
		}, nil
	}

	// Check if session is valid
	if !session.IsValid() {
		s.logger.WarnContext(ctx, "Invalid session")
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(errors.DomainError{
				Code:    "INVALID_REFRESH_TOKEN",
				Message: "Invalid or expired refresh token",
				Field:   "refreshToken",
			})},
		}, nil
	}

	// Get user
	domainUser, err := s.userRepo.FindByID(ctx, session.UserID())
	if err != nil {
		s.logger.ErrorContext(ctx, "User not found for session", "error", err)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Check if user is active
	if !domainUser.IsActive() {
		s.logger.WarnContext(ctx, "Inactive user refresh attempt", "userID", domainUser.ID().String())
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(errors.DomainError{
				Code:    "ACCOUNT_INACTIVE",
				Message: "Account is inactive",
				Field:   "account",
			})},
		}, nil
	}

	// Revoke old session
	session.Revoke()
	if err := s.sessionRepo.Update(ctx, session); err != nil {
		s.logger.ErrorContext(ctx, "Failed to revoke old session", "error", err)
		// Continue anyway
	}

	// Generate new tokens
	tokenPair, err := s.jwtService.GenerateTokenPair(domainUser.ID().String(), domainUser.Email().String())
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate tokens", "error", err)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(errors.DomainError{
				Code:    "TOKEN_GENERATION_FAILED",
				Message: "Failed to generate authentication tokens",
				Field:   "tokens",
			})},
		}, nil
	}

	// Create new session
	newRefreshTokenHash := s.hashRefreshToken(tokenPair.RefreshToken)
	newSession, err := auth.NewSession(
		domainUser.ID(),
		newRefreshTokenHash,
		time.Now().Add(s.jwtService.GetRefreshTokenTTL()),
		req.DeviceInfo,
		req.IPAddress,
	)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create new session", "error", err)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	if err := s.sessionRepo.Create(ctx, newSession); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create new session in repository", "error", err)
		return &AuthResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	s.logger.InfoContext(ctx, "Successfully refreshed token", "userID", domainUser.ID().String())
	return &AuthResponse{
		User:         mapDomainUserToDTO(domainUser),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}

// Logout invalidates a refresh token
func (s *service) Logout(ctx context.Context, req LogoutRequest) (*LogoutResponse, error) {
	s.logger.InfoContext(ctx, "User logout")

	// Validate request
	if err := s.validateLogoutRequest(req); err != nil {
		s.logger.WarnContext(ctx, "Invalid logout request", "error", err)
		return &LogoutResponse{
			Success: false,
			Errors:  []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Find and revoke session
	refreshTokenHash := s.hashRefreshToken(req.RefreshToken)
	refreshTokenHashVO, err := auth.NewRefreshTokenHash(refreshTokenHash)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to create refresh token hash", "error", err)
		return &LogoutResponse{
			Success: false,
			Errors:  []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	session, err := s.sessionRepo.FindByRefreshTokenHash(ctx, refreshTokenHashVO)
	if err != nil {
		s.logger.WarnContext(ctx, "Session not found for logout", "error", err)
		// Return success even if session not found (already logged out)
		return &LogoutResponse{Success: true}, nil
	}

	session.Revoke()
	if err := s.sessionRepo.Update(ctx, session); err != nil {
		s.logger.ErrorContext(ctx, "Failed to revoke session", "error", err)
		return &LogoutResponse{
			Success: false,
			Errors:  []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	s.logger.InfoContext(ctx, "Successfully logged out user")
	return &LogoutResponse{Success: true}, nil
}

// ValidateAccessToken validates an access token and returns user claims
func (s *service) ValidateAccessToken(ctx context.Context, token string) (*TokenClaims, error) {
	claims, err := s.jwtService.ValidateAccessToken(token)
	if err != nil {
		return nil, err
	}

	return &TokenClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
	}, nil
}

// Helper methods

func (s *service) hashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Validation methods

func (s *service) validateRegisterRequest(req RegisterRequest) error {
	if req.Email == "" {
		return errors.ErrInvalidEmail
	}
	if req.Name == "" {
		return errors.ErrInvalidName
	}
	if req.Password == "" {
		return errors.DomainError{
			Code:    "INVALID_PASSWORD",
			Message: "Password cannot be empty",
			Field:   "password",
		}
	}
	return nil
}

func (s *service) validateLoginRequest(req LoginRequest) error {
	if req.Email == "" {
		return errors.ErrInvalidEmail
	}
	if req.Password == "" {
		return errors.DomainError{
			Code:    "INVALID_PASSWORD",
			Message: "Password cannot be empty",
			Field:   "password",
		}
	}
	return nil
}

func (s *service) validateRefreshTokenRequest(req RefreshTokenRequest) error {
	if req.RefreshToken == "" {
		return errors.DomainError{
			Code:    "INVALID_REFRESH_TOKEN",
			Message: "Refresh token cannot be empty",
			Field:   "refreshToken",
		}
	}
	return nil
}

func (s *service) validateLogoutRequest(req LogoutRequest) error {
	if req.RefreshToken == "" {
		return errors.DomainError{
			Code:    "INVALID_REFRESH_TOKEN",
			Message: "Refresh token cannot be empty",
			Field:   "refreshToken",
		}
	}
	return nil
}
