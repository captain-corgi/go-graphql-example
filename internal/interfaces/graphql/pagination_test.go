package graphql_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/captain-corgi/go-graphql-example/internal/application/user"
	"github.com/captain-corgi/go-graphql-example/internal/application/user/mocks"
	"github.com/captain-corgi/go-graphql-example/internal/interfaces/graphql/resolver"
)

// TestGraphQLPaginationEdgeCases tests pagination edge cases
func TestGraphQLPaginationEdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	resolver := resolver.NewResolver(mockUserService, logger)

	testServer := createTestServer(resolver)
	defer testServer.Close()

	t.Run("Empty Result Set", func(t *testing.T) {
		emptyConnection := &user.UserConnectionDTO{
			Edges: []*user.UserEdgeDTO{},
			PageInfo: &user.PageInfoDTO{
				HasNextPage:     false,
				HasPreviousPage: false,
				StartCursor:     nil,
				EndCursor:       nil,
			},
		}

		mockUserService.EXPECT().
			ListUsers(gomock.Any(), user.ListUsersRequest{First: 10, After: ""}).
			Return(&user.ListUsersResponse{Users: emptyConnection}, nil)

		query := `
			query ListUsers {
				users {
					edges {
						node {
							id
						}
					}
					pageInfo {
						hasNextPage
						hasPreviousPage
						startCursor
						endCursor
					}
				}
			}
		`

		response := executeGraphQLRequest(t, testServer.URL, query, nil)

		// Verify response
		assert.Nil(t, response.Errors)
		assert.NotNil(t, response.Data)

		usersData := response.Data.(map[string]interface{})["users"].(map[string]interface{})
		edges := usersData["edges"].([]interface{})
		assert.Len(t, edges, 0)

		pageInfo := usersData["pageInfo"].(map[string]interface{})
		assert.False(t, pageInfo["hasNextPage"].(bool))
		assert.False(t, pageInfo["hasPreviousPage"].(bool))
		assert.Nil(t, pageInfo["startCursor"])
		assert.Nil(t, pageInfo["endCursor"])
	})

	t.Run("Maximum Page Size", func(t *testing.T) {
		query := `
			query ListUsers($first: Int) {
				users(first: $first) {
					edges {
						node {
							id
						}
					}
				}
			}
		`

		variables := map[string]interface{}{
			"first": 101, // Exceeds maximum of 100
		}

		response := executeGraphQLRequest(t, testServer.URL, query, variables)

		// Verify validation error
		assert.NotNil(t, response.Errors)
		assert.Len(t, response.Errors, 1)
		assert.Contains(t, response.Errors[0].Message, "First parameter cannot exceed 100")
	})

	t.Run("Negative Page Size", func(t *testing.T) {
		query := `
			query ListUsers($first: Int) {
				users(first: $first) {
					edges {
						node {
							id
						}
					}
				}
			}
		`

		variables := map[string]interface{}{
			"first": -5, // Negative value
		}

		response := executeGraphQLRequest(t, testServer.URL, query, variables)

		// Verify validation error
		assert.NotNil(t, response.Errors)
		assert.Len(t, response.Errors, 1)
		assert.Contains(t, response.Errors[0].Message, "First parameter must be non-negative")
	})

	t.Run("Pagination with After Cursor", func(t *testing.T) {
		users := []*user.UserEdgeDTO{
			{
				Node: &user.UserDTO{
					ID:        "user-3",
					Email:     "user3@example.com",
					Name:      "User Three",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Cursor: "user-3",
			},
			{
				Node: &user.UserDTO{
					ID:        "user-4",
					Email:     "user4@example.com",
					Name:      "User Four",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Cursor: "user-4",
			},
		}

		connection := &user.UserConnectionDTO{
			Edges: users,
			PageInfo: &user.PageInfoDTO{
				HasNextPage:     true,
				HasPreviousPage: true, // Has previous page since we're using after cursor
				StartCursor:     &users[0].Cursor,
				EndCursor:       &users[1].Cursor,
			},
		}

		mockUserService.EXPECT().
			ListUsers(gomock.Any(), user.ListUsersRequest{First: 2, After: "user-2"}).
			Return(&user.ListUsersResponse{Users: connection}, nil)

		query := `
			query ListUsers($first: Int, $after: String) {
				users(first: $first, after: $after) {
					edges {
						node {
							id
							email
						}
						cursor
					}
					pageInfo {
						hasNextPage
						hasPreviousPage
						startCursor
						endCursor
					}
				}
			}
		`

		variables := map[string]interface{}{
			"first": 2,
			"after": "user-2",
		}

		response := executeGraphQLRequest(t, testServer.URL, query, variables)

		// Verify response
		assert.Nil(t, response.Errors)
		assert.NotNil(t, response.Data)

		usersData := response.Data.(map[string]interface{})["users"].(map[string]interface{})
		edges := usersData["edges"].([]interface{})
		assert.Len(t, edges, 2)

		// Verify users
		firstEdge := edges[0].(map[string]interface{})
		firstNode := firstEdge["node"].(map[string]interface{})
		assert.Equal(t, "user-3", firstNode["id"])

		secondEdge := edges[1].(map[string]interface{})
		secondNode := secondEdge["node"].(map[string]interface{})
		assert.Equal(t, "user-4", secondNode["id"])

		// Verify pagination info
		pageInfo := usersData["pageInfo"].(map[string]interface{})
		assert.True(t, pageInfo["hasNextPage"].(bool))
		assert.True(t, pageInfo["hasPreviousPage"].(bool))
		assert.Equal(t, "user-3", pageInfo["startCursor"])
		assert.Equal(t, "user-4", pageInfo["endCursor"])
	})

	t.Run("Empty After Cursor", func(t *testing.T) {
		query := `
			query ListUsers($after: String) {
				users(after: $after) {
					edges {
						node {
							id
						}
					}
				}
			}
		`

		variables := map[string]interface{}{
			"after": "", // Empty cursor
		}

		response := executeGraphQLRequest(t, testServer.URL, query, variables)

		// Verify validation error for empty cursor
		assert.NotNil(t, response.Errors)
		assert.Len(t, response.Errors, 1)
		assert.Contains(t, response.Errors[0].Message, "Cursor cannot be empty")
	})

	t.Run("Large Page Size Within Limit", func(t *testing.T) {
		// Create 50 users for the response
		users := make([]*user.UserEdgeDTO, 50)
		for i := 0; i < 50; i++ {
			users[i] = &user.UserEdgeDTO{
				Node: &user.UserDTO{
					ID:        fmt.Sprintf("user-%d", i+1),
					Email:     fmt.Sprintf("user%d@example.com", i+1),
					Name:      fmt.Sprintf("User %d", i+1),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Cursor: fmt.Sprintf("user-%d", i+1),
			}
		}

		connection := &user.UserConnectionDTO{
			Edges: users,
			PageInfo: &user.PageInfoDTO{
				HasNextPage:     false,
				HasPreviousPage: false,
				StartCursor:     &users[0].Cursor,
				EndCursor:       &users[49].Cursor,
			},
		}

		mockUserService.EXPECT().
			ListUsers(gomock.Any(), user.ListUsersRequest{First: 50, After: ""}).
			Return(&user.ListUsersResponse{Users: connection}, nil)

		query := `
			query ListUsers($first: Int) {
				users(first: $first) {
					edges {
						node {
							id
						}
					}
					pageInfo {
						hasNextPage
						hasPreviousPage
					}
				}
			}
		`

		variables := map[string]interface{}{
			"first": 50, // Large but valid page size
		}

		response := executeGraphQLRequest(t, testServer.URL, query, variables)

		// Verify successful response
		assert.Nil(t, response.Errors)
		assert.NotNil(t, response.Data)

		usersData := response.Data.(map[string]interface{})["users"].(map[string]interface{})
		edges := usersData["edges"].([]interface{})
		assert.Len(t, edges, 50)
	})
}

// TestGraphQLInputSanitization tests input sanitization
func TestGraphQLInputSanitization(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	resolver := resolver.NewResolver(mockUserService, logger)

	testServer := createTestServer(resolver)
	defer testServer.Close()

	t.Run("Whitespace Trimming in User ID", func(t *testing.T) {
		expectedUser := &user.UserDTO{
			ID:        "user-123",
			Email:     "test@example.com",
			Name:      "Test User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Expect the service to be called with trimmed ID
		mockUserService.EXPECT().
			GetUser(gomock.Any(), user.GetUserRequest{ID: "user-123"}).
			Return(&user.GetUserResponse{User: expectedUser}, nil)

		query := `
			query GetUser($id: ID!) {
				user(id: $id) {
					id
					email
					name
				}
			}
		`

		variables := map[string]interface{}{
			"id": "  user-123  ", // ID with extra whitespace
		}

		response := executeGraphQLRequest(t, testServer.URL, query, variables)

		// Verify successful response (whitespace should be trimmed)
		assert.Nil(t, response.Errors)
		assert.NotNil(t, response.Data)

		userData := response.Data.(map[string]interface{})["user"].(map[string]interface{})
		assert.Equal(t, expectedUser.ID, userData["id"])
	})

	t.Run("Whitespace Trimming in Create User Input", func(t *testing.T) {
		expectedUser := &user.UserDTO{
			ID:        "new-user-id",
			Email:     "new@example.com",
			Name:      "New User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Expect the service to be called with trimmed values
		mockUserService.EXPECT().
			CreateUser(gomock.Any(), user.CreateUserRequest{
				Email: "new@example.com",
				Name:  "New User",
			}).
			Return(&user.CreateUserResponse{User: expectedUser}, nil)

		mutation := `
			mutation CreateUser($input: CreateUserInput!) {
				createUser(input: $input) {
					user {
						id
						email
						name
					}
					errors {
						message
					}
				}
			}
		`

		variables := map[string]interface{}{
			"input": map[string]interface{}{
				"email": "  new@example.com  ", // Email with extra whitespace
				"name":  "  New User  ",        // Name with extra whitespace
			},
		}

		response := executeGraphQLRequest(t, testServer.URL, mutation, variables)

		// Verify successful response (whitespace should be trimmed)
		assert.Nil(t, response.Errors)
		assert.NotNil(t, response.Data)

		createUserData := response.Data.(map[string]interface{})["createUser"].(map[string]interface{})
		userData := createUserData["user"].(map[string]interface{})
		assert.Equal(t, expectedUser.ID, userData["id"])
		assert.Equal(t, expectedUser.Email, userData["email"])
		assert.Equal(t, expectedUser.Name, userData["name"])
	})
}
