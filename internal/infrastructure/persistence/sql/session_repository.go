package sql

import (
	"context"
	"database/sql"
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/domain/auth"
	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
)

// sessionRepository implements the auth.SessionRepository interface
type sessionRepository struct {
	db *sql.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *sql.DB) auth.SessionRepository {
	return &sessionRepository{db: db}
}

// Create persists a new session
func (r *sessionRepository) Create(ctx context.Context, session *auth.Session) error {
	query := `
		INSERT INTO user_sessions (id, user_id, refresh_token_hash, expires_at, is_revoked, device_info, ip_address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query,
		session.ID().String(),
		session.UserID().String(),
		session.RefreshTokenHash().String(),
		session.ExpiresAt(),
		session.IsRevoked(),
		session.DeviceInfo(),
		session.IPAddress(),
		session.CreatedAt(),
		session.UpdatedAt(),
	)

	if err != nil {
		return errors.DomainError{
			Code:    "SESSION_CREATION_FAILED",
			Message: "Failed to create session",
			Field:   "session",
		}
	}

	return nil
}

// FindByRefreshTokenHash retrieves a session by refresh token hash
func (r *sessionRepository) FindByRefreshTokenHash(ctx context.Context, refreshTokenHash auth.RefreshTokenHash) (*auth.Session, error) {
	query := `
		SELECT id, user_id, refresh_token_hash, expires_at, is_revoked, device_info, ip_address, created_at, updated_at
		FROM user_sessions
		WHERE refresh_token_hash = $1 AND is_revoked = false`

	row := r.db.QueryRowContext(ctx, query, refreshTokenHash.String())

	var id, userID, tokenHash string
	var expiresAt, createdAt, updatedAt time.Time
	var isRevoked bool
	var deviceInfo, ipAddress sql.NullString

	err := row.Scan(&id, &userID, &tokenHash, &expiresAt, &isRevoked, &deviceInfo, &ipAddress, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.DomainError{
				Code:    "SESSION_NOT_FOUND",
				Message: "Session not found",
				Field:   "refreshToken",
			}
		}
		return nil, errors.DomainError{
			Code:    "SESSION_QUERY_FAILED",
			Message: "Failed to query session",
			Field:   "refreshToken",
		}
	}

	var deviceInfoPtr, ipAddressPtr *string
	if deviceInfo.Valid {
		deviceInfoPtr = &deviceInfo.String
	}
	if ipAddress.Valid {
		ipAddressPtr = &ipAddress.String
	}

	session, err := auth.NewSessionWithID(id, userID, tokenHash, expiresAt, isRevoked, deviceInfoPtr, ipAddressPtr, createdAt, updatedAt)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// FindByUserID retrieves all sessions for a user
func (r *sessionRepository) FindByUserID(ctx context.Context, userID user.UserID) ([]*auth.Session, error) {
	query := `
		SELECT id, user_id, refresh_token_hash, expires_at, is_revoked, device_info, ip_address, created_at, updated_at
		FROM user_sessions
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, errors.DomainError{
			Code:    "SESSION_QUERY_FAILED",
			Message: "Failed to query sessions",
			Field:   "userId",
		}
	}
	defer rows.Close()

	var sessions []*auth.Session
	for rows.Next() {
		var id, uID, tokenHash string
		var expiresAt, createdAt, updatedAt time.Time
		var isRevoked bool
		var deviceInfo, ipAddress sql.NullString

		err := rows.Scan(&id, &uID, &tokenHash, &expiresAt, &isRevoked, &deviceInfo, &ipAddress, &createdAt, &updatedAt)
		if err != nil {
			return nil, errors.DomainError{
				Code:    "SESSION_SCAN_FAILED",
				Message: "Failed to scan session",
				Field:   "sessions",
			}
		}

		var deviceInfoPtr, ipAddressPtr *string
		if deviceInfo.Valid {
			deviceInfoPtr = &deviceInfo.String
		}
		if ipAddress.Valid {
			ipAddressPtr = &ipAddress.String
		}

		session, err := auth.NewSessionWithID(id, uID, tokenHash, expiresAt, isRevoked, deviceInfoPtr, ipAddressPtr, createdAt, updatedAt)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.DomainError{
			Code:    "SESSION_ITERATION_FAILED",
			Message: "Failed to iterate sessions",
			Field:   "sessions",
		}
	}

	return sessions, nil
}

// Update modifies an existing session
func (r *sessionRepository) Update(ctx context.Context, session *auth.Session) error {
	query := `
		UPDATE user_sessions
		SET refresh_token_hash = $2, expires_at = $3, is_revoked = $4, device_info = $5, ip_address = $6, updated_at = $7
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query,
		session.ID().String(),
		session.RefreshTokenHash().String(),
		session.ExpiresAt(),
		session.IsRevoked(),
		session.DeviceInfo(),
		session.IPAddress(),
		session.UpdatedAt(),
	)

	if err != nil {
		return errors.DomainError{
			Code:    "SESSION_UPDATE_FAILED",
			Message: "Failed to update session",
			Field:   "session",
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.DomainError{
			Code:    "SESSION_UPDATE_CHECK_FAILED",
			Message: "Failed to check session update result",
			Field:   "session",
		}
	}

	if rowsAffected == 0 {
		return errors.DomainError{
			Code:    "SESSION_NOT_FOUND",
			Message: "Session not found for update",
			Field:   "sessionId",
		}
	}

	return nil
}

// Delete removes a session by ID
func (r *sessionRepository) Delete(ctx context.Context, id auth.SessionID) error {
	query := `DELETE FROM user_sessions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return errors.DomainError{
			Code:    "SESSION_DELETE_FAILED",
			Message: "Failed to delete session",
			Field:   "sessionId",
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.DomainError{
			Code:    "SESSION_DELETE_CHECK_FAILED",
			Message: "Failed to check session deletion result",
			Field:   "sessionId",
		}
	}

	if rowsAffected == 0 {
		return errors.DomainError{
			Code:    "SESSION_NOT_FOUND",
			Message: "Session not found for deletion",
			Field:   "sessionId",
		}
	}

	return nil
}

// RevokeByUserID revokes all sessions for a user
func (r *sessionRepository) RevokeByUserID(ctx context.Context, userID user.UserID) error {
	query := `UPDATE user_sessions SET is_revoked = true, updated_at = NOW() WHERE user_id = $1 AND is_revoked = false`

	_, err := r.db.ExecContext(ctx, query, userID.String())
	if err != nil {
		return errors.DomainError{
			Code:    "SESSION_REVOKE_FAILED",
			Message: "Failed to revoke user sessions",
			Field:   "userId",
		}
	}

	return nil
}

// DeleteExpired removes all expired sessions
func (r *sessionRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM user_sessions WHERE expires_at < NOW() OR is_revoked = true`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return errors.DomainError{
			Code:    "SESSION_CLEANUP_FAILED",
			Message: "Failed to clean up expired sessions",
			Field:   "sessions",
		}
	}

	return nil
}

// Count returns the total number of active sessions
func (r *sessionRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM user_sessions WHERE is_revoked = false AND expires_at > NOW()`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, errors.DomainError{
			Code:    "SESSION_COUNT_FAILED",
			Message: "Failed to count sessions",
			Field:   "sessions",
		}
	}

	return count, nil
}

