package logging

import (
	"testing"

	"github.com/captain-corgi/go-graphql-example/internal/infrastructure/config"
)

func TestNewLoggerFactory(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "info",
		Format: "json",
	}

	factory := NewLoggerFactory(cfg)

	if factory == nil {
		t.Error("NewLoggerFactory returned nil")
		return
	}
	if factory.logger == nil {
		t.Error("Factory logger is nil")
	}
	if factory.config != cfg {
		t.Error("Factory config not set correctly")
	}
}

func TestLoggerFactory_GetLogger(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "debug",
		Format: "text",
	}

	factory := NewLoggerFactory(cfg)
	logger := factory.GetLogger()

	if logger == nil {
		t.Error("GetLogger returned nil")
		return
	}
	if logger.Logger == nil {
		t.Error("Logger.Logger is nil")
	}
}

func TestLoggerFactory_GetDomainLogger(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "info",
		Format: "json",
	}

	factory := NewLoggerFactory(cfg)
	domainLogger := factory.GetDomainLogger()

	if domainLogger == nil {
		t.Error("GetDomainLogger returned nil")
		return
	}
	if domainLogger.Logger == nil {
		t.Error("DomainLogger.Logger is nil")
	}
}

func TestLoggerFactory_GetApplicationLogger(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "info",
		Format: "json",
	}

	factory := NewLoggerFactory(cfg)
	appLogger := factory.GetApplicationLogger()

	if appLogger == nil {
		t.Error("GetApplicationLogger returned nil")
		return
	}
	if appLogger.Logger == nil {
		t.Error("ApplicationLogger.Logger is nil")
	}
}

func TestLoggerFactory_GetInfrastructureLogger(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "info",
		Format: "json",
	}

	factory := NewLoggerFactory(cfg)
	infraLogger := factory.GetInfrastructureLogger()

	if infraLogger == nil {
		t.Error("GetInfrastructureLogger returned nil")
		return
	}
	if infraLogger.Logger == nil {
		t.Error("InfrastructureLogger.Logger is nil")
	}
}

func TestLoggerFactory_GetInterfaceLogger(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "info",
		Format: "json",
	}

	factory := NewLoggerFactory(cfg)
	interfaceLogger := factory.GetInterfaceLogger()

	if interfaceLogger == nil {
		t.Error("GetInterfaceLogger returned nil")
		return
	}
	if interfaceLogger.Logger == nil {
		t.Error("InterfaceLogger.Logger is nil")
	}
}

func TestLoggerFactory_UpdateConfig(t *testing.T) {
	initialCfg := config.LoggingConfig{
		Level:  "info",
		Format: "json",
	}

	factory := NewLoggerFactory(initialCfg)
	initialLogger := factory.GetLogger()

	// Update configuration
	newCfg := config.LoggingConfig{
		Level:  "debug",
		Format: "text",
	}
	factory.UpdateConfig(newCfg)

	// Verify config was updated
	if factory.GetConfig() != newCfg {
		t.Error("Config was not updated correctly")
	}

	// Verify logger was recreated
	newLogger := factory.GetLogger()
	if newLogger == initialLogger {
		t.Error("Logger should have been recreated after config update")
	}
}

func TestLoggerFactory_GetConfig(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "warn",
		Format: "json",
	}

	factory := NewLoggerFactory(cfg)
	retrievedCfg := factory.GetConfig()

	if retrievedCfg != cfg {
		t.Error("GetConfig returned incorrect configuration")
	}
}

func TestLoggerFactory_ConcurrentAccess(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "info",
		Format: "json",
	}

	factory := NewLoggerFactory(cfg)

	// Test concurrent access to various methods
	done := make(chan bool, 4)

	// Concurrent GetLogger calls
	go func() {
		for i := 0; i < 100; i++ {
			logger := factory.GetLogger()
			if logger == nil {
				t.Error("GetLogger returned nil during concurrent access")
			}
		}
		done <- true
	}()

	// Concurrent GetDomainLogger calls
	go func() {
		for i := 0; i < 100; i++ {
			logger := factory.GetDomainLogger()
			if logger == nil {
				t.Error("GetDomainLogger returned nil during concurrent access")
			}
		}
		done <- true
	}()

	// Concurrent GetConfig calls
	go func() {
		for i := 0; i < 100; i++ {
			config := factory.GetConfig()
			if config.Level == "" {
				t.Error("GetConfig returned invalid config during concurrent access")
			}
		}
		done <- true
	}()

	// Concurrent UpdateConfig calls
	go func() {
		for i := 0; i < 10; i++ {
			newCfg := config.LoggingConfig{
				Level:  "debug",
				Format: "text",
			}
			factory.UpdateConfig(newCfg)
		}
		done <- true
	}()

	// Wait for all goroutines to complete
	for i := 0; i < 4; i++ {
		<-done
	}
}
