package logging

import (
	"sync"

	"github.com/captain-corgi/go-graphql-example/internal/infrastructure/config"
)

// LoggerFactory manages logger creation and configuration
type LoggerFactory struct {
	config config.LoggingConfig
	logger *Logger
	mu     sync.RWMutex
}

// NewLoggerFactory creates a new logger factory
func NewLoggerFactory(cfg config.LoggingConfig) *LoggerFactory {
	return &LoggerFactory{
		config: cfg,
		logger: NewLogger(cfg),
	}
}

// GetLogger returns the base logger
func (f *LoggerFactory) GetLogger() *Logger {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.logger
}

// GetDomainLogger returns a logger configured for domain layer
func (f *LoggerFactory) GetDomainLogger() *DomainLogger {
	return NewDomainLogger(f.GetLogger())
}

// GetApplicationLogger returns a logger configured for application layer
func (f *LoggerFactory) GetApplicationLogger() *ApplicationLogger {
	return NewApplicationLogger(f.GetLogger())
}

// GetInfrastructureLogger returns a logger configured for infrastructure layer
func (f *LoggerFactory) GetInfrastructureLogger() *InfrastructureLogger {
	return NewInfrastructureLogger(f.GetLogger())
}

// GetInterfaceLogger returns a logger configured for interface layer
func (f *LoggerFactory) GetInterfaceLogger() *InterfaceLogger {
	return NewInterfaceLogger(f.GetLogger())
}

// UpdateConfig updates the logging configuration and recreates the logger
func (f *LoggerFactory) UpdateConfig(cfg config.LoggingConfig) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.config = cfg
	f.logger = NewLogger(cfg)
}

// GetConfig returns the current logging configuration
func (f *LoggerFactory) GetConfig() config.LoggingConfig {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.config
}
