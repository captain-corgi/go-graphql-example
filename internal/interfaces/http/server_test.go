package http

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/captain-corgi/go-graphql-example/internal/application/user/mocks"
	"github.com/captain-corgi/go-graphql-example/internal/infrastructure/config"
	"github.com/captain-corgi/go-graphql-example/internal/interfaces/graphql/resolver"
)

// createTestResolver creates a resolver with mock dependencies for testing
func createTestResolver(t *testing.T) *resolver.Resolver {
	ctrl := gomock.NewController(t)
	mockUserService := mocks.NewMockService(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return resolver.NewResolver(mockUserService, logger)
}

func TestNewServer(t *testing.T) {
	cfg := &config.ServerConfig{
		Port:         "8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	resolver := createTestResolver(t)

	server := NewServer(cfg, resolver, logger)

	assert.NotNil(t, server)
	assert.Equal(t, cfg, server.config)
	assert.Equal(t, logger, server.logger)
	assert.Equal(t, resolver, server.resolver)
	assert.NotNil(t, server.router)
	assert.NotNil(t, server.server)
}

func TestServerRoutes(t *testing.T) {
	cfg := &config.ServerConfig{
		Port:         "8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	resolver := createTestResolver(t)

	server := NewServer(cfg, resolver, logger)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		checkHeaders   map[string]string
	}{
		{
			name:           "health check endpoint",
			method:         http.MethodGet,
			path:           "/health",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GraphQL playground endpoint",
			method:         http.MethodGet,
			path:           "/playground",
			expectedStatus: http.StatusOK,
			checkHeaders: map[string]string{
				"Content-Type": "text/html; charset=UTF-8",
			},
		},
		{
			name:           "root redirects to playground",
			method:         http.MethodGet,
			path:           "/",
			expectedStatus: http.StatusOK,
			checkHeaders: map[string]string{
				"Content-Type": "text/html; charset=UTF-8",
			},
		},
		{
			name:           "GraphQL query endpoint accepts POST",
			method:         http.MethodPost,
			path:           "/query",
			expectedStatus: http.StatusBadRequest, // No valid GraphQL query provided
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			server.router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			for header, expectedValue := range tt.checkHeaders {
				assert.Contains(t, w.Header().Get(header), expectedValue, "Header %s", header)
			}
		})
	}
}

func TestHealthHandler(t *testing.T) {
	cfg := &config.ServerConfig{
		Port:         "8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	resolver := createTestResolver(t)

	server := NewServer(cfg, resolver, logger)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"status\":\"ok\"")
	assert.Contains(t, w.Body.String(), "\"service\":\"graphql-service\"")
	assert.Contains(t, w.Body.String(), "\"timestamp\":")
}

func TestServerMiddleware(t *testing.T) {
	cfg := &config.ServerConfig{
		Port:         "8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	resolver := createTestResolver(t)

	server := NewServer(cfg, resolver, logger)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	// Check that middleware is working
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"), "Request ID middleware should add request ID header")
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"), "CORS middleware should add CORS headers")
}

func TestServerConfiguration(t *testing.T) {
	cfg := &config.ServerConfig{
		Port:         "9090",
		ReadTimeout:  45 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  180 * time.Second,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	resolver := createTestResolver(t)

	server := NewServer(cfg, resolver, logger)

	require.NotNil(t, server.server)
	assert.Equal(t, ":9090", server.server.Addr)
	assert.Equal(t, 45*time.Second, server.server.ReadTimeout)
	assert.Equal(t, 60*time.Second, server.server.WriteTimeout)
	assert.Equal(t, 180*time.Second, server.server.IdleTimeout)
}
