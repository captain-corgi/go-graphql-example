# Application User Mocks

This directory contains generated mocks and test utilities for the application user service interface.

## Generated Mocks

- `mock_service.go` - Mock implementation of `user.Service` interface

## Test Utilities

- `test_utils.go` - Helper functions and builders for common test scenarios
- `mock_generation_test.go` - Tests to verify mock generation is working correctly

## Usage

### Basic Mock Usage

```go
func TestSomething(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockService := mocks.NewMockService(ctrl)
    
    // Set up expectations
    req := user.GetUserRequest{ID: "test-id"}
    expectedResp := &user.GetUserResponse{
        User: &user.UserDTO{
            ID:    "test-id",
            Email: "test@example.com",
            Name:  "Test User",
        },
    }
    
    mockService.EXPECT().
        GetUser(gomock.Any(), req).
        Return(expectedResp, nil).
        Times(1)
    
    // Use the mock in your test
    result, err := mockService.GetUser(context.Background(), req)
    // ... assertions
}
```

### Using Test Utilities

```go
func TestWithUtilities(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockService := mocks.NewMockService(ctrl)
    utils := mocks.NewServiceTestUtils(mockService)
    
    // Use utility methods for common scenarios
    testUser := &user.UserDTO{
        ID:    "test-id",
        Email: "test@example.com",
        Name:  "Test User",
    }
    
    utils.ExpectGetUserSuccess("test-id", testUser)
    
    // Test your code
    req := user.GetUserRequest{ID: "test-id"}
    resp, err := mockService.GetUser(context.Background(), req)
    assert.NoError(t, err)
    assert.Equal(t, "test-id", resp.User.ID)
}
```

### Building Test Data

```go
func TestWithBuilder(t *testing.T) {
    // Create test DTOs with builder pattern
    userDTO := mocks.NewTestDTOBuilder().
        WithID("custom-id").
        WithEmail("custom@example.com").
        WithName("Custom User").
        BuildUserDTO()
    
    // Create user connections for pagination tests
    users := []*user.UserDTO{userDTO}
    connection := mocks.BuildTestUserConnection(users, true, false)
    
    // Use in your tests
    assert.Equal(t, 1, len(connection.Edges))
    assert.True(t, connection.PageInfo.HasNextPage)
}
```

## Regenerating Mocks

To regenerate mocks after interface changes:

```bash
# Regenerate all mocks in this package
go generate ./internal/application/user/...

# Or regenerate all mocks in the project
go generate ./...
```

## Available Test Utilities

### ServiceTestUtils

- `ExpectGetUserSuccess(userID, returnUser)` - Mock successful user retrieval
- `ExpectGetUserNotFound(userID)` - Mock user not found error
- `ExpectListUsersSuccess(first, after, returnUsers)` - Mock successful user listing
- `ExpectCreateUserSuccess(email, name, returnUser)` - Mock successful user creation
- `ExpectCreateUserDuplicateEmail(email, name)` - Mock duplicate email error
- `ExpectUpdateUserSuccess(userID, email, name, returnUser)` - Mock successful user update
- `ExpectUpdateUserNotFound(userID, email, name)` - Mock user not found for update
- `ExpectDeleteUserSuccess(userID)` - Mock successful user deletion
- `ExpectDeleteUserNotFound(userID)` - Mock user not found for deletion

### TestDTOBuilder

- `WithID(id)` - Set user ID
- `WithEmail(email)` - Set user email
- `WithName(name)` - Set user name
- `WithCreatedAt(time)` - Set creation timestamp
- `WithUpdatedAt(time)` - Set update timestamp
- `BuildUserDTO()` - Create UserDTO instance
- `BuildUserEdgeDTO()` - Create UserEdgeDTO instance

### Helper Functions

- `BuildTestUserConnection(users, hasNextPage, hasPreviousPage)` - Create test user connection for pagination
- `CommonTestContext()` - Get context with reasonable timeout for tests

## Testing Patterns

### Testing GraphQL Resolvers

```go
func TestUserResolver(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockService := mocks.NewMockService(ctrl)
    utils := mocks.NewServiceTestUtils(mockService)
    
    // Set up expectations
    testUser := mocks.NewTestDTOBuilder().
        WithID("test-id").
        BuildUserDTO()
    
    utils.ExpectGetUserSuccess("test-id", testUser)
    
    // Test resolver
    resolver := &Resolver{userService: mockService}
    result, err := resolver.User(context.Background(), "test-id")
    
    assert.NoError(t, err)
    assert.Equal(t, "test-id", result.ID)
}
```

### Testing Pagination

```go
func TestPagination(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockService := mocks.NewMockService(ctrl)
    utils := mocks.NewServiceTestUtils(mockService)
    
    // Create test data
    users := []*user.UserDTO{
        mocks.NewTestDTOBuilder().WithID("1").BuildUserDTO(),
        mocks.NewTestDTOBuilder().WithID("2").BuildUserDTO(),
    }
    
    connection := mocks.BuildTestUserConnection(users, true, false)
    utils.ExpectListUsersSuccess(10, "", connection)
    
    // Test your pagination logic
    // ...
}
```
