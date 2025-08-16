package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/infrastructure/config"
	"github.com/stretchr/testify/require"
)

// TestDBConfig holds configuration for test database setup
type TestDBConfig struct {
	BaseURL         string
	TestDBName      string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultTestDBConfig returns default test database configuration
func DefaultTestDBConfig() TestDBConfig {
	baseURL := os.Getenv("TEST_DATABASE_URL")
	if baseURL == "" {
		baseURL = "postgres://user:password@localhost:5432/postgres?sslmode=disable"
	}

	return TestDBConfig{
		BaseURL:         baseURL,
		TestDBName:      fmt.Sprintf("test_db_%d", time.Now().UnixNano()),
		MaxOpenConns:    5,
		MaxIdleConns:    2,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}
}

// TestDBManager manages test database lifecycle
type TestDBManager struct {
	config    TestDBConfig
	adminDB   *sql.DB
	testDB    *DB
	logger    *slog.Logger
	dbCreated bool
}

// NewTestDBManager creates a new test database manager
func NewTestDBManager(t *testing.T, cfg TestDBConfig) *TestDBManager {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	// Connect to admin database (usually postgres database)
	adminDB, err := sql.Open("postgres", cfg.BaseURL)
	require.NoError(t, err)

	// Test admin connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = adminDB.PingContext(ctx)
	require.NoError(t, err)

	return &TestDBManager{
		config:  cfg,
		adminDB: adminDB,
		logger:  logger,
	}
}

// CreateTestDB creates a new test database
func (tm *TestDBManager) CreateTestDB(t *testing.T) *DB {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create test database
	createSQL := fmt.Sprintf("CREATE DATABASE %s", tm.config.TestDBName)
	_, err := tm.adminDB.ExecContext(ctx, createSQL)
	require.NoError(t, err)
	tm.dbCreated = true

	// Build test database URL
	testDBURL := tm.buildTestDBURL()

	// Connect to test database
	dbConfig := config.DatabaseConfig{
		URL:             testDBURL,
		MaxOpenConns:    tm.config.MaxOpenConns,
		MaxIdleConns:    tm.config.MaxIdleConns,
		ConnMaxLifetime: tm.config.ConnMaxLifetime,
		ConnMaxIdleTime: tm.config.ConnMaxIdleTime,
	}

	testDB, err := NewConnection(dbConfig, tm.logger)
	require.NoError(t, err)

	tm.testDB = testDB
	return testDB
}

// GetTestDB returns the test database connection (creates if not exists)
func (tm *TestDBManager) GetTestDB(t *testing.T) *DB {
	if tm.testDB == nil {
		return tm.CreateTestDB(t)
	}
	return tm.testDB
}

// RunMigrations runs migrations on the test database
func (tm *TestDBManager) RunMigrations(t *testing.T, migrationsPath string) {
	testDB := tm.GetTestDB(t)

	migrationManager := NewMigrationManager(testDB, tm.logger)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err := migrationManager.RunMigrations(ctx, migrationsPath)
	require.NoError(t, err)
}

// CleanupTestDB drops the test database and closes connections
func (tm *TestDBManager) CleanupTestDB(t *testing.T) {
	// Close test database connection first
	if tm.testDB != nil {
		tm.testDB.Close()
		tm.testDB = nil
	}

	// Drop test database if it was created
	if tm.dbCreated {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Terminate any remaining connections to the test database
		terminateSQL := fmt.Sprintf(`
			SELECT pg_terminate_backend(pid)
			FROM pg_stat_activity
			WHERE datname = '%s' AND pid <> pg_backend_pid()
		`, tm.config.TestDBName)

		_, _ = tm.adminDB.ExecContext(ctx, terminateSQL)

		// Drop the test database
		dropSQL := fmt.Sprintf("DROP DATABASE IF EXISTS %s", tm.config.TestDBName)
		_, err := tm.adminDB.ExecContext(ctx, dropSQL)
		if err != nil {
			t.Logf("Warning: Failed to drop test database %s: %v", tm.config.TestDBName, err)
		}
		tm.dbCreated = false
	}

	// Close admin database connection
	if tm.adminDB != nil {
		tm.adminDB.Close()
		tm.adminDB = nil
	}
}

// buildTestDBURL builds the test database URL from base URL and test DB name
func (tm *TestDBManager) buildTestDBURL() string {
	baseURL := tm.config.BaseURL

	// Parse the base URL to replace the database name
	// This is a simple approach - for production use, consider using url.Parse
	if strings.Contains(baseURL, "/postgres?") {
		return strings.Replace(baseURL, "/postgres?", "/"+tm.config.TestDBName+"?", 1)
	} else if strings.Contains(baseURL, "/postgres") && strings.HasSuffix(baseURL, "/postgres") {
		return strings.Replace(baseURL, "/postgres", "/"+tm.config.TestDBName, 1)
	} else {
		// Fallback: append database name
		separator := "/"
		if strings.HasSuffix(baseURL, "/") {
			separator = ""
		}
		return baseURL + separator + tm.config.TestDBName
	}
}

// TruncateAllTables truncates all tables in the test database (useful for test cleanup)
func (tm *TestDBManager) TruncateAllTables(t *testing.T) {
	testDB := tm.GetTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all table names (excluding system tables)
	query := `
		SELECT tablename 
		FROM pg_tables 
		WHERE schemaname = 'public' 
		AND tablename != 'schema_migrations'
	`

	rows, err := testDB.QueryContext(ctx, query)
	require.NoError(t, err)
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		require.NoError(t, err)
		tables = append(tables, tableName)
	}

	// Truncate all tables
	if len(tables) > 0 {
		truncateSQL := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE",
			strings.Join(tables, ", "))
		_, err = testDB.ExecContext(ctx, truncateSQL)
		require.NoError(t, err)
	}
}

// ExecuteSQL executes arbitrary SQL on the test database (useful for test setup)
func (tm *TestDBManager) ExecuteSQL(t *testing.T, sql string, args ...interface{}) {
	testDB := tm.GetTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := testDB.ExecContext(ctx, sql, args...)
	require.NoError(t, err)
}

// QueryRow executes a query that returns a single row
func (tm *TestDBManager) QueryRow(t *testing.T, sql string, args ...interface{}) *sql.Row {
	testDB := tm.GetTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return testDB.QueryRowContext(ctx, sql, args...)
}

// TestDBSetup is a helper function for setting up test database in test suites
func TestDBSetup(t *testing.T, migrationsPath string) (*DB, func()) {
	// Skip if no test database URL is provided
	if os.Getenv("TEST_DATABASE_URL") == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping database integration tests")
	}

	cfg := DefaultTestDBConfig()
	manager := NewTestDBManager(t, cfg)

	// Create test database and run migrations
	db := manager.CreateTestDB(t)
	if migrationsPath != "" {
		manager.RunMigrations(t, migrationsPath)
	}

	// Return database and cleanup function
	cleanup := func() {
		manager.CleanupTestDB(t)
	}

	return db, cleanup
}

// TestDBSetupWithConfig is a helper function for setting up test database with custom config
func TestDBSetupWithConfig(t *testing.T, cfg TestDBConfig, migrationsPath string) (*DB, func()) {
	// Skip if no test database URL is provided
	if os.Getenv("TEST_DATABASE_URL") == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping database integration tests")
	}

	manager := NewTestDBManager(t, cfg)

	// Create test database and run migrations
	db := manager.CreateTestDB(t)
	if migrationsPath != "" {
		manager.RunMigrations(t, migrationsPath)
	}

	// Return database and cleanup function
	cleanup := func() {
		manager.CleanupTestDB(t)
	}

	return db, cleanup
}
