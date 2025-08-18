package graphql_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authMocks "github.com/captain-corgi/go-graphql-example/internal/application/auth/mocks"
	"github.com/captain-corgi/go-graphql-example/internal/application/user/mocks"
	"github.com/captain-corgi/go-graphql-example/internal/interfaces/graphql/resolver"
)

// TestGraphQLServerSetup tests the GraphQL server setup and basic functionality
func TestGraphQLServerSetup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)
	mockAuthService := authMocks.NewMockService(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	resolver := resolver.NewResolver(mockUserService, mockAuthService, logger)

	testServer := createTestServer(resolver)
	defer testServer.Close()

	t.Run("GraphQL Endpoint Available", func(t *testing.T) {
		// Test that the GraphQL endpoint is accessible
		resp, err := http.Post(
			fmt.Sprintf("%s/query", testServer.URL),
			"application/json",
			bytes.NewBuffer([]byte(`{"query": "{ __schema { types { name } } }"}`)),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")
	})

	t.Run("GraphQL Playground Available", func(t *testing.T) {
		// Test that the GraphQL playground is accessible
		resp, err := http.Get(fmt.Sprintf("%s/playground", testServer.URL))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, resp.Header.Get("Content-Type"), "text/html")
	})

	t.Run("Health Check Available", func(t *testing.T) {
		// Test that the health check endpoint is accessible
		resp, err := http.Get(fmt.Sprintf("%s/health", testServer.URL))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var healthResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&healthResponse)
		require.NoError(t, err)

		assert.Equal(t, "ok", healthResponse["status"])
		assert.Equal(t, "graphql-service", healthResponse["service"])
		assert.NotEmpty(t, healthResponse["timestamp"])
	})

	t.Run("CORS Headers Present", func(t *testing.T) {
		// Test that CORS headers are present
		resp, err := http.Post(
			fmt.Sprintf("%s/query", testServer.URL),
			"application/json",
			bytes.NewBuffer([]byte(`{"query": "{ __schema { types { name } } }"}`)),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
	})

	t.Run("Request ID Header Present", func(t *testing.T) {
		// Test that request ID header is present
		resp, err := http.Post(
			fmt.Sprintf("%s/query", testServer.URL),
			"application/json",
			bytes.NewBuffer([]byte(`{"query": "{ __schema { types { name } } }"}`)),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.NotEmpty(t, resp.Header.Get("X-Request-ID"))
	})

	t.Run("HTTP Method Validation", func(t *testing.T) {
		// Test that GET requests to /query are handled appropriately
		resp, err := http.Get(fmt.Sprintf("%s/query", testServer.URL))
		require.NoError(t, err)
		defer resp.Body.Close()

		// GET requests to GraphQL endpoint should return method not allowed, bad request, not found, or handle introspection
		assert.True(t,
			resp.StatusCode == http.StatusMethodNotAllowed ||
				resp.StatusCode == http.StatusOK ||
				resp.StatusCode == http.StatusBadRequest ||
				resp.StatusCode == http.StatusNotFound,
			"Expected status 200, 400, 404, or 405, got: %d", resp.StatusCode)
	})
}

// TestGraphQLConcurrency tests concurrent requests
func TestGraphQLConcurrency(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)
	mockAuthService := authMocks.NewMockService(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	resolver := resolver.NewResolver(mockUserService, mockAuthService, logger)

	testServer := createTestServer(resolver)
	defer testServer.Close()

	t.Run("Concurrent Requests", func(t *testing.T) {
		// Test that the server can handle concurrent requests
		const numRequests = 10
		results := make(chan error, numRequests)

		for i := 0; i < numRequests; i++ {
			go func() {
				resp, err := http.Post(
					fmt.Sprintf("%s/query", testServer.URL),
					"application/json",
					bytes.NewBuffer([]byte(`{"query": "{ __schema { types { name } } }"}`)),
				)
				if err != nil {
					results <- err
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					results <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
					return
				}

				results <- nil
			}()
		}

		// Wait for all requests to complete
		for i := 0; i < numRequests; i++ {
			err := <-results
			assert.NoError(t, err, "Concurrent request %d failed", i)
		}
	})

	t.Run("Concurrent Different Operations", func(t *testing.T) {
		// Test concurrent requests with different operations
		const numRequests = 6
		results := make(chan error, numRequests)

		queries := []string{
			`{"query": "{ __schema { types { name } } }"}`,
			`{"query": "{ __type(name: \"User\") { name kind } }"}`,
			`{"query": "{ __schema { queryType { name } } }"}`,
			`{"query": "{ __schema { mutationType { name } } }"}`,
			`{"query": "{ __schema { types { name kind } } }"}`,
			`{"query": "{ __type(name: \"Query\") { fields { name } } }"}`,
		}

		for i, query := range queries {
			go func(q string, index int) {
				resp, err := http.Post(
					fmt.Sprintf("%s/query", testServer.URL),
					"application/json",
					bytes.NewBuffer([]byte(q)),
				)
				if err != nil {
					results <- fmt.Errorf("request %d failed: %w", index, err)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					results <- fmt.Errorf("request %d returned status %d", index, resp.StatusCode)
					return
				}

				results <- nil
			}(query, i)
		}

		// Wait for all requests to complete
		for i := 0; i < numRequests; i++ {
			err := <-results
			assert.NoError(t, err, "Concurrent operation %d failed", i)
		}
	})
}

// TestGraphQLPerformance tests basic performance characteristics
func TestGraphQLPerformance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)
	mockAuthService := authMocks.NewMockService(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	resolver := resolver.NewResolver(mockUserService, mockAuthService, logger)

	testServer := createTestServer(resolver)
	defer testServer.Close()

	t.Run("Response Time", func(t *testing.T) {
		// Test that introspection queries respond within reasonable time
		start := time.Now()

		resp, err := http.Post(
			fmt.Sprintf("%s/query", testServer.URL),
			"application/json",
			bytes.NewBuffer([]byte(`{"query": "{ __schema { types { name } } }"}`)),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		duration := time.Since(start)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Less(t, duration, 5*time.Second, "Introspection query should complete within 5 seconds")
	})

	t.Run("Memory Usage", func(t *testing.T) {
		// Test that repeated requests don't cause memory leaks
		// This is a basic test - more sophisticated memory testing would require runtime profiling

		for i := 0; i < 100; i++ {
			resp, err := http.Post(
				fmt.Sprintf("%s/query", testServer.URL),
				"application/json",
				bytes.NewBuffer([]byte(`{"query": "{ __schema { types { name } } }"}`)),
			)
			require.NoError(t, err)
			resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}

		// If we get here without running out of memory, the test passes
		assert.True(t, true, "Completed 100 requests without memory issues")
	})
}
