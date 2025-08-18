package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/golang-jwt/jwt/v5"
)

// JWTService handles JWT token operations
type JWTService struct {
	secretKey       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	issuer          string
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWT service
func NewJWTService(secretKey string, accessTokenTTL, refreshTokenTTL time.Duration, issuer string) *JWTService {
	return &JWTService{
		secretKey:       []byte(secretKey),
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		issuer:          issuer,
	}
}

// GenerateTokenPair generates a new access and refresh token pair
func (j *JWTService) GenerateTokenPair(userID, email string) (*TokenPair, error) {
	now := time.Now()

	// Generate access token
	accessTokenClaims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   userID,
			Audience:  []string{"api"},
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(j.secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token (random string)
	refreshToken, err := j.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
		ExpiresAt:    now.Add(j.accessTokenTTL),
	}, nil
}

// ValidateAccessToken validates and parses an access token
func (j *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, errors.DomainError{
			Code:    "INVALID_TOKEN",
			Message: "Invalid or expired access token",
			Field:   "accessToken",
		}
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.DomainError{
			Code:    "INVALID_TOKEN",
			Message: "Invalid token claims",
			Field:   "accessToken",
		}
	}

	// Check if token is expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.DomainError{
			Code:    "TOKEN_EXPIRED",
			Message: "Access token has expired",
			Field:   "accessToken",
		}
	}

	return claims, nil
}

// generateRefreshToken generates a cryptographically secure random string
func (j *JWTService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetRefreshTokenTTL returns the refresh token TTL
func (j *JWTService) GetRefreshTokenTTL() time.Duration {
	return j.refreshTokenTTL
}

