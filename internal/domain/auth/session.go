package auth

import (
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
	"github.com/google/uuid"
)

// SessionID represents a unique identifier for a session
type SessionID struct {
	value string
}

// NewSessionID creates a new SessionID from a string
func NewSessionID(id string) (SessionID, error) {
	if id == "" {
		return SessionID{}, errors.DomainError{
			Code:    "INVALID_SESSION_ID",
			Message: "Session ID cannot be empty",
			Field:   "sessionId",
		}
	}

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return SessionID{}, errors.DomainError{
			Code:    "INVALID_SESSION_ID",
			Message: "Session ID must be a valid UUID",
			Field:   "sessionId",
		}
	}

	return SessionID{value: id}, nil
}

// GenerateSessionID creates a new random SessionID
func GenerateSessionID() SessionID {
	return SessionID{value: uuid.New().String()}
}

// String returns the string representation of the SessionID
func (id SessionID) String() string {
	return id.value
}

// Equals checks if two SessionIDs are equal
func (id SessionID) Equals(other SessionID) bool {
	return id.value == other.value
}

// RefreshTokenHash represents a hashed refresh token
type RefreshTokenHash struct {
	value string
}

// NewRefreshTokenHash creates a new RefreshTokenHash
func NewRefreshTokenHash(hash string) (RefreshTokenHash, error) {
	if hash == "" {
		return RefreshTokenHash{}, errors.DomainError{
			Code:    "INVALID_REFRESH_TOKEN_HASH",
			Message: "Refresh token hash cannot be empty",
			Field:   "refreshTokenHash",
		}
	}

	return RefreshTokenHash{value: hash}, nil
}

// String returns the string representation of the RefreshTokenHash
func (r RefreshTokenHash) String() string {
	return r.value
}

// Equals checks if two RefreshTokenHashes are equal
func (r RefreshTokenHash) Equals(other RefreshTokenHash) bool {
	return r.value == other.value
}

// Session represents a user session with refresh token
type Session struct {
	id               SessionID
	userID           user.UserID
	refreshTokenHash RefreshTokenHash
	expiresAt        time.Time
	isRevoked        bool
	deviceInfo       *string
	ipAddress        *string
	createdAt        time.Time
	updatedAt        time.Time
}

// NewSession creates a new Session entity
func NewSession(userID user.UserID, refreshTokenHash string, expiresAt time.Time, deviceInfo, ipAddress *string) (*Session, error) {
	refreshTokenHashVO, err := NewRefreshTokenHash(refreshTokenHash)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	if expiresAt.Before(now) {
		return nil, errors.DomainError{
			Code:    "INVALID_EXPIRY_TIME",
			Message: "Session expiry time cannot be in the past",
			Field:   "expiresAt",
		}
	}

	return &Session{
		id:               GenerateSessionID(),
		userID:           userID,
		refreshTokenHash: refreshTokenHashVO,
		expiresAt:        expiresAt,
		isRevoked:        false,
		deviceInfo:       deviceInfo,
		ipAddress:        ipAddress,
		createdAt:        now,
		updatedAt:        now,
	}, nil
}

// NewSessionWithID creates a Session entity with a specific ID (for reconstruction from persistence)
func NewSessionWithID(id, userID, refreshTokenHash string, expiresAt time.Time, isRevoked bool, deviceInfo, ipAddress *string, createdAt, updatedAt time.Time) (*Session, error) {
	sessionID, err := NewSessionID(id)
	if err != nil {
		return nil, err
	}

	userIDVO, err := user.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	refreshTokenHashVO, err := NewRefreshTokenHash(refreshTokenHash)
	if err != nil {
		return nil, err
	}

	return &Session{
		id:               sessionID,
		userID:           userIDVO,
		refreshTokenHash: refreshTokenHashVO,
		expiresAt:        expiresAt,
		isRevoked:        isRevoked,
		deviceInfo:       deviceInfo,
		ipAddress:        ipAddress,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}, nil
}

// ID returns the session's ID
func (s *Session) ID() SessionID {
	return s.id
}

// UserID returns the user ID
func (s *Session) UserID() user.UserID {
	return s.userID
}

// RefreshTokenHash returns the refresh token hash
func (s *Session) RefreshTokenHash() RefreshTokenHash {
	return s.refreshTokenHash
}

// ExpiresAt returns when the session expires
func (s *Session) ExpiresAt() time.Time {
	return s.expiresAt
}

// IsRevoked returns whether the session is revoked
func (s *Session) IsRevoked() bool {
	return s.isRevoked
}

// DeviceInfo returns the device information
func (s *Session) DeviceInfo() *string {
	return s.deviceInfo
}

// IPAddress returns the IP address
func (s *Session) IPAddress() *string {
	return s.ipAddress
}

// CreatedAt returns when the session was created
func (s *Session) CreatedAt() time.Time {
	return s.createdAt
}

// UpdatedAt returns when the session was last updated
func (s *Session) UpdatedAt() time.Time {
	return s.updatedAt
}

// IsExpired checks if the session is expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.expiresAt)
}

// IsValid checks if the session is valid (not revoked and not expired)
func (s *Session) IsValid() bool {
	return !s.isRevoked && !s.IsExpired()
}

// Revoke revokes the session
func (s *Session) Revoke() {
	s.isRevoked = true
	s.updatedAt = time.Now()
}

// Validate performs comprehensive validation of the session entity
func (s *Session) Validate() error {
	if s.id.String() == "" {
		return errors.DomainError{
			Code:    "INVALID_SESSION_ID",
			Message: "Session ID cannot be empty",
			Field:   "sessionId",
		}
	}

	if s.userID.String() == "" {
		return errors.DomainError{
			Code:    "INVALID_USER_ID",
			Message: "User ID cannot be empty",
			Field:   "userId",
		}
	}

	if s.refreshTokenHash.String() == "" {
		return errors.DomainError{
			Code:    "INVALID_REFRESH_TOKEN_HASH",
			Message: "Refresh token hash cannot be empty",
			Field:   "refreshTokenHash",
		}
	}

	if s.expiresAt.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_EXPIRY_TIME",
			Message: "Expiry time cannot be zero",
			Field:   "expiresAt",
		}
	}

	if s.createdAt.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_CREATED_AT",
			Message: "Created at timestamp cannot be zero",
			Field:   "createdAt",
		}
	}

	if s.updatedAt.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_UPDATED_AT",
			Message: "Updated at timestamp cannot be zero",
			Field:   "updatedAt",
		}
	}

	if s.updatedAt.Before(s.createdAt) {
		return errors.DomainError{
			Code:    "INVALID_TIMESTAMPS",
			Message: "Updated at cannot be before created at",
			Field:   "updatedAt",
		}
	}

	return nil
}

