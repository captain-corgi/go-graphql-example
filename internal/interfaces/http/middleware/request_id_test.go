package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		existingHeader string
		expectGenerate bool
	}{
		{
			name:           "generates new request ID when none provided",
			existingHeader: "",
			expectGenerate: true,
		},
		{
			name:           "uses existing request ID when provided",
			existingHeader: "existing-request-id",
			expectGenerate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := gin.New()
			router.Use(RequestID())

			var capturedRequestID string
			var capturedContextID string

			router.GET("/test", func(c *gin.Context) {
				capturedRequestID = c.Writer.Header().Get(RequestIDHeader)
				capturedContextID = GetRequestID(c.Request.Context())
				c.Status(http.StatusOK)
			})

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.existingHeader != "" {
				req.Header.Set(RequestIDHeader, tt.existingHeader)
			}

			// Execute
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, http.StatusOK, w.Code)

			if tt.expectGenerate {
				assert.NotEmpty(t, capturedRequestID)
				// When generating, the ID should be different from empty string
				assert.NotEmpty(t, capturedRequestID)
			} else {
				assert.Equal(t, tt.existingHeader, capturedRequestID)
			}

			assert.Equal(t, capturedRequestID, capturedContextID)
			assert.Equal(t, capturedRequestID, w.Header().Get(RequestIDHeader))
		})
	}
}

func TestGenerateRequestID(t *testing.T) {
	// Test that generateRequestID produces unique IDs
	id1 := generateRequestID()
	id2 := generateRequestID()

	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)

	// Test that IDs are hex strings (when crypto/rand works)
	assert.Regexp(t, "^[a-f0-9]+$", id1)
}

func TestGetRequestID(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "returns request ID from context",
			ctx:      context.WithValue(context.Background(), RequestIDKey, "test-id"),
			expected: "test-id",
		},
		{
			name:     "returns empty string when no request ID in context",
			ctx:      context.Background(),
			expected: "",
		},
		{
			name:     "returns empty string when request ID is wrong type",
			ctx:      context.WithValue(context.Background(), RequestIDKey, 123),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRequestID(tt.ctx)
			assert.Equal(t, tt.expected, result)
		})
	}
}
