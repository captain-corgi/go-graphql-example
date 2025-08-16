package user

import "context"

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// Repository defines the interface for user persistence operations
type Repository interface {
	// FindByID retrieves a user by their ID
	FindByID(ctx context.Context, id UserID) (*User, error)

	// FindByEmail retrieves a user by their email address
	FindByEmail(ctx context.Context, email Email) (*User, error)

	// FindAll retrieves users with pagination support
	// limit: maximum number of users to return
	// cursor: pagination cursor (empty string for first page)
	// Returns users and next cursor (empty if no more pages)
	FindAll(ctx context.Context, limit int, cursor string) ([]*User, string, error)

	// Create persists a new user
	Create(ctx context.Context, user *User) error

	// Update modifies an existing user
	Update(ctx context.Context, user *User) error

	// Delete removes a user by their ID
	Delete(ctx context.Context, id UserID) error

	// ExistsByEmail checks if a user with the given email exists
	ExistsByEmail(ctx context.Context, email Email) (bool, error)

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)
}
