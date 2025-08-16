# Domain User Mocks

This directory contains generated mocks and test utilities for the domain user interfaces.

## Generated Mocks

- `mock_repository.go` - Mock implementation of `user.Repository` interface
- `mock_service.go` - Mock implementation of `user.DomainService` interface

## Test Utilities

- `test_utils.go` - Helper functions and builders for common test scenarios
- `mock_generation_test.go` - Tests to verify mock generation is working correctly

## Usage

### Basic Mock Usage

```go
func TestSomething(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockRepository(ctrl)
    
    // Set up expectations
    userID, _ := user.NewUserID("test-id")
    testUser, _ := user.NewUser("test@example.com", "Test User")
    
    mockRepo.EXPECT().
        FindByID(gomock.Any(), userID).
        Return(testUser, nil).
        Times(1)
    
    // Use the mock in your test
    result, err := mockRepo.FindByID(context.Background(), userID)
    // ... assertions
}
```

### Using Test Utilities

```go
func TestWithUtilities(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockRepository(ctrl)
    utils := mocks.NewRepositoryTestUtils(mockRepo)
    
    // Use utility methods for common scenarios
    userID, _ := user.NewUserID("test-id")
    utils.ExpectFindByIDNotFound(userID)
    
    // Test your code
    _, err := mockRepo.FindByID(context.Background(), userID)
    assert.Equal(t, errors.ErrUserNotFound, err)
}
```

### Building Test Data

```go
func TestWithBuilder(t *testing.T) {
    // Create test users with builder pattern
    testUser := mocks.NewTestUserBuilder().
        WithID("custom-id").
        WithEmail("custom@example.com").
        WithName("Custom User").
        Build()
    
    // Use in your tests
    assert.Equal(t, "custom-id", testUser.ID().String())
}
```

## Regenerating Mocks

To regenerate mocks after interface changes:

```bash
# Regenerate all mocks in this package
go generate ./internal/domain/user/...

# Or regenerate all mocks in the project
go generate ./...
```

## Available Test Utilities

### RepositoryTestUtils

- `ExpectFindByIDSuccess(userID, returnUser)` - Mock successful user lookup
- `ExpectFindByIDNotFound(userID)` - Mock user not found
- `ExpectFindByEmailSuccess(email, returnUser)` - Mock successful email lookup
- `ExpectFindByEmailNotFound(email)` - Mock email not found
- `ExpectCreateSuccess(user)` - Mock successful user creation
- `ExpectUpdateSuccess(user)` - Mock successful user update
- `ExpectDeleteSuccess(userID)` - Mock successful user deletion
- `ExpectExistsByEmailTrue(email)` - Mock email exists check (true)
- `ExpectExistsByEmailFalse(email)` - Mock email exists check (false)
- `ExpectFindAllSuccess(limit, cursor, users, nextCursor)` - Mock successful user listing

### DomainServiceTestUtils

- `ExpectValidateUniqueEmailSuccess(email, excludeUserID)` - Mock successful email validation
- `ExpectValidateUniqueEmailDuplicate(email, excludeUserID)` - Mock duplicate email error
- `ExpectCanDeleteUserSuccess(userID)` - Mock successful deletion validation

### TestUserBuilder

- `WithID(id)` - Set user ID
- `WithEmail(email)` - Set user email
- `WithName(name)` - Set user name
- `WithCreatedAt(time)` - Set creation timestamp
- `WithUpdatedAt(time)` - Set update timestamp
- `Build()` - Create the user instance

## Common Test Context

Use `CommonTestContext()` to get a context with a reasonable timeout for tests.
