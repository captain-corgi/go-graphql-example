package config

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_WithDefaults(t *testing.T) {
	// Clear any existing environment variables
	clearEnvVars(t)

	// Reset viper to clean state
	viper.Reset()

	// Change to a temporary directory to avoid loading config files
	originalDir, err := os.Getwd()
	require.NoError(t, err)

	tempDir := t.TempDir()
	err = os.Chdir(tempDir)
	require.NoError(t, err)

	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Logf("Failed to restore working directory: %v", err)
		}
	}()

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Check defaults
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, 30*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, 30*time.Second, cfg.Server.WriteTimeout)
	assert.Equal(t, 120*time.Second, cfg.Server.IdleTimeout)

	assert.Equal(t, "postgres://user:password@localhost:5432/graphql_service?sslmode=disable", cfg.Database.URL)
	assert.Equal(t, 25, cfg.Database.MaxOpenConns)
	assert.Equal(t, 5, cfg.Database.MaxIdleConns)
	assert.Equal(t, 5*time.Minute, cfg.Database.ConnMaxLifetime)
	assert.Equal(t, 5*time.Minute, cfg.Database.ConnMaxIdleTime)

	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "json", cfg.Logging.Format)
}

func TestLoad_WithEnvironmentVariables(t *testing.T) {
	// Clear any existing environment variables
	clearEnvVars(t)

	// Set environment variables
	envVars := map[string]string{
		"GRAPHQL_SERVICE_SERVER_PORT":             "9090",
		"GRAPHQL_SERVICE_SERVER_READ_TIMEOUT":     "60s",
		"GRAPHQL_SERVICE_DATABASE_URL":            "postgres://test:test@localhost:5432/test_db",
		"GRAPHQL_SERVICE_DATABASE_MAX_OPEN_CONNS": "50",
		"GRAPHQL_SERVICE_DATABASE_MAX_IDLE_CONNS": "10",
		"GRAPHQL_SERVICE_LOGGING_LEVEL":           "debug",
		"GRAPHQL_SERVICE_LOGGING_FORMAT":          "text",
	}

	for key, value := range envVars {
		t.Setenv(key, value)
	}

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Check environment variable overrides
	assert.Equal(t, "9090", cfg.Server.Port)
	assert.Equal(t, 60*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, "postgres://test:test@localhost:5432/test_db", cfg.Database.URL)
	assert.Equal(t, 50, cfg.Database.MaxOpenConns)
	assert.Equal(t, 10, cfg.Database.MaxIdleConns)
	assert.Equal(t, "debug", cfg.Logging.Level)
	assert.Equal(t, "text", cfg.Logging.Format)
}

func TestLoad_ValidationFailure(t *testing.T) {
	// Clear any existing environment variables
	clearEnvVars(t)

	// Set invalid environment variables
	t.Setenv("GRAPHQL_SERVICE_SERVER_PORT", "")
	t.Setenv("GRAPHQL_SERVICE_LOGGING_LEVEL", "invalid")

	cfg, err := Load()
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "config validation failed")
}

func TestMustLoad_Success(t *testing.T) {
	// Clear any existing environment variables
	clearEnvVars(t)

	// Reset viper to clean state
	viper.Reset()

	// Change to a temporary directory to avoid loading config files
	originalDir, err := os.Getwd()
	require.NoError(t, err)

	tempDir := t.TempDir()
	err = os.Chdir(tempDir)
	require.NoError(t, err)

	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Logf("Failed to restore working directory: %v", err)
		}
	}()

	cfg := MustLoad()
	require.NotNil(t, cfg)
	assert.Equal(t, "8080", cfg.Server.Port)
}

func TestMustLoad_Panic(t *testing.T) {
	// Clear any existing environment variables
	clearEnvVars(t)

	// Set invalid environment variables that will cause validation to fail
	// Setting invalid log level should cause validation to fail
	t.Setenv("GRAPHQL_SERVICE_LOGGING_LEVEL", "invalid_level")

	assert.Panics(t, func() {
		MustLoad()
	})
}

// clearEnvVars clears all GRAPHQL_SERVICE environment variables
func clearEnvVars(t *testing.T) {
	envVars := []string{
		"GRAPHQL_SERVICE_SERVER_PORT",
		"GRAPHQL_SERVICE_SERVER_READ_TIMEOUT",
		"GRAPHQL_SERVICE_SERVER_WRITE_TIMEOUT",
		"GRAPHQL_SERVICE_SERVER_IDLE_TIMEOUT",
		"GRAPHQL_SERVICE_DATABASE_URL",
		"GRAPHQL_SERVICE_DATABASE_MAX_OPEN_CONNS",
		"GRAPHQL_SERVICE_DATABASE_MAX_IDLE_CONNS",
		"GRAPHQL_SERVICE_DATABASE_CONN_MAX_LIFETIME",
		"GRAPHQL_SERVICE_DATABASE_CONN_MAX_IDLE_TIME",
		"GRAPHQL_SERVICE_LOGGING_LEVEL",
		"GRAPHQL_SERVICE_LOGGING_FORMAT",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}
