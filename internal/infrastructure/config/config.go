package config

import (
	"fmt"
	"time"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Auth     AuthConfig     `mapstructure:"auth"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Name         string        `mapstructure:"name"`
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	URL             string        `mapstructure:"url"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret       string        `mapstructure:"jwt_secret"`
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
}

// Validate validates the configuration and returns an error if invalid
func (c *Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server config validation failed: %w", err)
	}

	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database config validation failed: %w", err)
	}

	if err := c.Logging.Validate(); err != nil {
		return fmt.Errorf("logging config validation failed: %w", err)
	}

	if err := c.Auth.Validate(); err != nil {
		return fmt.Errorf("auth config validation failed: %w", err)
	}

	return nil
}

// Validate validates server configuration
func (s *ServerConfig) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("server name is required")
	}

	if s.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if s.ReadTimeout <= 0 {
		return fmt.Errorf("server read timeout must be positive")
	}

	if s.WriteTimeout <= 0 {
		return fmt.Errorf("server write timeout must be positive")
	}

	if s.IdleTimeout <= 0 {
		return fmt.Errorf("server idle timeout must be positive")
	}

	return nil
}

// Validate validates database configuration
func (d *DatabaseConfig) Validate() error {
	if d.URL == "" {
		return fmt.Errorf("database URL is required")
	}

	if d.MaxOpenConns <= 0 {
		return fmt.Errorf("database max open connections must be positive")
	}

	if d.MaxIdleConns <= 0 {
		return fmt.Errorf("database max idle connections must be positive")
	}

	if d.MaxIdleConns > d.MaxOpenConns {
		return fmt.Errorf("database max idle connections cannot exceed max open connections")
	}

	if d.ConnMaxLifetime <= 0 {
		return fmt.Errorf("database connection max lifetime must be positive")
	}

	if d.ConnMaxIdleTime <= 0 {
		return fmt.Errorf("database connection max idle time must be positive")
	}

	return nil
}

// Validate validates logging configuration
func (l *LoggingConfig) Validate() error {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLevels[l.Level] {
		return fmt.Errorf("invalid log level: %s (must be one of: debug, info, warn, error)", l.Level)
	}

	validFormats := map[string]bool{
		"json": true,
		"text": true,
	}

	if !validFormats[l.Format] {
		return fmt.Errorf("invalid log format: %s (must be one of: json, text)", l.Format)
	}

	return nil
}

// Validate validates authentication configuration
func (a *AuthConfig) Validate() error {
	if a.JWTSecret == "" {
		return fmt.Errorf("JWT secret cannot be empty")
	}

	if len(a.JWTSecret) < 32 {
		return fmt.Errorf("JWT secret must be at least 32 characters long")
	}

	if a.AccessTokenTTL <= 0 {
		return fmt.Errorf("access token TTL must be positive")
	}

	if a.RefreshTokenTTL <= 0 {
		return fmt.Errorf("refresh token TTL must be positive")
	}

	if a.RefreshTokenTTL <= a.AccessTokenTTL {
		return fmt.Errorf("refresh token TTL must be greater than access token TTL")
	}

	return nil
}
