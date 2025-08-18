package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ConfigIntegrationTestSuite defines the integration test suite for configuration loading
type ConfigIntegrationTestSuite struct {
	suite.Suite
	tempDir    string
	configPath string
}

// SetupSuite sets up the test suite
func (suite *ConfigIntegrationTestSuite) SetupSuite() {
	// Create temporary directory for test config files
	tempDir, err := os.MkdirTemp("", "config_integration_test")
	require.NoError(suite.T(), err)

	suite.tempDir = tempDir
	suite.configPath = filepath.Join(tempDir, "config.yaml")
}

// TearDownSuite cleans up the test suite
func (suite *ConfigIntegrationTestSuite) TearDownSuite() {
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// SetupTest sets up each test
func (suite *ConfigIntegrationTestSuite) SetupTest() {
	// Clear environment variables before each test
	suite.clearEnvVars()

	// Remove config file if it exists
	if _, err := os.Stat(suite.configPath); err == nil {
		os.Remove(suite.configPath)
	}
}

// TestLoadConfigFromFile tests loading configuration from YAML file
func (suite *ConfigIntegrationTestSuite) TestLoadConfigFromFile() {
	// Create test config file
	configContent := `
server:
  name: "test-server"
  port: "9000"
  read_timeout: "45s"
  write_timeout: "45s"
  idle_timeout: "180s"

database:
  url: "postgres://testuser:testpass@localhost:5432/testdb"
  max_open_conns: 30
  max_idle_conns: 8
  conn_max_lifetime: "10m"
  conn_max_idle_time: "8m"

logging:
  level: "debug"
  format: "text"

auth:
  jwt_secret: "test-secret-key-32-chars-minimum-length"
  access_token_ttl: "15m"
  refresh_token_ttl: "24h"
`

	err := os.WriteFile(suite.configPath, []byte(configContent), 0644)
	require.NoError(suite.T(), err)

	// Change working directory to temp dir so config file can be found
	originalWd, err := os.Getwd()
	require.NoError(suite.T(), err)
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			suite.T().Logf("Failed to restore working directory: %v", err)
		}
	}()

	err = os.Chdir(suite.tempDir)
	require.NoError(suite.T(), err)

	// Load configuration
	cfg, err := Load()
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), cfg)

	// Verify loaded values
	assert.Equal(suite.T(), "9000", cfg.Server.Port)
	assert.Equal(suite.T(), 45*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(suite.T(), 45*time.Second, cfg.Server.WriteTimeout)
	assert.Equal(suite.T(), 180*time.Second, cfg.Server.IdleTimeout)

	assert.Equal(suite.T(), "postgres://testuser:testpass@localhost:5432/testdb", cfg.Database.URL)
	assert.Equal(suite.T(), 30, cfg.Database.MaxOpenConns)
	assert.Equal(suite.T(), 8, cfg.Database.MaxIdleConns)
	assert.Equal(suite.T(), 10*time.Minute, cfg.Database.ConnMaxLifetime)
	assert.Equal(suite.T(), 8*time.Minute, cfg.Database.ConnMaxIdleTime)

	assert.Equal(suite.T(), "debug", cfg.Logging.Level)
	assert.Equal(suite.T(), "text", cfg.Logging.Format)
}

// TestLoadConfigWithEnvironmentOverrides tests environment variable overrides
func (suite *ConfigIntegrationTestSuite) TestLoadConfigWithEnvironmentOverrides() {
	// Create base config file
	configContent := `
server:
  name: "test-server"
  port: "8080"
  read_timeout: "30s"

database:
  url: "postgres://localhost:5432/base_db"
  max_open_conns: 25

logging:
  level: "info"
  format: "json"

auth:
  jwt_secret: "test-secret-key-32-chars-minimum-length"
  access_token_ttl: "15m"
  refresh_token_ttl: "24h"
`

	err := os.WriteFile(suite.configPath, []byte(configContent), 0644)
	require.NoError(suite.T(), err)

	// Set environment variable overrides
	envOverrides := map[string]string{
		"GRAPHQL_SERVICE_SERVER_PORT":             "9090",
		"GRAPHQL_SERVICE_DATABASE_URL":            "postgres://override:pass@localhost:5432/override_db",
		"GRAPHQL_SERVICE_DATABASE_MAX_OPEN_CONNS": "50",
		"GRAPHQL_SERVICE_LOGGING_LEVEL":           "debug",
	}

	for key, value := range envOverrides {
		suite.T().Setenv(key, value)
	}

	// Change working directory to temp dir
	originalWd, err := os.Getwd()
	require.NoError(suite.T(), err)
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			suite.T().Logf("Failed to restore working directory: %v", err)
		}
	}()

	err = os.Chdir(suite.tempDir)
	require.NoError(suite.T(), err)

	// Load configuration
	cfg, err := Load()
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), cfg)

	// Verify environment overrides took effect
	assert.Equal(suite.T(), "9090", cfg.Server.Port)                                                 // overridden
	assert.Equal(suite.T(), 30*time.Second, cfg.Server.ReadTimeout)                                  // from file
	assert.Equal(suite.T(), "postgres://override:pass@localhost:5432/override_db", cfg.Database.URL) // overridden
	assert.Equal(suite.T(), 50, cfg.Database.MaxOpenConns)                                           // overridden
	assert.Equal(suite.T(), "debug", cfg.Logging.Level)                                              // overridden
	assert.Equal(suite.T(), "json", cfg.Logging.Format)                                              // from file
}

// TestLoadConfigWithPartialFile tests loading with partial config file
func (suite *ConfigIntegrationTestSuite) TestLoadConfigWithPartialFile() {
	// Create partial config file (only server section)
	configContent := `
server:
  name: "test-server"
  port: "7000"
  read_timeout: "60s"

auth:
  jwt_secret: "test-secret-key-32-chars-minimum-length"
  access_token_ttl: "15m"
  refresh_token_ttl: "24h"
`

	err := os.WriteFile(suite.configPath, []byte(configContent), 0644)
	require.NoError(suite.T(), err)

	// Change working directory to temp dir
	originalWd, err := os.Getwd()
	require.NoError(suite.T(), err)
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			suite.T().Logf("Failed to restore working directory: %v", err)
		}
	}()

	err = os.Chdir(suite.tempDir)
	require.NoError(suite.T(), err)

	// Load configuration
	cfg, err := Load()
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), cfg)

	// Verify partial file values and defaults
	assert.Equal(suite.T(), "7000", cfg.Server.Port)                                                                     // from file
	assert.Equal(suite.T(), 60*time.Second, cfg.Server.ReadTimeout)                                                      // from file
	assert.Equal(suite.T(), 30*time.Second, cfg.Server.WriteTimeout)                                                     // default
	assert.Equal(suite.T(), "postgres://user:password@localhost:5432/graphql_service?sslmode=disable", cfg.Database.URL) // default
	assert.Equal(suite.T(), "info", cfg.Logging.Level)                                                                   // default
}

// TestLoadConfigWithInvalidFile tests handling of invalid config file
func (suite *ConfigIntegrationTestSuite) TestLoadConfigWithInvalidFile() {
	// Create invalid YAML file
	configContent := `
server:
  name: "test-server"
  port: "8080"
  invalid_yaml: [unclosed array
`

	err := os.WriteFile(suite.configPath, []byte(configContent), 0644)
	require.NoError(suite.T(), err)

	// Change working directory to temp dir
	originalWd, err := os.Getwd()
	require.NoError(suite.T(), err)
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			suite.T().Logf("Failed to restore working directory: %v", err)
		}
	}()

	err = os.Chdir(suite.tempDir)
	require.NoError(suite.T(), err)

	// Load configuration should fail
	cfg, err := Load()
	require.Error(suite.T(), err)
	assert.Nil(suite.T(), cfg)
	assert.Contains(suite.T(), err.Error(), "failed to read config file")
}

// TestLoadConfigWithValidationErrors tests configuration validation
func (suite *ConfigIntegrationTestSuite) TestLoadConfigWithValidationErrors() {
	// Create config file with invalid values
	configContent := `
server:
  name: "test-server"
  port: ""  # invalid: empty port
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "120s"

database:
  url: "postgres://localhost:5432/test"
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: "5m"
  conn_max_idle_time: "5m"

logging:
  level: "invalid_level"  # invalid: unsupported level
  format: "json"

auth:
  jwt_secret: "test-secret-key-32-chars-minimum-length"
  access_token_ttl: "15m"
  refresh_token_ttl: "24h"
`

	err := os.WriteFile(suite.configPath, []byte(configContent), 0644)
	require.NoError(suite.T(), err)

	// Change working directory to temp dir
	originalWd, err := os.Getwd()
	require.NoError(suite.T(), err)
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			suite.T().Logf("Failed to restore working directory: %v", err)
		}
	}()

	err = os.Chdir(suite.tempDir)
	require.NoError(suite.T(), err)

	// Load configuration should fail validation
	cfg, err := Load()
	require.Error(suite.T(), err)
	assert.Nil(suite.T(), cfg)
	assert.Contains(suite.T(), err.Error(), "config validation failed")
}

// TestLoadConfigNoFile tests loading with no config file (defaults only)
func (suite *ConfigIntegrationTestSuite) TestLoadConfigNoFile() {
	// Create a separate temp directory to ensure no config file exists
	emptyTempDir, err := os.MkdirTemp("", "config_no_file_test")
	require.NoError(suite.T(), err)
	defer os.RemoveAll(emptyTempDir)

	// Change to empty temp dir where no config file exists
	originalWd, err := os.Getwd()
	require.NoError(suite.T(), err)
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			suite.T().Logf("Failed to restore working directory: %v", err)
		}
	}()

	err = os.Chdir(emptyTempDir)
	require.NoError(suite.T(), err)

	// Set required environment variables for validation
	suite.T().Setenv("GRAPHQL_SERVICE_SERVER_NAME", "test-server")
	suite.T().Setenv("GRAPHQL_SERVICE_AUTH_JWT_SECRET", "test-secret-key-32-chars-minimum-length")
	suite.T().Setenv("GRAPHQL_SERVICE_AUTH_ACCESS_TOKEN_TTL", "15m")
	suite.T().Setenv("GRAPHQL_SERVICE_AUTH_REFRESH_TOKEN_TTL", "24h")

	// Reset viper to ensure clean state
	viper.Reset()

	// Load configuration (should use defaults)
	cfg, err := Load()
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), cfg)

	// Verify all defaults are used (server name is overridden by env var)
	assert.Equal(suite.T(), "test-server", cfg.Server.Name)
	assert.Equal(suite.T(), "8080", cfg.Server.Port)
	assert.Equal(suite.T(), 30*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(suite.T(), 30*time.Second, cfg.Server.WriteTimeout)
	assert.Equal(suite.T(), 120*time.Second, cfg.Server.IdleTimeout)

	assert.Equal(suite.T(), "postgres://user:password@localhost:5432/graphql_service?sslmode=disable", cfg.Database.URL)
	assert.Equal(suite.T(), 25, cfg.Database.MaxOpenConns)
	assert.Equal(suite.T(), 5, cfg.Database.MaxIdleConns)
	assert.Equal(suite.T(), 5*time.Minute, cfg.Database.ConnMaxLifetime)
	assert.Equal(suite.T(), 5*time.Minute, cfg.Database.ConnMaxIdleTime)

	assert.Equal(suite.T(), "info", cfg.Logging.Level)
	assert.Equal(suite.T(), "json", cfg.Logging.Format)
}

// TestLoadConfigComplexEnvironmentVariables tests complex environment variable scenarios
func (suite *ConfigIntegrationTestSuite) TestLoadConfigComplexEnvironmentVariables() {
	// Set all possible environment variables
	envVars := map[string]string{
		"GRAPHQL_SERVICE_SERVER_NAME":                 "test-server",
		"GRAPHQL_SERVICE_SERVER_PORT":                 "3000",
		"GRAPHQL_SERVICE_SERVER_READ_TIMEOUT":         "120s",
		"GRAPHQL_SERVICE_SERVER_WRITE_TIMEOUT":        "90s",
		"GRAPHQL_SERVICE_SERVER_IDLE_TIMEOUT":         "300s",
		"GRAPHQL_SERVICE_DATABASE_URL":                "postgres://envuser:envpass@envhost:5432/envdb?sslmode=require",
		"GRAPHQL_SERVICE_DATABASE_MAX_OPEN_CONNS":     "100",
		"GRAPHQL_SERVICE_DATABASE_MAX_IDLE_CONNS":     "20",
		"GRAPHQL_SERVICE_DATABASE_CONN_MAX_LIFETIME":  "30m",
		"GRAPHQL_SERVICE_DATABASE_CONN_MAX_IDLE_TIME": "15m",
		"GRAPHQL_SERVICE_LOGGING_LEVEL":               "warn",
		"GRAPHQL_SERVICE_LOGGING_FORMAT":              "text",
		"GRAPHQL_SERVICE_AUTH_JWT_SECRET":             "test-secret-key-32-chars-minimum-length",
		"GRAPHQL_SERVICE_AUTH_ACCESS_TOKEN_TTL":       "15m",
		"GRAPHQL_SERVICE_AUTH_REFRESH_TOKEN_TTL":      "24h",
	}

	for key, value := range envVars {
		suite.T().Setenv(key, value)
	}

	// Load configuration
	cfg, err := Load()
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), cfg)

	// Verify all environment variables were applied
	assert.Equal(suite.T(), "test-server", cfg.Server.Name)
	assert.Equal(suite.T(), "3000", cfg.Server.Port)
	assert.Equal(suite.T(), 120*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(suite.T(), 90*time.Second, cfg.Server.WriteTimeout)
	assert.Equal(suite.T(), 300*time.Second, cfg.Server.IdleTimeout)

	assert.Equal(suite.T(), "postgres://envuser:envpass@envhost:5432/envdb?sslmode=require", cfg.Database.URL)
	assert.Equal(suite.T(), 100, cfg.Database.MaxOpenConns)
	assert.Equal(suite.T(), 20, cfg.Database.MaxIdleConns)
	assert.Equal(suite.T(), 30*time.Minute, cfg.Database.ConnMaxLifetime)
	assert.Equal(suite.T(), 15*time.Minute, cfg.Database.ConnMaxIdleTime)

	assert.Equal(suite.T(), "warn", cfg.Logging.Level)
	assert.Equal(suite.T(), "text", cfg.Logging.Format)
}

// clearEnvVars clears all GRAPHQL_SERVICE environment variables
func (suite *ConfigIntegrationTestSuite) clearEnvVars() {
	envVars := []string{
		"GRAPHQL_SERVICE_SERVER_NAME",
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
		"GRAPHQL_SERVICE_AUTH_JWT_SECRET",
		"GRAPHQL_SERVICE_AUTH_ACCESS_TOKEN_TTL",
		"GRAPHQL_SERVICE_AUTH_REFRESH_TOKEN_TTL",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}

// TestConfigIntegration runs the integration test suite
func TestConfigIntegration(t *testing.T) {
	suite.Run(t, new(ConfigIntegrationTestSuite))
}
