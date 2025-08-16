package logging

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/captain-corgi/go-graphql-example/internal/infrastructure/config"
)

// Logger wraps slog.Logger with additional functionality
type Logger struct {
	*slog.Logger
}

// NewLogger creates a new structured logger based on configuration
func NewLogger(cfg config.LoggingConfig) *Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:     parseLevel(cfg.Level),
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize attribute formatting if needed
			return a
		},
	}

	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	return &Logger{Logger: logger}
}

// parseLevel converts string level to slog.Level
func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// WithRequestID adds request ID to logger context
func (l *Logger) WithRequestID(ctx context.Context) *Logger {
	if requestID := GetRequestID(ctx); requestID != "" {
		return &Logger{Logger: l.Logger.With("request_id", requestID)}
	}
	return l
}

// WithCorrelationID adds correlation ID to logger context
func (l *Logger) WithCorrelationID(correlationID string) *Logger {
	return &Logger{Logger: l.Logger.With("correlation_id", correlationID)}
}

// WithComponent adds component name to logger context
func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{Logger: l.Logger.With("component", component)}
}

// WithError adds error to logger context
func (l *Logger) WithError(err error) *Logger {
	return &Logger{Logger: l.Logger.With("error", err.Error())}
}

// WithFields adds multiple fields to logger context
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{Logger: l.Logger.With(args...)}
}
