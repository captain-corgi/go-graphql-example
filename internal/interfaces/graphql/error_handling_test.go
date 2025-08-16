package graphql_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/captain-corgi/go-graphql-example/internal/application/user/mocks"
	"github.com/captain-corgi/go-graphql-example/internal/interfaces/graphql/generated"
	"github.com/captain-corgi/go-graphql-example/internal/interfaces/graphql/resolver"
)

// TestGraphQLErrorHandling tests various error scenarios
func TestGraphQLErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	resolver := resolver.NewResolver(mockUserService, logger)

	testServer := createTestServer(resolver)
	defer testServer.Close()

	t.Run("Invalid GraphQL Syntax", func(t *testing.T) {
		invalidQuery := `
			query {
				user(id: "test") {
					id
					email
					// Missing closing brace
		`

		response := executeGraphQLRequest(t, testServer.URL, invalidQuery, nil)

		// Verify syntax error
		assert.NotNil(t, response.Errors)
		assert.Len(t, response.Errors, 1)
		// GraphQL syntax errors can have different messages
		errorMsg := strings.ToLower(response.Errors[0].Message)
		assert.True(t,
			strings.Contains(errorMsg, "syntax") ||
				strings.Contains(errorMsg, "expected") ||
				strings.Contains(errorMsg, "invalid") ||
				strings.Contains(errorMsg, "parse"),
			"Expected syntax error, got: %s", response.Errors[0].Message)
	})

	t.Run("Unknown Field", func(t *testing.T) {
		query := `
			query {
				user(id: "test") {
					id
					email
					unknownField
				}
			}
		`

		response := executeGraphQLRequest(t, testServer.URL, query, nil)

		// Verify field error
		assert.NotNil(t, response.Errors)
		assert.Len(t, response.Errors, 1)
		assert.Contains(t, strings.ToLower(response.Errors[0].Message), "field")
	})

	t.Run("Missing Required Argument", func(t *testing.T) {
		query := `
			query {
				user {
					id
					email
				}
			}
		`

		response := executeGraphQLRequest(t, testServer.URL, query, nil)

		// Verify argument error
		assert.NotNil(t, response.Errors)
		assert.Len(t, response.Errors, 1)
		assert.Contains(t, strings.ToLower(response.Errors[0].Message), "argument")
	})

	t.Run("Invalid JSON Request", func(t *testing.T) {
		invalidJSON := `{"query": "{ user(id: "test") { id } }", "variables": invalid}`

		resp, err := http.Post(
			fmt.Sprintf("%s/query", testServer.URL),
			"application/json",
			bytes.NewBuffer([]byte(invalidJSON)),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return a 400 Bad Request for invalid JSON
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Empty Request Body", func(t *testing.T) {
		resp, err := http.Post(
			fmt.Sprintf("%s/query", testServer.URL),
			"application/json",
			bytes.NewBuffer([]byte("")),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return a 400 Bad Request for empty body
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Missing Required Variable", func(t *testing.T) {
		query := `
			query GetUser($id: ID!) {
				user(id: $id) {
					id
					email
				}
			}
		`

		// No variables provided
		response := executeGraphQLRequest(t, testServer.URL, query, nil)

		// Should return validation error for missing required variable
		assert.NotNil(t, response.Errors)
		assert.Len(t, response.Errors, 1)
		errorMsg := strings.ToLower(response.Errors[0].Message)
		assert.True(t,
			strings.Contains(errorMsg, "variable") ||
				strings.Contains(errorMsg, "must be defined") ||
				strings.Contains(errorMsg, "required"),
			"Expected variable error, got: %s", response.Errors[0].Message)
	})
}

// TestGraphQLIntrospection tests GraphQL introspection queries
func TestGraphQLIntrospection(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	resolver := resolver.NewResolver(mockUserService, logger)

	// Create a direct GraphQL handler for introspection testing
	schema := generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	})
	graphqlHandler := handler.NewDefaultServer(schema)

	t.Run("Schema Introspection", func(t *testing.T) {
		query := `
			query IntrospectionQuery {
				__schema {
					types {
						name
						kind
					}
				}
			}
		`

		req := httptest.NewRequest("POST", "/query", bytes.NewBuffer([]byte(fmt.Sprintf(`{"query": %q}`, query))))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		graphqlHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response GraphQLResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Nil(t, response.Errors)
		assert.NotNil(t, response.Data)

		// Verify that our custom types are present in the schema
		schemaData := response.Data.(map[string]interface{})["__schema"].(map[string]interface{})
		types := schemaData["types"].([]interface{})

		typeNames := make([]string, 0)
		for _, typeInterface := range types {
			typeMap := typeInterface.(map[string]interface{})
			typeNames = append(typeNames, typeMap["name"].(string))
		}

		// Check for our custom types
		assert.Contains(t, typeNames, "User")
		assert.Contains(t, typeNames, "UserConnection")
		assert.Contains(t, typeNames, "CreateUserInput")
		assert.Contains(t, typeNames, "CreateUserPayload")
		assert.Contains(t, typeNames, "Query")
		assert.Contains(t, typeNames, "Mutation")
	})

	t.Run("Type Introspection", func(t *testing.T) {
		query := `
			query TypeIntrospection {
				__type(name: "User") {
					name
					kind
					fields {
						name
						type {
							name
							kind
						}
					}
				}
			}
		`

		req := httptest.NewRequest("POST", "/query", bytes.NewBuffer([]byte(fmt.Sprintf(`{"query": %q}`, query))))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		graphqlHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response GraphQLResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Nil(t, response.Errors)
		assert.NotNil(t, response.Data)

		// Verify User type structure
		typeData := response.Data.(map[string]interface{})["__type"].(map[string]interface{})
		assert.Equal(t, "User", typeData["name"])
		assert.Equal(t, "OBJECT", typeData["kind"])

		fields := typeData["fields"].([]interface{})
		fieldNames := make([]string, 0)
		for _, fieldInterface := range fields {
			fieldMap := fieldInterface.(map[string]interface{})
			fieldNames = append(fieldNames, fieldMap["name"].(string))
		}

		// Check for User fields
		assert.Contains(t, fieldNames, "id")
		assert.Contains(t, fieldNames, "email")
		assert.Contains(t, fieldNames, "name")
		assert.Contains(t, fieldNames, "createdAt")
		assert.Contains(t, fieldNames, "updatedAt")
	})
}
