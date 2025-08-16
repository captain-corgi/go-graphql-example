package user

import (
	"context"
	"log/slog"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// Service defines the interface for user application services
type Service interface {
	// GetUser retrieves a user by ID
	GetUser(ctx context.Context, req GetUserRequest) (*GetUserResponse, error)

	// ListUsers retrieves a paginated list of users
	ListUsers(ctx context.Context, req ListUsersRequest) (*ListUsersResponse, error)

	// CreateUser creates a new user
	CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error)

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, req UpdateUserRequest) (*UpdateUserResponse, error)

	// DeleteUser deletes a user by ID
	DeleteUser(ctx context.Context, req DeleteUserRequest) (*DeleteUserResponse, error)
}

// service implements the Service interface
type service struct {
	userRepo user.Repository
	logger   *slog.Logger
}

// NewService creates a new user service
func NewService(userRepo user.Repository, logger *slog.Logger) Service {
	return &service{
		userRepo: userRepo,
		logger:   logger,
	}
}

// GetUser retrieves a user by ID
func (s *service) GetUser(ctx context.Context, req GetUserRequest) (*GetUserResponse, error) {
	s.logger.InfoContext(ctx, "Getting user", "userID", req.ID)

	// Validate request
	if err := s.validateGetUserRequest(req); err != nil {
		s.logger.WarnContext(ctx, "Invalid get user request", "error", err, "userID", req.ID)
		return &GetUserResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Convert string ID to domain UserID
	userID, err := user.NewUserID(req.ID)
	if err != nil {
		s.logger.WarnContext(ctx, "Invalid user ID format", "error", err, "userID", req.ID)
		return &GetUserResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Retrieve user from repository
	domainUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get user from repository", "error", err, "userID", req.ID)
		return &GetUserResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	s.logger.InfoContext(ctx, "Successfully retrieved user", "userID", req.ID)
	return &GetUserResponse{
		User: mapDomainUserToDTO(domainUser),
	}, nil
}

// ListUsers retrieves a paginated list of users
func (s *service) ListUsers(ctx context.Context, req ListUsersRequest) (*ListUsersResponse, error) {
	s.logger.InfoContext(ctx, "Listing users", "first", req.First, "after", req.After)

	// Validate request
	if err := s.validateListUsersRequest(req); err != nil {
		s.logger.WarnContext(ctx, "Invalid list users request", "error", err)
		return &ListUsersResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Set default limit if not provided
	limit := req.First
	if limit == 0 {
		limit = 10 // Default page size
	}

	// Retrieve users from repository
	domainUsers, nextCursor, err := s.userRepo.FindAll(ctx, limit+1, req.After) // +1 to check if there's a next page
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list users from repository", "error", err)
		return &ListUsersResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Build connection response
	connection := s.buildUserConnection(domainUsers, limit, req.After, nextCursor)

	s.logger.InfoContext(ctx, "Successfully listed users", "count", len(connection.Edges))
	return &ListUsersResponse{
		Users: connection,
	}, nil
}

// CreateUser creates a new user
func (s *service) CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	s.logger.InfoContext(ctx, "Creating user", "email", req.Email, "name", req.Name)

	// Validate request
	if err := s.validateCreateUserRequest(req); err != nil {
		s.logger.WarnContext(ctx, "Invalid create user request", "error", err)
		return &CreateUserResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Check if user with email already exists
	email, err := user.NewEmail(req.Email)
	if err != nil {
		s.logger.WarnContext(ctx, "Invalid email format", "error", err, "email", req.Email)
		return &CreateUserResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check if user exists", "error", err, "email", req.Email)
		return &CreateUserResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	if exists {
		s.logger.WarnContext(ctx, "User with email already exists", "email", req.Email)
		return &CreateUserResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(errors.ErrDuplicateEmail)},
		}, nil
	}

	// Create domain user
	domainUser, err := user.NewUser(req.Email, req.Name)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to create domain user", "error", err)
		return &CreateUserResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Persist user
	if err := s.userRepo.Create(ctx, domainUser); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create user in repository", "error", err)
		return &CreateUserResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	s.logger.InfoContext(ctx, "Successfully created user", "userID", domainUser.ID().String(), "email", req.Email)
	return &CreateUserResponse{
		User: mapDomainUserToDTO(domainUser),
	}, nil
}

// UpdateUser updates an existing user
func (s *service) UpdateUser(ctx context.Context, req UpdateUserRequest) (*UpdateUserResponse, error) {
	s.logger.InfoContext(ctx, "Updating user", "userID", req.ID)

	// Validate request
	if err := s.validateUpdateUserRequest(req); err != nil {
		s.logger.WarnContext(ctx, "Invalid update user request", "error", err)
		return &UpdateUserResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Convert string ID to domain UserID
	userID, err := user.NewUserID(req.ID)
	if err != nil {
		s.logger.WarnContext(ctx, "Invalid user ID format", "error", err, "userID", req.ID)
		return &UpdateUserResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Retrieve existing user
	domainUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get user for update", "error", err, "userID", req.ID)
		return &UpdateUserResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Update email if provided
	if req.Email != nil {
		// Check if new email is already taken by another user
		newEmail, err := user.NewEmail(*req.Email)
		if err != nil {
			s.logger.WarnContext(ctx, "Invalid email format", "error", err, "email", *req.Email)
			return &UpdateUserResponse{
				Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
			}, nil
		}

		// Only check for duplicates if the email is actually changing
		if !domainUser.Email().Equals(newEmail) {
			exists, err := s.userRepo.ExistsByEmail(ctx, newEmail)
			if err != nil {
				s.logger.ErrorContext(ctx, "Failed to check if email exists", "error", err, "email", *req.Email)
				return &UpdateUserResponse{
					Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
				}, nil
			}

			if exists {
				s.logger.WarnContext(ctx, "Email already exists", "email", *req.Email)
				return &UpdateUserResponse{
					Errors: []ErrorDTO{mapDomainErrorToDTO(errors.ErrDuplicateEmail)},
				}, nil
			}
		}

		if err := domainUser.UpdateEmail(*req.Email); err != nil {
			s.logger.WarnContext(ctx, "Failed to update user email", "error", err)
			return &UpdateUserResponse{
				Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
			}, nil
		}
	}

	// Update name if provided
	if req.Name != nil {
		if err := domainUser.UpdateName(*req.Name); err != nil {
			s.logger.WarnContext(ctx, "Failed to update user name", "error", err)
			return &UpdateUserResponse{
				Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
			}, nil
		}
	}

	// Persist updated user
	if err := s.userRepo.Update(ctx, domainUser); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update user in repository", "error", err)
		return &UpdateUserResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	s.logger.InfoContext(ctx, "Successfully updated user", "userID", req.ID)
	return &UpdateUserResponse{
		User: mapDomainUserToDTO(domainUser),
	}, nil
}

// DeleteUser deletes a user by ID
func (s *service) DeleteUser(ctx context.Context, req DeleteUserRequest) (*DeleteUserResponse, error) {
	s.logger.InfoContext(ctx, "Deleting user", "userID", req.ID)

	// Validate request
	if err := s.validateDeleteUserRequest(req); err != nil {
		s.logger.WarnContext(ctx, "Invalid delete user request", "error", err)
		return &DeleteUserResponse{
			Success: false,
			Errors:  []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Convert string ID to domain UserID
	userID, err := user.NewUserID(req.ID)
	if err != nil {
		s.logger.WarnContext(ctx, "Invalid user ID format", "error", err, "userID", req.ID)
		return &DeleteUserResponse{
			Success: false,
			Errors:  []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Check if user exists before deletion
	_, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "User not found for deletion", "error", err, "userID", req.ID)
		return &DeleteUserResponse{
			Success: false,
			Errors:  []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Delete user
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete user from repository", "error", err)
		return &DeleteUserResponse{
			Success: false,
			Errors:  []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	s.logger.InfoContext(ctx, "Successfully deleted user", "userID", req.ID)
	return &DeleteUserResponse{
		Success: true,
	}, nil
}

// Validation methods

func (s *service) validateGetUserRequest(req GetUserRequest) error {
	if req.ID == "" {
		return errors.ErrInvalidUserID
	}
	return nil
}

func (s *service) validateListUsersRequest(req ListUsersRequest) error {
	if req.First < 0 {
		return errors.DomainError{
			Code:    "INVALID_FIRST",
			Message: "First parameter must be non-negative",
			Field:   "first",
		}
	}
	if req.First > 100 {
		return errors.DomainError{
			Code:    "INVALID_FIRST",
			Message: "First parameter cannot exceed 100",
			Field:   "first",
		}
	}
	return nil
}

func (s *service) validateCreateUserRequest(req CreateUserRequest) error {
	if req.Email == "" {
		return errors.ErrInvalidEmail
	}
	if req.Name == "" {
		return errors.ErrInvalidName
	}
	return nil
}

func (s *service) validateUpdateUserRequest(req UpdateUserRequest) error {
	if req.ID == "" {
		return errors.ErrInvalidUserID
	}
	if req.Email != nil && *req.Email == "" {
		return errors.ErrInvalidEmail
	}
	if req.Name != nil && *req.Name == "" {
		return errors.ErrInvalidName
	}
	return nil
}

func (s *service) validateDeleteUserRequest(req DeleteUserRequest) error {
	if req.ID == "" {
		return errors.ErrInvalidUserID
	}
	return nil
}

// Helper methods

func (s *service) buildUserConnection(users []*user.User, limit int, after string, nextCursor string) *UserConnectionDTO {
	hasNextPage := len(users) > limit
	if hasNextPage {
		users = users[:limit] // Remove the extra user we fetched
	}

	edges := make([]*UserEdgeDTO, len(users))
	for i, u := range users {
		edges[i] = &UserEdgeDTO{
			Node:   mapDomainUserToDTO(u),
			Cursor: u.ID().String(), // Using user ID as cursor for simplicity
		}
	}

	var startCursor, endCursor *string
	if len(edges) > 0 {
		start := edges[0].Cursor
		end := edges[len(edges)-1].Cursor
		startCursor = &start
		endCursor = &end
	}

	return &UserConnectionDTO{
		Edges: edges,
		PageInfo: &PageInfoDTO{
			HasNextPage:     hasNextPage,
			HasPreviousPage: after != "",
			StartCursor:     startCursor,
			EndCursor:       endCursor,
		},
	}
}
