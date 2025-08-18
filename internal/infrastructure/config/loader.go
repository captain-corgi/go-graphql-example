package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Load loads configuration from files and environment variables
func Load() (*Config, error) {
	// Set up Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Set defaults
	setDefaults()

	// Environment variable overrides
	viper.SetEnvPrefix("GRAPHQL_SERVICE")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file (optional - will use defaults and env vars if not found)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found is OK, we'll use defaults and env vars
	}

	// Unmarshal into config struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.name", "graphql-service")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "120s")

	// Database defaults
	viper.SetDefault("database.url", "postgres://user:password@localhost:5432/graphql_service?sslmode=disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "5m")
	viper.SetDefault("database.conn_max_idle_time", "5m")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")

	// Auth defaults
	viper.SetDefault("auth.jwt_secret", "default-secret-key-32-chars-minimum-length")
	viper.SetDefault("auth.access_token_ttl", "15m")
	viper.SetDefault("auth.refresh_token_ttl", "24h")
}

// MustLoad loads configuration and panics if it fails
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load configuration: %v", err))
	}
	return cfg
}
