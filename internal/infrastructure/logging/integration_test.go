package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/infrastructure/config"
)

func TestLoggingIntegration(t *testing.T) {
	ctx := WithRequestID(context.Background(), "req-integration-test")

	var buf bytes.Buffer

	// Create a custom logger that writes to our buffer for testing
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	baseLogger := slog.New(handler)
	testLogger := &Logger{Logger: baseLogger}
	testDomainLogger := &DomainLogger{Logger: testLogger.WithComponent("domain")}

	// Test logging with request ID
	testDomainLogger.LogEntityCreated(ctx, "User", "user-123")

	// Parse the JSON output
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	// Verify the log entry contains expected fields
	expectedFields := map[string]string{
		"component":   "domain",
		"request_id":  "req-integration-test",
		"entity_type": "User",
		"entity_id":   "user-123",
		"msg":         "Domain entity created",
	}

	for field, expected := range expectedFields {
		if actual, ok := logEntry[field].(string); !ok || actual != expected {
			t.Errorf("Expected %s '%s', got %v", field, expected, logEntry[field])
		}
	}

	// Verify log level is info
	if level, ok := logEntry["level"].(string); !ok || level != "INFO" {
		t.Errorf("Expected level 'INFO', got %v", logEntry["level"])
	}
}

func TestLayerLoggersIntegration(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "info",
		Format: "text",
	}

	factory := NewLoggerFactory(cfg)
	ctx := WithRequestID(context.Background(), "req-layer-test")

	// Test all layer loggers
	domainLogger := factory.GetDomainLogger()
	appLogger := factory.GetApplicationLogger()
	infraLogger := factory.GetInfrastructureLogger()
	interfaceLogger := factory.GetInterfaceLogger()

	// These should not panic and should work correctly
	domainLogger.LogEntityCreated(ctx, "User", "user-456")
	appLogger.LogUseCaseStarted(ctx, "CreateUser", map[string]interface{}{"email": "test@example.com"})
	infraLogger.LogDatabaseQuery(ctx, "SELECT * FROM users", time.Millisecond*10)
	interfaceLogger.LogHTTPRequest(ctx, "POST", "/query", 200, time.Millisecond*50)

	// If we get here without panicking, the integration is working
}

func TestConfigurationValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  config.LoggingConfig
		wantErr bool
	}{
		{
			name: "valid json config",
			config: config.LoggingConfig{
				Level:  "info",
				Format: "json",
			},
			wantErr: false,
		},
		{
			name: "valid text config",
			config: config.LoggingConfig{
				Level:  "debug",
				Format: "text",
			},
			wantErr: false,
		},
		{
			name: "invalid level",
			config: config.LoggingConfig{
				Level:  "invalid",
				Format: "json",
			},
			wantErr: true,
		},
		{
			name: "invalid format",
			config: config.LoggingConfig{
				Level:  "info",
				Format: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoggingConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			// If config is valid, test that logger factory can be created
			if !tt.wantErr {
				factory := NewLoggerFactory(tt.config)
				if factory == nil {
					t.Error("NewLoggerFactory returned nil for valid config")
				}
			}
		})
	}
}

func TestRequestTracingFlow(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "debug",
		Format: "json",
	}

	factory := NewLoggerFactory(cfg)

	// Simulate a request flow with tracing
	requestID := GenerateRequestID()
	correlationID := GenerateCorrelationID()

	ctx := context.Background()
	ctx = WithRequestID(ctx, requestID)
	ctx = WithCorrelationID(ctx, correlationID)

	// Test that context values are preserved through the flow
	if retrievedReqID := GetRequestID(ctx); retrievedReqID != requestID {
		t.Errorf("Request ID not preserved: expected %s, got %s", requestID, retrievedReqID)
	}

	if retrievedCorrID := GetCorrelationID(ctx); retrievedCorrID != correlationID {
		t.Errorf("Correlation ID not preserved: expected %s, got %s", correlationID, retrievedCorrID)
	}

	// Test that loggers can use the context
	logger := factory.GetLogger()
	enrichedLogger := logger.WithRequestID(ctx).WithCorrelationID(correlationID)

	// This should work without errors
	enrichedLogger.Info("Request processed successfully")
}

func TestTextFormatOutput(t *testing.T) {
	var buf bytes.Buffer

	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	baseLogger := slog.New(handler)
	logger := &Logger{Logger: baseLogger}

	ctx := WithRequestID(context.Background(), "req-text-test")

	logger.WithRequestID(ctx).WithComponent("test").Info("Text format test message")

	output := buf.String()

	// Verify text format contains expected elements
	if !strings.Contains(output, "Text format test message") {
		t.Error("Log output should contain the message")
	}
	if !strings.Contains(output, "component=test") {
		t.Error("Log output should contain the component field")
	}
	if !strings.Contains(output, "request_id=req-text-test") {
		t.Error("Log output should contain the request_id field")
	}
}
