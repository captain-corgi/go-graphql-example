package sql

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
	"github.com/captain-corgi/go-graphql-example/internal/infrastructure/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// UserRepositoryTestSuite defines the test suite for user repository
type UserRepositoryTestSuite struct {
	suite.Suite
	db         *database.DB
	repository user.Repository
	ctx        context.Context
	cleanup    func()
}

// SetupSuite sets up the test suite
func (suite *UserRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	// Use test database utilities for setup
	db, cleanup := database.TestDBSetup(suite.T(), "../../../../migrations")
	suite.db = db
	suite.cleanup = cleanup

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	suite.repository = NewUserRepository(db, logger)
}

// TearDownSuite cleans up the test suite
func (suite *UserRepositoryTestSuite) TearDownSuite() {
	if suite.cleanup != nil {
		suite.cleanup()
	}
}

// SetupTest sets up each test
func (suite *UserRepositoryTestSuite) SetupTest() {
	// Clean up users table before each test
	_, err := suite.db.ExecContext(suite.ctx, "DELETE FROM users")
	require.NoError(suite.T(), err)
}

// TestCreateUser tests user creation
func (suite *UserRepositoryTestSuite) TestCreateUser() {
	// Create a test user
	testUser, err := user.NewUser("test@example.com", "Test User")
	require.NoError(suite.T(), err)

	// Create user in repository
	err = suite.repository.Create(suite.ctx, testUser)
	assert.NoError(suite.T(), err)

	// Verify user was created by finding it
	foundUser, err := suite.repository.FindByID(suite.ctx, testUser.ID())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), testUser.ID().String(), foundUser.ID().String())
	assert.Equal(suite.T(), testUser.Email().String(), foundUser.Email().String())
	assert.Equal(suite.T(), testUser.Name().String(), foundUser.Name().String())
}

// TestCreateUserDuplicateEmail tests duplicate email constraint
func (suite *UserRepositoryTestSuite) TestCreateUserDuplicateEmail() {
	// Create first user
	user1, err := user.NewUser("test@example.com", "Test User 1")
	require.NoError(suite.T(), err)

	err = suite.repository.Create(suite.ctx, user1)
	require.NoError(suite.T(), err)

	// Try to create second user with same email
	user2, err := user.NewUser("test@example.com", "Test User 2")
	require.NoError(suite.T(), err)

	err = suite.repository.Create(suite.ctx, user2)
	assert.Equal(suite.T(), errors.ErrDuplicateEmail, err)
}

// TestFindByID tests finding user by ID
func (suite *UserRepositoryTestSuite) TestFindByID() {
	// Create a test user
	testUser, err := user.NewUser("test@example.com", "Test User")
	require.NoError(suite.T(), err)

	err = suite.repository.Create(suite.ctx, testUser)
	require.NoError(suite.T(), err)

	// Find user by ID
	foundUser, err := suite.repository.FindByID(suite.ctx, testUser.ID())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), testUser.ID().String(), foundUser.ID().String())
	assert.Equal(suite.T(), testUser.Email().String(), foundUser.Email().String())
	assert.Equal(suite.T(), testUser.Name().String(), foundUser.Name().String())
}

// TestFindByIDNotFound tests finding non-existent user
func (suite *UserRepositoryTestSuite) TestFindByIDNotFound() {
	nonExistentID := user.GenerateUserID()
	_, err := suite.repository.FindByID(suite.ctx, nonExistentID)
	assert.Equal(suite.T(), errors.ErrUserNotFound, err)
}

// TestFindByEmail tests finding user by email
func (suite *UserRepositoryTestSuite) TestFindByEmail() {
	// Create a test user
	testUser, err := user.NewUser("test@example.com", "Test User")
	require.NoError(suite.T(), err)

	err = suite.repository.Create(suite.ctx, testUser)
	require.NoError(suite.T(), err)

	// Find user by email
	foundUser, err := suite.repository.FindByEmail(suite.ctx, testUser.Email())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), testUser.ID().String(), foundUser.ID().String())
	assert.Equal(suite.T(), testUser.Email().String(), foundUser.Email().String())
	assert.Equal(suite.T(), testUser.Name().String(), foundUser.Name().String())
}

// TestUpdateUser tests user update
func (suite *UserRepositoryTestSuite) TestUpdateUser() {
	// Create a test user
	testUser, err := user.NewUser("test@example.com", "Test User")
	require.NoError(suite.T(), err)

	err = suite.repository.Create(suite.ctx, testUser)
	require.NoError(suite.T(), err)

	// Update user
	err = testUser.UpdateEmail("updated@example.com")
	require.NoError(suite.T(), err)
	err = testUser.UpdateName("Updated User")
	require.NoError(suite.T(), err)

	err = suite.repository.Update(suite.ctx, testUser)
	assert.NoError(suite.T(), err)

	// Verify update
	foundUser, err := suite.repository.FindByID(suite.ctx, testUser.ID())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "updated@example.com", foundUser.Email().String())
	assert.Equal(suite.T(), "Updated User", foundUser.Name().String())
}

// TestDeleteUser tests user deletion
func (suite *UserRepositoryTestSuite) TestDeleteUser() {
	// Create a test user
	testUser, err := user.NewUser("test@example.com", "Test User")
	require.NoError(suite.T(), err)

	err = suite.repository.Create(suite.ctx, testUser)
	require.NoError(suite.T(), err)

	// Delete user
	err = suite.repository.Delete(suite.ctx, testUser.ID())
	assert.NoError(suite.T(), err)

	// Verify deletion
	_, err = suite.repository.FindByID(suite.ctx, testUser.ID())
	assert.Equal(suite.T(), errors.ErrUserNotFound, err)
}

// TestExistsByEmail tests checking user existence by email
func (suite *UserRepositoryTestSuite) TestExistsByEmail() {
	// Create a test user
	testUser, err := user.NewUser("test@example.com", "Test User")
	require.NoError(suite.T(), err)

	// Check existence before creation
	exists, err := suite.repository.ExistsByEmail(suite.ctx, testUser.Email())
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)

	// Create user
	err = suite.repository.Create(suite.ctx, testUser)
	require.NoError(suite.T(), err)

	// Check existence after creation
	exists, err = suite.repository.ExistsByEmail(suite.ctx, testUser.Email())
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)
}

// TestFindAll tests pagination
func (suite *UserRepositoryTestSuite) TestFindAll() {
	// Create multiple test users
	users := make([]*user.User, 5)
	for i := 0; i < 5; i++ {
		u, err := user.NewUser(fmt.Sprintf("user%d@example.com", i), fmt.Sprintf("User %d", i))
		require.NoError(suite.T(), err)

		err = suite.repository.Create(suite.ctx, u)
		require.NoError(suite.T(), err)

		users[i] = u

		// Add small delay to ensure different created_at timestamps
		time.Sleep(10 * time.Millisecond)
	}

	// Test first page
	foundUsers, nextCursor, err := suite.repository.FindAll(suite.ctx, 3, "")
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), foundUsers, 3)
	assert.NotEmpty(suite.T(), nextCursor)

	// Test second page
	remainingUsers, finalCursor, err := suite.repository.FindAll(suite.ctx, 3, nextCursor)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), remainingUsers, 2)
	assert.NotEmpty(suite.T(), finalCursor)
}

// TestCount tests user count
func (suite *UserRepositoryTestSuite) TestCount() {
	// Initial count should be 0
	count, err := suite.repository.Count(suite.ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(0), count)

	// Create some users
	for i := 0; i < 3; i++ {
		testUser, err := user.NewUser(fmt.Sprintf("user%d@example.com", i), fmt.Sprintf("User %d", i))
		require.NoError(suite.T(), err)

		err = suite.repository.Create(suite.ctx, testUser)
		require.NoError(suite.T(), err)
	}

	// Count should be 3
	count, err = suite.repository.Count(suite.ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), count)
}

// TestConcurrentOperations tests concurrent repository operations
func (suite *UserRepositoryTestSuite) TestConcurrentOperations() {
	// Create multiple users concurrently
	const numUsers = 10
	userChan := make(chan *user.User, numUsers)
	errChan := make(chan error, numUsers)

	// Create users concurrently
	for i := 0; i < numUsers; i++ {
		go func(index int) {
			testUser, err := user.NewUser(fmt.Sprintf("concurrent%d@example.com", index), fmt.Sprintf("Concurrent User %d", index))
			if err != nil {
				errChan <- err
				return
			}

			err = suite.repository.Create(suite.ctx, testUser)
			if err != nil {
				errChan <- err
				return
			}

			userChan <- testUser
		}(i)
	}

	// Collect results
	var createdUsers []*user.User
	for i := 0; i < numUsers; i++ {
		select {
		case u := <-userChan:
			createdUsers = append(createdUsers, u)
		case err := <-errChan:
			suite.T().Errorf("Concurrent operation failed: %v", err)
		case <-time.After(10 * time.Second):
			suite.T().Fatal("Timeout waiting for concurrent operations")
		}
	}

	// Verify all users were created
	assert.Len(suite.T(), createdUsers, numUsers)

	// Verify count
	count, err := suite.repository.Count(suite.ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(numUsers), count)
}

// TestRepositoryWithDatabaseFailure tests repository behavior during database failures
func (suite *UserRepositoryTestSuite) TestRepositoryWithDatabaseFailure() {
	// Create a user first
	testUser, err := user.NewUser("test@example.com", "Test User")
	require.NoError(suite.T(), err)

	err = suite.repository.Create(suite.ctx, testUser)
	require.NoError(suite.T(), err)

	// Test with cancelled context (simulates timeout/cancellation)
	cancelledCtx, cancel := context.WithCancel(suite.ctx)
	cancel() // Cancel immediately

	// Operations should fail with cancelled context
	_, err = suite.repository.FindByID(cancelledCtx, testUser.ID())
	assert.Error(suite.T(), err)

	// Test with very short timeout
	timeoutCtx, cancel := context.WithTimeout(suite.ctx, 1*time.Nanosecond)
	defer cancel()

	_, _ = suite.repository.FindByID(timeoutCtx, testUser.ID())
	// This might or might not fail depending on timing, but should not panic
}

// TestRepositoryTransactionBehavior tests repository behavior within transactions
func (suite *UserRepositoryTestSuite) TestRepositoryTransactionBehavior() {
	// Start a transaction
	tx, err := suite.db.BeginTx(suite.ctx, nil)
	require.NoError(suite.T(), err)

	// Note: This test demonstrates the concept, but our current repository doesn't support transactions directly
	// In a real implementation, you might have a repository that accepts a transaction or context with transaction

	// Create user within transaction
	testUser, err := user.NewUser("tx@example.com", "Transaction User")
	require.NoError(suite.T(), err)

	// Insert directly using transaction (simulating transactional repository)
	_, err = tx.ExecContext(suite.ctx, `
		INSERT INTO users (id, email, name, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5)`,
		testUser.ID().String(),
		testUser.Email().String(),
		testUser.Name().String(),
		testUser.CreatedAt(),
		testUser.UpdatedAt(),
	)
	require.NoError(suite.T(), err)

	// Rollback transaction
	err = tx.Rollback()
	require.NoError(suite.T(), err)

	// User should not exist after rollback
	_, err = suite.repository.FindByID(suite.ctx, testUser.ID())
	assert.Equal(suite.T(), errors.ErrUserNotFound, err)
}

// TestRepositoryPerformance tests repository performance with larger datasets
func (suite *UserRepositoryTestSuite) TestRepositoryPerformance() {
	const numUsers = 100

	// Measure creation time
	start := time.Now()

	for i := 0; i < numUsers; i++ {
		testUser, err := user.NewUser(fmt.Sprintf("perf%d@example.com", i), fmt.Sprintf("Performance User %d", i))
		require.NoError(suite.T(), err)

		err = suite.repository.Create(suite.ctx, testUser)
		require.NoError(suite.T(), err)

		// Add small delay to avoid overwhelming the database
		if i%10 == 0 {
			time.Sleep(1 * time.Millisecond)
		}
	}

	creationTime := time.Since(start)
	suite.T().Logf("Created %d users in %v (avg: %v per user)", numUsers, creationTime, creationTime/numUsers)

	// Measure query time
	start = time.Now()
	users, _, err := suite.repository.FindAll(suite.ctx, numUsers, "")
	queryTime := time.Since(start)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, numUsers)
	suite.T().Logf("Queried %d users in %v", numUsers, queryTime)

	// Performance assertions (adjust thresholds based on your requirements)
	assert.Less(suite.T(), creationTime, 30*time.Second, "User creation took too long")
	assert.Less(suite.T(), queryTime, 1*time.Second, "User query took too long")
}

// TestRepositoryDataIntegrity tests data integrity constraints
func (suite *UserRepositoryTestSuite) TestRepositoryDataIntegrity() {
	// Test email uniqueness constraint
	user1, err := user.NewUser("unique@example.com", "User 1")
	require.NoError(suite.T(), err)

	err = suite.repository.Create(suite.ctx, user1)
	require.NoError(suite.T(), err)

	// Try to create another user with same email
	user2, err := user.NewUser("unique@example.com", "User 2")
	require.NoError(suite.T(), err)

	err = suite.repository.Create(suite.ctx, user2)
	assert.Equal(suite.T(), errors.ErrDuplicateEmail, err)

	// Test that partial updates don't violate constraints
	user3, err := user.NewUser("another@example.com", "User 3")
	require.NoError(suite.T(), err)

	err = suite.repository.Create(suite.ctx, user3)
	require.NoError(suite.T(), err)

	// Try to update user3 to have same email as user1
	err = user3.UpdateEmail("unique@example.com")
	require.NoError(suite.T(), err)

	err = suite.repository.Update(suite.ctx, user3)
	assert.Equal(suite.T(), errors.ErrDuplicateEmail, err)
}

// TestRepositoryEdgeCases tests edge cases and boundary conditions
func (suite *UserRepositoryTestSuite) TestRepositoryEdgeCases() {
	// Test with very long email (within limits)
	longEmail := strings.Repeat("a", 240) + "@example.com" // 251 chars total
	longUser, err := user.NewUser(longEmail, "Long Email User")
	require.NoError(suite.T(), err)

	err = suite.repository.Create(suite.ctx, longUser)
	assert.NoError(suite.T(), err)

	// Test with very long name (within limits)
	longName := strings.Repeat("A", 250) // 250 chars
	longNameUser, err := user.NewUser("longname@example.com", longName)
	require.NoError(suite.T(), err)

	err = suite.repository.Create(suite.ctx, longNameUser)
	assert.NoError(suite.T(), err)

	// Test pagination with edge cases
	// Empty result set
	users, cursor, err := suite.repository.FindAll(suite.ctx, 10, "nonexistent-cursor")
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), users)
	assert.Empty(suite.T(), cursor)

	// Zero limit (should handle gracefully)
	users, cursor, err = suite.repository.FindAll(suite.ctx, 0, "")
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), users)
	assert.Empty(suite.T(), cursor)
}

// TestUserRepositoryIntegration runs the integration test suite
func TestUserRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
