package database

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
)

func TestMaskDatabaseURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "long URL",
			url:      "postgres://user:password@localhost:5432/dbname",
			expected: "postgres:/***/dbname",
		},
		{
			name:     "short URL",
			url:      "short",
			expected: "***",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskDatabaseURL(tt.url)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewConnection_InvalidURL(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	cfg := config.DatabaseConfig{
		URL:             "invalid://url",
		MaxOpenConns:    5,
		MaxIdleConns:    2,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}

	_, err := NewConnection(cfg, logger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping database")
}
