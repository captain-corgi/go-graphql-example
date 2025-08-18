package sql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
	"github.com/captain-corgi/go-graphql-example/internal/infrastructure/database"
	"github.com/lib/pq"
)

// userRepository implements the user.Repository interface using SQL
type userRepository struct {
	db     *database.DB
	logger *slog.Logger
}

// NewUserRepository creates a new SQL-based user repository
func NewUserRepository(db *database.DB, logger *slog.Logger) user.Repository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

// FindByID retrieves a user by their ID
func (r *userRepository) FindByID(ctx context.Context, id user.UserID) (*user.User, error) {
	r.logger.DebugContext(ctx, "Finding user by ID", "user_id", id.String())

	query := `
		SELECT id, email, name, password_hash, is_active, last_login_at, created_at, updated_at 
		FROM users 
		WHERE id = $1`

	var userID, email, name, passwordHash string
	var isActive bool
	var lastLoginAt sql.NullTime
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&userID, &email, &name, &passwordHash, &isActive, &lastLoginAt, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.DebugContext(ctx, "User not found", "user_id", id.String())
			return nil, errors.ErrUserNotFound
		}
		r.logger.ErrorContext(ctx, "Failed to find user by ID", "error", err, "user_id", id.String())
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	var lastLoginAtPtr *time.Time
	if lastLoginAt.Valid {
		lastLoginAtPtr = &lastLoginAt.Time
	}

	domainUser, err := user.NewUserWithFullDetails(userID, email, name, passwordHash, isActive, lastLoginAtPtr, createdAt, updatedAt)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to create domain user from database record", "error", err)
		return nil, fmt.Errorf("failed to create domain user: %w", err)
	}

	r.logger.DebugContext(ctx, "Successfully found user by ID", "user_id", id.String())
	return domainUser, nil
}

// FindByEmail retrieves a user by their email address
func (r *userRepository) FindByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	r.logger.DebugContext(ctx, "Finding user by email", "email", email.String())

	query := `
		SELECT id, email, name, password_hash, is_active, last_login_at, created_at, updated_at 
		FROM users 
		WHERE email = $1`

	var userID, userEmail, name, passwordHash string
	var isActive bool
	var lastLoginAt sql.NullTime
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(
		&userID, &userEmail, &name, &passwordHash, &isActive, &lastLoginAt, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.DebugContext(ctx, "User not found by email", "email", email.String())
			return nil, errors.ErrUserNotFound
		}
		r.logger.ErrorContext(ctx, "Failed to find user by email", "error", err, "email", email.String())
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	var lastLoginAtPtr *time.Time
	if lastLoginAt.Valid {
		lastLoginAtPtr = &lastLoginAt.Time
	}

	domainUser, err := user.NewUserWithFullDetails(userID, userEmail, name, passwordHash, isActive, lastLoginAtPtr, createdAt, updatedAt)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to create domain user from database record", "error", err)
		return nil, fmt.Errorf("failed to create domain user: %w", err)
	}

	r.logger.DebugContext(ctx, "Successfully found user by email", "email", email.String())
	return domainUser, nil
}

// FindAll retrieves users with pagination support
func (r *userRepository) FindAll(ctx context.Context, limit int, cursor string) ([]*user.User, string, error) {
	r.logger.DebugContext(ctx, "Finding all users", "limit", limit, "cursor", cursor)

	var query string
	var args []interface{}

	if cursor == "" {
		// First page
		query = `
			SELECT id, email, name, password_hash, is_active, last_login_at, created_at, updated_at 
			FROM users 
			ORDER BY created_at DESC, id DESC 
			LIMIT $1`
		args = []interface{}{limit}
	} else {
		// Subsequent pages - cursor-based pagination using created_at and id
		query = `
			SELECT id, email, name, password_hash, is_active, last_login_at, created_at, updated_at 
			FROM users 
			WHERE (created_at, id) < (
				SELECT created_at, id FROM users WHERE id = $1
			)
			ORDER BY created_at DESC, id DESC 
			LIMIT $2`
		args = []interface{}{cursor, limit}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to query users", "error", err)
		return nil, "", fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*user.User
	var nextCursor string

	for rows.Next() {
		var userID, email, name, passwordHash string
		var isActive bool
		var lastLoginAt sql.NullTime
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&userID, &email, &name, &passwordHash, &isActive, &lastLoginAt, &createdAt, &updatedAt); err != nil {
			r.logger.ErrorContext(ctx, "Failed to scan user row", "error", err)
			return nil, "", fmt.Errorf("failed to scan user row: %w", err)
		}

		var lastLoginAtPtr *time.Time
		if lastLoginAt.Valid {
			lastLoginAtPtr = &lastLoginAt.Time
		}

		domainUser, err := user.NewUserWithFullDetails(userID, email, name, passwordHash, isActive, lastLoginAtPtr, createdAt, updatedAt)
		if err != nil {
			r.logger.ErrorContext(ctx, "Failed to create domain user from database record", "error", err)
			return nil, "", fmt.Errorf("failed to create domain user: %w", err)
		}

		users = append(users, domainUser)
		nextCursor = userID // Use the last user's ID as the next cursor
	}

	if err := rows.Err(); err != nil {
		r.logger.ErrorContext(ctx, "Error iterating over user rows", "error", err)
		return nil, "", fmt.Errorf("error iterating over user rows: %w", err)
	}

	r.logger.DebugContext(ctx, "Successfully found users", "count", len(users), "next_cursor", nextCursor)
	return users, nextCursor, nil
}

// Create persists a new user
func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	r.logger.DebugContext(ctx, "Creating user", "user_id", u.ID().String(), "email", u.Email().String())

	query := `
		INSERT INTO users (id, email, name, password_hash, is_active, last_login_at, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.ExecContext(ctx, query,
		u.ID().String(),
		u.Email().String(),
		u.Name().String(),
		u.PasswordHash().String(),
		u.IsActive(),
		u.LastLoginAt(),
		u.CreatedAt(),
		u.UpdatedAt(),
	)

	if err != nil {
		// Check for unique constraint violation (duplicate email)
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" && pqErr.Constraint == "users_email_key" {
				r.logger.WarnContext(ctx, "Duplicate email constraint violation", "email", u.Email().String())
				return errors.ErrDuplicateEmail
			}
		}
		r.logger.ErrorContext(ctx, "Failed to create user", "error", err, "user_id", u.ID().String())
		return fmt.Errorf("failed to create user: %w", err)
	}

	r.logger.InfoContext(ctx, "Successfully created user", "user_id", u.ID().String(), "email", u.Email().String())
	return nil
}

// Update modifies an existing user
func (r *userRepository) Update(ctx context.Context, u *user.User) error {
	r.logger.DebugContext(ctx, "Updating user", "user_id", u.ID().String(), "email", u.Email().String())

	query := `
		UPDATE users 
		SET email = $2, name = $3, password_hash = $4, is_active = $5, last_login_at = $6, updated_at = $7 
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query,
		u.ID().String(),
		u.Email().String(),
		u.Name().String(),
		u.PasswordHash().String(),
		u.IsActive(),
		u.LastLoginAt(),
		u.UpdatedAt(),
	)

	if err != nil {
		// Check for unique constraint violation (duplicate email)
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" && pqErr.Constraint == "users_email_key" {
				r.logger.WarnContext(ctx, "Duplicate email constraint violation during update", "email", u.Email().String())
				return errors.ErrDuplicateEmail
			}
		}
		r.logger.ErrorContext(ctx, "Failed to update user", "error", err, "user_id", u.ID().String())
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to get rows affected after update", "error", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		r.logger.WarnContext(ctx, "No rows affected during user update", "user_id", u.ID().String())
		return errors.ErrUserNotFound
	}

	r.logger.InfoContext(ctx, "Successfully updated user", "user_id", u.ID().String(), "email", u.Email().String())
	return nil
}

// Delete removes a user by their ID
func (r *userRepository) Delete(ctx context.Context, id user.UserID) error {
	r.logger.DebugContext(ctx, "Deleting user", "user_id", id.String())

	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to delete user", "error", err, "user_id", id.String())
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to get rows affected after delete", "error", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		r.logger.WarnContext(ctx, "No rows affected during user delete", "user_id", id.String())
		return errors.ErrUserNotFound
	}

	r.logger.InfoContext(ctx, "Successfully deleted user", "user_id", id.String())
	return nil
}

// ExistsByEmail checks if a user with the given email exists
func (r *userRepository) ExistsByEmail(ctx context.Context, email user.Email) (bool, error) {
	r.logger.DebugContext(ctx, "Checking if user exists by email", "email", email.String())

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(&exists)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to check if user exists by email", "error", err, "email", email.String())
		return false, fmt.Errorf("failed to check if user exists by email: %w", err)
	}

	r.logger.DebugContext(ctx, "Successfully checked user existence by email", "email", email.String(), "exists", exists)
	return exists, nil
}

// Count returns the total number of users
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	r.logger.DebugContext(ctx, "Counting total users")

	query := `SELECT COUNT(*) FROM users`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to count users", "error", err)
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	r.logger.DebugContext(ctx, "Successfully counted users", "count", count)
	return count, nil
}
