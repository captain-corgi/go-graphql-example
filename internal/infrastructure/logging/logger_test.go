package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/captain-corgi/go-graphql-example/internal/infrastructure/config"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name   string
		config config.LoggingConfig
	}{
		{
			name: "JSON format with info level",
			config: config.LoggingConfig{
				Level:  "info",
				Format: "json",
			},
		},
		{
			name: "Text format with debug level",
			config: config.LoggingConfig{
				Level:  "debug",
				Format: "text",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.config)
			if logger == nil {
				t.Error("NewLogger returned nil")
				return
			}
			if logger.Logger == nil {
				t.Error("Logger.Logger is nil")
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"DEBUG", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"INFO", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"WARN", slog.LevelWarn},
		{"error", slog.LevelError},
		{"ERROR", slog.LevelError},
		{"invalid", slog.LevelInfo}, // default
		{"", slog.LevelInfo},        // default
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseLevel(tt.input)
			if result != tt.expected {
				t.Errorf("parseLevel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLoggerWithRequestID(t *testing.T) {
	var buf bytes.Buffer

	// Create a logger that writes to our buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	baseLogger := slog.New(handler)
	logger := &Logger{Logger: baseLogger}

	// Create context with request ID
	ctx := WithRequestID(context.Background(), "test-request-123")

	// Log with request ID
	logger.WithRequestID(ctx).Info("test message")

	// Parse the JSON output
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	// Check that request_id is included
	if requestID, ok := logEntry["request_id"].(string); !ok || requestID != "test-request-123" {
		t.Errorf("Expected request_id 'test-request-123', got %v", logEntry["request_id"])
	}

	// Check that message is correct
	if msg, ok := logEntry["msg"].(string); !ok || msg != "test message" {
		t.Errorf("Expected message 'test message', got %v", logEntry["msg"])
	}
}

func TestLoggerWithCorrelationID(t *testing.T) {
	var buf bytes.Buffer

	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	baseLogger := slog.New(handler)
	logger := &Logger{Logger: baseLogger}

	// Log with correlation ID
	logger.WithCorrelationID("corr-456").Info("test message")

	// Parse the JSON output
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	// Check that correlation_id is included
	if corrID, ok := logEntry["correlation_id"].(string); !ok || corrID != "corr-456" {
		t.Errorf("Expected correlation_id 'corr-456', got %v", logEntry["correlation_id"])
	}
}

func TestLoggerWithComponent(t *testing.T) {
	var buf bytes.Buffer

	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	baseLogger := slog.New(handler)
	logger := &Logger{Logger: baseLogger}

	// Log with component
	logger.WithComponent("test-component").Info("test message")

	// Parse the JSON output
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	// Check that component is included
	if component, ok := logEntry["component"].(string); !ok || component != "test-component" {
		t.Errorf("Expected component 'test-component', got %v", logEntry["component"])
	}
}

func TestLoggerWithFields(t *testing.T) {
	var buf bytes.Buffer

	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	baseLogger := slog.New(handler)
	logger := &Logger{Logger: baseLogger}

	// Log with multiple fields
	fields := map[string]interface{}{
		"user_id": "123",
		"action":  "create",
		"count":   42,
	}
	logger.WithFields(fields).Info("test message")

	// Parse the JSON output
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	// Check that all fields are included
	if userID, ok := logEntry["user_id"].(string); !ok || userID != "123" {
		t.Errorf("Expected user_id '123', got %v", logEntry["user_id"])
	}
	if action, ok := logEntry["action"].(string); !ok || action != "create" {
		t.Errorf("Expected action 'create', got %v", logEntry["action"])
	}
	if count, ok := logEntry["count"].(float64); !ok || count != 42 {
		t.Errorf("Expected count 42, got %v", logEntry["count"])
	}
}

func TestLoggerChaining(t *testing.T) {
	var buf bytes.Buffer

	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	baseLogger := slog.New(handler)
	logger := &Logger{Logger: baseLogger}

	ctx := WithRequestID(context.Background(), "req-123")

	// Chain multiple context additions
	logger.WithRequestID(ctx).
		WithCorrelationID("corr-456").
		WithComponent("test").
		Info("chained message")

	// Parse the JSON output
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	// Check that all chained values are included
	expectedFields := map[string]string{
		"request_id":     "req-123",
		"correlation_id": "corr-456",
		"component":      "test",
		"msg":            "chained message",
	}

	for field, expected := range expectedFields {
		if actual, ok := logEntry[field].(string); !ok || actual != expected {
			t.Errorf("Expected %s '%s', got %v", field, expected, logEntry[field])
		}
	}
}

func TestLoggerTextFormat(t *testing.T) {
	var buf bytes.Buffer

	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	baseLogger := slog.New(handler)
	logger := &Logger{Logger: baseLogger}

	logger.WithComponent("test").Info("text format message")

	output := buf.String()
	if !strings.Contains(output, "text format message") {
		t.Error("Log output should contain the message")
	}
	if !strings.Contains(output, "component=test") {
		t.Error("Log output should contain the component field")
	}
}
