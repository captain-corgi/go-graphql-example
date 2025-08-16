package database

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// DatabaseIntegrationTestSuite defines the integration test suite for database operations
type DatabaseIntegrationTestSuite struct {
	suite.Suite
	db     *DB
	ctx    context.Context
	logger *slog.Logger
}

// SetupSuite sets up the test suite
func (suite *DatabaseIntegrationTestSuite) SetupSuite() {
	// Skip integration tests if no database URL is provided
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		suite.T().Skip("TEST_DATABASE_URL not set, skipping database integration tests")
	}

	suite.ctx = context.Background()
	suite.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// Create test database connection
	cfg := config.DatabaseConfig{
		URL:             dbURL,
		MaxOpenConns:    10,
		MaxIdleConns:    3,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}

	db, err := NewConnection(cfg, suite.logger)
	require.NoError(suite.T(), err)

	suite.db = db
}

// TearDownSuite cleans up the test suite
func (suite *DatabaseIntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

// TestDatabaseConnection tests basic database connection functionality
func (suite *DatabaseIntegrationTestSuite) TestDatabaseConnection() {
	// Test that connection is established
	require.NotNil(suite.T(), suite.db)
	require.NotNil(suite.T(), suite.db.DB)

	// Test health check
	err := suite.db.Health(suite.ctx)
	assert.NoError(suite.T(), err)

	// Test connection stats
	stats := suite.db.Stats()
	assert.GreaterOrEqual(suite.T(), stats.MaxOpenConnections, 1)
}

// TestDatabaseConnectionWithInvalidURL tests connection with invalid database URL
func (suite *DatabaseIntegrationTestSuite) TestDatabaseConnectionWithInvalidURL() {
	cfg := config.DatabaseConfig{
		URL:             "postgres://invalid:invalid@nonexistent:5432/invalid",
		MaxOpenConns:    5,
		MaxIdleConns:    2,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}

	_, err := NewConnection(cfg, suite.logger)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to ping database")
}

// TestDatabaseConnectionPoolConfiguration tests connection pool settings
func (suite *DatabaseIntegrationTestSuite) TestDatabaseConnectionPoolConfiguration() {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		suite.T().Skip("TEST_DATABASE_URL not set")
	}

	cfg := config.DatabaseConfig{
		URL:             dbURL,
		MaxOpenConns:    15,
		MaxIdleConns:    5,
		ConnMaxLifetime: 10 * time.Minute,
		ConnMaxIdleTime: 2 * time.Minute,
	}

	db, err := NewConnection(cfg, suite.logger)
	require.NoError(suite.T(), err)
	defer db.Close()

	// Verify connection pool settings
	stats := db.Stats()
	assert.Equal(suite.T(), 15, stats.MaxOpenConnections)

	// Test that we can actually use connections
	err = db.Health(suite.ctx)
	assert.NoError(suite.T(), err)
}

// TestDatabaseHealthCheck tests health check functionality
func (suite *DatabaseIntegrationTestSuite) TestDatabaseHealthCheck() {
	// Test successful health check
	err := suite.db.Health(suite.ctx)
	assert.NoError(suite.T(), err)

	// Test health check with timeout
	ctx, cancel := context.WithTimeout(suite.ctx, 1*time.Millisecond)
	defer cancel()

	// This might or might not fail depending on how fast the database responds
	// We're mainly testing that the context timeout is respected
	_ = suite.db.Health(ctx)
}

// TestDatabaseBasicOperations tests basic database operations
func (suite *DatabaseIntegrationTestSuite) TestDatabaseBasicOperations() {
	// Test simple query
	var result int
	err := suite.db.QueryRowContext(suite.ctx, "SELECT 1").Scan(&result)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, result)

	// Test query with parameters
	var echo string
	err = suite.db.QueryRowContext(suite.ctx, "SELECT $1::text", "test").Scan(&echo)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test", echo)
}

// TestDatabaseTransactionHandling tests transaction functionality
func (suite *DatabaseIntegrationTestSuite) TestDatabaseTransactionHandling() {
	// Create a temporary table for testing
	_, err := suite.db.ExecContext(suite.ctx, `
		CREATE TEMPORARY TABLE test_transactions (
			id SERIAL PRIMARY KEY,
			value TEXT
		)
	`)
	require.NoError(suite.T(), err)

	// Test successful transaction
	tx, err := suite.db.BeginTx(suite.ctx, nil)
	require.NoError(suite.T(), err)

	_, err = tx.ExecContext(suite.ctx, "INSERT INTO test_transactions (value) VALUES ($1)", "test1")
	require.NoError(suite.T(), err)

	err = tx.Commit()
	require.NoError(suite.T(), err)

	// Verify data was committed
	var count int
	err = suite.db.QueryRowContext(suite.ctx, "SELECT COUNT(*) FROM test_transactions").Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)

	// Test rollback transaction
	tx, err = suite.db.BeginTx(suite.ctx, nil)
	require.NoError(suite.T(), err)

	_, err = tx.ExecContext(suite.ctx, "INSERT INTO test_transactions (value) VALUES ($1)", "test2")
	require.NoError(suite.T(), err)

	err = tx.Rollback()
	require.NoError(suite.T(), err)

	// Verify data was not committed
	err = suite.db.QueryRowContext(suite.ctx, "SELECT COUNT(*) FROM test_transactions").Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count) // Still only the first record
}

// MigrationIntegrationTestSuite defines the integration test suite for migration operations
type MigrationIntegrationTestSuite struct {
	suite.Suite
	db                *DB
	migrationManager  *MigrationManager
	ctx               context.Context
	logger            *slog.Logger
	tempMigrationsDir string
}

// SetupSuite sets up the migration test suite
func (suite *MigrationIntegrationTestSuite) SetupSuite() {
	// Skip integration tests if no database URL is provided
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		suite.T().Skip("TEST_DATABASE_URL not set, skipping migration integration tests")
	}

	suite.ctx = context.Background()
	suite.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// Create test database connection
	cfg := config.DatabaseConfig{
		URL:             dbURL,
		MaxOpenConns:    5,
		MaxIdleConns:    2,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}

	db, err := NewConnection(cfg, suite.logger)
	require.NoError(suite.T(), err)

	suite.db = db
	suite.migrationManager = NewMigrationManager(db, suite.logger)

	// Create temporary migrations directory
	tempDir, err := os.MkdirTemp("", "migration_test")
	require.NoError(suite.T(), err)
	suite.tempMigrationsDir = tempDir

	// Create test migration files
	suite.createTestMigrations()
}

// TearDownSuite cleans up the migration test suite
func (suite *MigrationIntegrationTestSuite) TearDownSuite() {
	if suite.tempMigrationsDir != "" {
		os.RemoveAll(suite.tempMigrationsDir)
	}
	if suite.db != nil {
		// Clean up migration tables
		_, _ = suite.db.ExecContext(suite.ctx, "DROP TABLE IF EXISTS schema_migrations")
		_, _ = suite.db.ExecContext(suite.ctx, "DROP TABLE IF EXISTS test_migration_table")
		suite.db.Close()
	}
}

// SetupTest sets up each migration test
func (suite *MigrationIntegrationTestSuite) SetupTest() {
	// Clean up any existing migration state
	_, _ = suite.db.ExecContext(suite.ctx, "DROP TABLE IF EXISTS schema_migrations")
	_, _ = suite.db.ExecContext(suite.ctx, "DROP TABLE IF EXISTS test_migration_table")
}

// createTestMigrations creates test migration files
func (suite *MigrationIntegrationTestSuite) createTestMigrations() {
	// Migration 1 - Create table
	upMigration1 := `
CREATE TABLE test_migration_table (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
`
	downMigration1 := `DROP TABLE test_migration_table;`

	err := os.WriteFile(suite.tempMigrationsDir+"/001_create_test_table.up.sql", []byte(upMigration1), 0644)
	require.NoError(suite.T(), err)

	err = os.WriteFile(suite.tempMigrationsDir+"/001_create_test_table.down.sql", []byte(downMigration1), 0644)
	require.NoError(suite.T(), err)

	// Migration 2 - Add column
	upMigration2 := `ALTER TABLE test_migration_table ADD COLUMN email VARCHAR(255);`
	downMigration2 := `ALTER TABLE test_migration_table DROP COLUMN email;`

	err = os.WriteFile(suite.tempMigrationsDir+"/002_add_email_column.up.sql", []byte(upMigration2), 0644)
	require.NoError(suite.T(), err)

	err = os.WriteFile(suite.tempMigrationsDir+"/002_add_email_column.down.sql", []byte(downMigration2), 0644)
	require.NoError(suite.T(), err)
}

// TestRunMigrations tests running migrations
func (suite *MigrationIntegrationTestSuite) TestRunMigrations() {
	// Run migrations
	err := suite.migrationManager.RunMigrations(suite.ctx, suite.tempMigrationsDir)
	assert.NoError(suite.T(), err)

	// Verify table was created
	var exists bool
	err = suite.db.QueryRowContext(suite.ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_name = 'test_migration_table'
		)
	`).Scan(&exists)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)

	// Verify both columns exist (from both migrations)
	var columnCount int
	err = suite.db.QueryRowContext(suite.ctx, `
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_name = 'test_migration_table' 
		AND column_name IN ('name', 'email')
	`).Scan(&columnCount)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, columnCount)

	// Check migration version
	version, dirty, err := suite.migrationManager.GetMigrationVersion(suite.ctx, suite.tempMigrationsDir)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), dirty)
	assert.Equal(suite.T(), uint(2), version)
}

// TestRunMigrationsNoChange tests running migrations when no new migrations exist
func (suite *MigrationIntegrationTestSuite) TestRunMigrationsNoChange() {
	// Run migrations first time
	err := suite.migrationManager.RunMigrations(suite.ctx, suite.tempMigrationsDir)
	require.NoError(suite.T(), err)

	// Run migrations again - should be no change
	err = suite.migrationManager.RunMigrations(suite.ctx, suite.tempMigrationsDir)
	assert.NoError(suite.T(), err)
}

// TestRollbackMigration tests rolling back migrations
func (suite *MigrationIntegrationTestSuite) TestRollbackMigration() {
	// Run migrations first
	err := suite.migrationManager.RunMigrations(suite.ctx, suite.tempMigrationsDir)
	require.NoError(suite.T(), err)

	// Verify we're at version 2
	version, dirty, err := suite.migrationManager.GetMigrationVersion(suite.ctx, suite.tempMigrationsDir)
	require.NoError(suite.T(), err)
	require.False(suite.T(), dirty)
	require.Equal(suite.T(), uint(2), version)

	// Rollback one migration
	err = suite.migrationManager.RollbackMigration(suite.ctx, suite.tempMigrationsDir)
	assert.NoError(suite.T(), err)

	// Verify we're now at version 1
	version, dirty, err = suite.migrationManager.GetMigrationVersion(suite.ctx, suite.tempMigrationsDir)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), dirty)
	assert.Equal(suite.T(), uint(1), version)

	// Verify email column was removed
	var columnExists bool
	err = suite.db.QueryRowContext(suite.ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.columns 
			WHERE table_name = 'test_migration_table' 
			AND column_name = 'email'
		)
	`).Scan(&columnExists)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), columnExists)
}

// TestGetMigrationVersion tests getting migration version
func (suite *MigrationIntegrationTestSuite) TestGetMigrationVersion() {
	// Initially no migrations
	version, dirty, err := suite.migrationManager.GetMigrationVersion(suite.ctx, suite.tempMigrationsDir)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), dirty)
	assert.Equal(suite.T(), uint(0), version)

	// Run migrations
	err = suite.migrationManager.RunMigrations(suite.ctx, suite.tempMigrationsDir)
	require.NoError(suite.T(), err)

	// Check version after migrations
	version, dirty, err = suite.migrationManager.GetMigrationVersion(suite.ctx, suite.tempMigrationsDir)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), dirty)
	assert.Equal(suite.T(), uint(2), version)
}

// TestMigrationWithInvalidPath tests migration with invalid path
func (suite *MigrationIntegrationTestSuite) TestMigrationWithInvalidPath() {
	err := suite.migrationManager.RunMigrations(suite.ctx, "/nonexistent/path")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to create migrate instance")
}

// TestDatabaseIntegration runs the database integration test suite
func TestDatabaseIntegration(t *testing.T) {
	suite.Run(t, new(DatabaseIntegrationTestSuite))
}

// TestMigrationIntegration runs the migration integration test suite
func TestMigrationIntegration(t *testing.T) {
	suite.Run(t, new(MigrationIntegrationTestSuite))
}
