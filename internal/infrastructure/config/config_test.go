package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: Config{
				Server: ServerConfig{
					Port:         "8080",
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
					IdleTimeout:  120 * time.Second,
				},
				Database: DatabaseConfig{
					URL:             "postgres://localhost/test",
					MaxOpenConns:    25,
					MaxIdleConns:    5,
					ConnMaxLifetime: 5 * time.Minute,
					ConnMaxIdleTime: 5 * time.Minute,
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "json",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid server config - empty port",
			config: Config{
				Server: ServerConfig{
					Port:         "",
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
					IdleTimeout:  120 * time.Second,
				},
				Database: DatabaseConfig{
					URL:             "postgres://localhost/test",
					MaxOpenConns:    25,
					MaxIdleConns:    5,
					ConnMaxLifetime: 5 * time.Minute,
					ConnMaxIdleTime: 5 * time.Minute,
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "json",
				},
			},
			wantErr: true,
			errMsg:  "server config validation failed",
		},
		{
			name: "invalid database config - empty URL",
			config: Config{
				Server: ServerConfig{
					Port:         "8080",
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
					IdleTimeout:  120 * time.Second,
				},
				Database: DatabaseConfig{
					URL:             "",
					MaxOpenConns:    25,
					MaxIdleConns:    5,
					ConnMaxLifetime: 5 * time.Minute,
					ConnMaxIdleTime: 5 * time.Minute,
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "json",
				},
			},
			wantErr: true,
			errMsg:  "database config validation failed",
		},
		{
			name: "invalid logging config - invalid level",
			config: Config{
				Server: ServerConfig{
					Port:         "8080",
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
					IdleTimeout:  120 * time.Second,
				},
				Database: DatabaseConfig{
					URL:             "postgres://localhost/test",
					MaxOpenConns:    25,
					MaxIdleConns:    5,
					ConnMaxLifetime: 5 * time.Minute,
					ConnMaxIdleTime: 5 * time.Minute,
				},
				Logging: LoggingConfig{
					Level:  "invalid",
					Format: "json",
				},
			},
			wantErr: true,
			errMsg:  "logging config validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestServerConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ServerConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid server config",
			config: ServerConfig{
				Port:         "8080",
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 30 * time.Second,
				IdleTimeout:  120 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "empty port",
			config: ServerConfig{
				Port:         "",
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 30 * time.Second,
				IdleTimeout:  120 * time.Second,
			},
			wantErr: true,
			errMsg:  "server port is required",
		},
		{
			name: "zero read timeout",
			config: ServerConfig{
				Port:         "8080",
				ReadTimeout:  0,
				WriteTimeout: 30 * time.Second,
				IdleTimeout:  120 * time.Second,
			},
			wantErr: true,
			errMsg:  "server read timeout must be positive",
		},
		{
			name: "zero write timeout",
			config: ServerConfig{
				Port:         "8080",
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 0,
				IdleTimeout:  120 * time.Second,
			},
			wantErr: true,
			errMsg:  "server write timeout must be positive",
		},
		{
			name: "zero idle timeout",
			config: ServerConfig{
				Port:         "8080",
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 30 * time.Second,
				IdleTimeout:  0,
			},
			wantErr: true,
			errMsg:  "server idle timeout must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDatabaseConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  DatabaseConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid database config",
			config: DatabaseConfig{
				URL:             "postgres://localhost/test",
				MaxOpenConns:    25,
				MaxIdleConns:    5,
				ConnMaxLifetime: 5 * time.Minute,
				ConnMaxIdleTime: 5 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "empty URL",
			config: DatabaseConfig{
				URL:             "",
				MaxOpenConns:    25,
				MaxIdleConns:    5,
				ConnMaxLifetime: 5 * time.Minute,
				ConnMaxIdleTime: 5 * time.Minute,
			},
			wantErr: true,
			errMsg:  "database URL is required",
		},
		{
			name: "zero max open connections",
			config: DatabaseConfig{
				URL:             "postgres://localhost/test",
				MaxOpenConns:    0,
				MaxIdleConns:    5,
				ConnMaxLifetime: 5 * time.Minute,
				ConnMaxIdleTime: 5 * time.Minute,
			},
			wantErr: true,
			errMsg:  "database max open connections must be positive",
		},
		{
			name: "zero max idle connections",
			config: DatabaseConfig{
				URL:             "postgres://localhost/test",
				MaxOpenConns:    25,
				MaxIdleConns:    0,
				ConnMaxLifetime: 5 * time.Minute,
				ConnMaxIdleTime: 5 * time.Minute,
			},
			wantErr: true,
			errMsg:  "database max idle connections must be positive",
		},
		{
			name: "max idle connections exceeds max open connections",
			config: DatabaseConfig{
				URL:             "postgres://localhost/test",
				MaxOpenConns:    5,
				MaxIdleConns:    10,
				ConnMaxLifetime: 5 * time.Minute,
				ConnMaxIdleTime: 5 * time.Minute,
			},
			wantErr: true,
			errMsg:  "database max idle connections cannot exceed max open connections",
		},
		{
			name: "zero connection max lifetime",
			config: DatabaseConfig{
				URL:             "postgres://localhost/test",
				MaxOpenConns:    25,
				MaxIdleConns:    5,
				ConnMaxLifetime: 0,
				ConnMaxIdleTime: 5 * time.Minute,
			},
			wantErr: true,
			errMsg:  "database connection max lifetime must be positive",
		},
		{
			name: "zero connection max idle time",
			config: DatabaseConfig{
				URL:             "postgres://localhost/test",
				MaxOpenConns:    25,
				MaxIdleConns:    5,
				ConnMaxLifetime: 5 * time.Minute,
				ConnMaxIdleTime: 0,
			},
			wantErr: true,
			errMsg:  "database connection max idle time must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLoggingConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  LoggingConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid logging config - debug level",
			config: LoggingConfig{
				Level:  "debug",
				Format: "json",
			},
			wantErr: false,
		},
		{
			name: "valid logging config - info level",
			config: LoggingConfig{
				Level:  "info",
				Format: "json",
			},
			wantErr: false,
		},
		{
			name: "valid logging config - warn level",
			config: LoggingConfig{
				Level:  "warn",
				Format: "json",
			},
			wantErr: false,
		},
		{
			name: "valid logging config - error level",
			config: LoggingConfig{
				Level:  "error",
				Format: "json",
			},
			wantErr: false,
		},
		{
			name: "valid logging config - text format",
			config: LoggingConfig{
				Level:  "info",
				Format: "text",
			},
			wantErr: false,
		},
		{
			name: "invalid log level",
			config: LoggingConfig{
				Level:  "invalid",
				Format: "json",
			},
			wantErr: true,
			errMsg:  "invalid log level: invalid",
		},
		{
			name: "invalid log format",
			config: LoggingConfig{
				Level:  "info",
				Format: "invalid",
			},
			wantErr: true,
			errMsg:  "invalid log format: invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
