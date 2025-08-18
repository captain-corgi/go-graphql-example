package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/captain-corgi/go-graphql-example/internal/application/auth"
	"github.com/gin-gonic/gin"
)

// AuthKey is the context key for authentication information
type AuthKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey AuthKey = "userID"
	// UserEmailKey is the context key for user email
	UserEmailKey AuthKey = "userEmail"
)

// AuthMiddleware provides JWT authentication middleware
type AuthMiddleware struct {
	authService auth.Service
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(authService auth.Service) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth middleware that requires a valid JWT token
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		token := extractTokenFromHeader(c.GetHeader("Authorization"))
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid authorization header"})
			c.Abort()
			return
		}

		claims, err := m.authService.ValidateAccessToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Add user info to context
		ctx := context.WithValue(c.Request.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	})
}

// OptionalAuth middleware that optionally validates JWT token
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		token := extractTokenFromHeader(c.GetHeader("Authorization"))
		if token != "" {
			claims, err := m.authService.ValidateAccessToken(c.Request.Context(), token)
			if err == nil {
				// Add user info to context if token is valid
				ctx := context.WithValue(c.Request.Context(), UserIDKey, claims.UserID)
				ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
				c.Request = c.Request.WithContext(ctx)
			}
		}

		c.Next()
	})
}

// extractTokenFromHeader extracts the JWT token from the Authorization header
func extractTokenFromHeader(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	// Check for "Bearer " prefix
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return ""
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// GetUserEmailFromContext retrieves the user email from the request context
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailKey).(string)
	return email, ok
}

