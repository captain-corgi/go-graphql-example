package auth

import (
	"context"

	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// SessionRepository defines the interface for session persistence operations
type SessionRepository interface {
	// Create persists a new session
	Create(ctx context.Context, session *Session) error

	// FindByRefreshTokenHash retrieves a session by refresh token hash
	FindByRefreshTokenHash(ctx context.Context, refreshTokenHash RefreshTokenHash) (*Session, error)

	// FindByUserID retrieves all sessions for a user
	FindByUserID(ctx context.Context, userID user.UserID) ([]*Session, error)

	// Update modifies an existing session
	Update(ctx context.Context, session *Session) error

	// Delete removes a session by ID
	Delete(ctx context.Context, id SessionID) error

	// RevokeByUserID revokes all sessions for a user
	RevokeByUserID(ctx context.Context, userID user.UserID) error

	// DeleteExpired removes all expired sessions
	DeleteExpired(ctx context.Context) error

	// Count returns the total number of active sessions
	Count(ctx context.Context) (int64, error)
}

