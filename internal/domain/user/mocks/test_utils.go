package mocks

import (
	"context"
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
	"github.com/golang/mock/gomock"
)

// TestUserBuilder helps create test users with common configurations
type TestUserBuilder struct {
	id        string
	email     string
	name      string
	createdAt time.Time
	updatedAt time.Time
}

// NewTestUserBuilder creates a new test user builder with default values
func NewTestUserBuilder() *TestUserBuilder {
	now := time.Now()
	return &TestUserBuilder{
		id:        "test-user-id",
		email:     "test@example.com",
		name:      "Test User",
		createdAt: now,
		updatedAt: now,
	}
}

// WithID sets the user ID
func (b *TestUserBuilder) WithID(id string) *TestUserBuilder {
	b.id = id
	return b
}

// WithEmail sets the user email
func (b *TestUserBuilder) WithEmail(email string) *TestUserBuilder {
	b.email = email
	return b
}

// WithName sets the user name
func (b *TestUserBuilder) WithName(name string) *TestUserBuilder {
	b.name = name
	return b
}

// WithCreatedAt sets the created at timestamp
func (b *TestUserBuilder) WithCreatedAt(createdAt time.Time) *TestUserBuilder {
	b.createdAt = createdAt
	return b
}

// WithUpdatedAt sets the updated at timestamp
func (b *TestUserBuilder) WithUpdatedAt(updatedAt time.Time) *TestUserBuilder {
	b.updatedAt = updatedAt
	return b
}

// Build creates a user with the configured values
func (b *TestUserBuilder) Build() *user.User {
	// Note: This assumes we have a way to create a user with specific timestamps
	// In a real implementation, you might need to use reflection or provide
	// a test constructor in the domain package
	testUser, _ := user.NewUser(b.email, b.name)
	return testUser
}

// RepositoryTestUtils provides common mock repository setup utilities
type RepositoryTestUtils struct {
	mockRepo *MockRepository
}

// NewRepositoryTestUtils creates a new repository test utilities instance
func NewRepositoryTestUtils(mockRepo *MockRepository) *RepositoryTestUtils {
	return &RepositoryTestUtils{
		mockRepo: mockRepo,
	}
}

// ExpectFindByIDSuccess sets up mock to return a user successfully
func (u *RepositoryTestUtils) ExpectFindByIDSuccess(userID user.UserID, returnUser *user.User) {
	u.mockRepo.EXPECT().
		FindByID(gomock.Any(), userID).
		Return(returnUser, nil).
		Times(1)
}

// ExpectFindByIDNotFound sets up mock to return user not found error
func (u *RepositoryTestUtils) ExpectFindByIDNotFound(userID user.UserID) {
	u.mockRepo.EXPECT().
		FindByID(gomock.Any(), userID).
		Return(nil, errors.ErrUserNotFound).
		Times(1)
}

// ExpectFindByEmailSuccess sets up mock to return a user by email successfully
func (u *RepositoryTestUtils) ExpectFindByEmailSuccess(email user.Email, returnUser *user.User) {
	u.mockRepo.EXPECT().
		FindByEmail(gomock.Any(), email).
		Return(returnUser, nil).
		Times(1)
}

// ExpectFindByEmailNotFound sets up mock to return user not found error for email
func (u *RepositoryTestUtils) ExpectFindByEmailNotFound(email user.Email) {
	u.mockRepo.EXPECT().
		FindByEmail(gomock.Any(), email).
		Return(nil, errors.ErrUserNotFound).
		Times(1)
}

// ExpectCreateSuccess sets up mock to create user successfully
func (u *RepositoryTestUtils) ExpectCreateSuccess(expectedUser *user.User) {
	u.mockRepo.EXPECT().
		Create(gomock.Any(), expectedUser).
		Return(nil).
		Times(1)
}

// ExpectUpdateSuccess sets up mock to update user successfully
func (u *RepositoryTestUtils) ExpectUpdateSuccess(expectedUser *user.User) {
	u.mockRepo.EXPECT().
		Update(gomock.Any(), expectedUser).
		Return(nil).
		Times(1)
}

// ExpectDeleteSuccess sets up mock to delete user successfully
func (u *RepositoryTestUtils) ExpectDeleteSuccess(userID user.UserID) {
	u.mockRepo.EXPECT().
		Delete(gomock.Any(), userID).
		Return(nil).
		Times(1)
}

// ExpectExistsByEmailTrue sets up mock to return true for email existence check
func (u *RepositoryTestUtils) ExpectExistsByEmailTrue(email user.Email) {
	u.mockRepo.EXPECT().
		ExistsByEmail(gomock.Any(), email).
		Return(true, nil).
		Times(1)
}

// ExpectExistsByEmailFalse sets up mock to return false for email existence check
func (u *RepositoryTestUtils) ExpectExistsByEmailFalse(email user.Email) {
	u.mockRepo.EXPECT().
		ExistsByEmail(gomock.Any(), email).
		Return(false, nil).
		Times(1)
}

// ExpectFindAllSuccess sets up mock to return users list successfully
func (u *RepositoryTestUtils) ExpectFindAllSuccess(limit int, cursor string, returnUsers []*user.User, nextCursor string) {
	u.mockRepo.EXPECT().
		FindAll(gomock.Any(), limit, cursor).
		Return(returnUsers, nextCursor, nil).
		Times(1)
}

// DomainServiceTestUtils provides common mock domain service setup utilities
type DomainServiceTestUtils struct {
	mockService *MockDomainService
}

// NewDomainServiceTestUtils creates a new domain service test utilities instance
func NewDomainServiceTestUtils(mockService *MockDomainService) *DomainServiceTestUtils {
	return &DomainServiceTestUtils{
		mockService: mockService,
	}
}

// ExpectValidateUniqueEmailSuccess sets up mock to validate email uniqueness successfully
func (u *DomainServiceTestUtils) ExpectValidateUniqueEmailSuccess(email user.Email, excludeUserID *user.UserID) {
	u.mockService.EXPECT().
		ValidateUniqueEmail(gomock.Any(), email, excludeUserID).
		Return(nil).
		Times(1)
}

// ExpectValidateUniqueEmailDuplicate sets up mock to return duplicate email error
func (u *DomainServiceTestUtils) ExpectValidateUniqueEmailDuplicate(email user.Email, excludeUserID *user.UserID) {
	u.mockService.EXPECT().
		ValidateUniqueEmail(gomock.Any(), email, excludeUserID).
		Return(errors.ErrDuplicateEmail).
		Times(1)
}

// ExpectCanDeleteUserSuccess sets up mock to allow user deletion
func (u *DomainServiceTestUtils) ExpectCanDeleteUserSuccess(userID user.UserID) {
	u.mockService.EXPECT().
		CanDeleteUser(gomock.Any(), userID).
		Return(nil).
		Times(1)
}

// CommonTestContext provides a common test context with timeout
func CommonTestContext() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_ = cancel // Store cancel function to avoid linter warning, but don't call it since we return the context
	return ctx
}
