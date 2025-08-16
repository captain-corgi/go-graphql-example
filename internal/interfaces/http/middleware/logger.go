package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

// Logger creates a structured logging middleware using slog
func Logger(logger *slog.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Extract request ID from context
		requestID := ""
		if param.Request != nil {
			if id := GetRequestID(param.Request.Context()); id != "" {
				requestID = id
			}
		}

		// Log the request with structured logging
		logger.Info("HTTP request",
			slog.String("method", param.Method),
			slog.String("path", param.Path),
			slog.Int("status", param.StatusCode),
			slog.Duration("latency", param.Latency),
			slog.String("ip", param.ClientIP),
			slog.String("user_agent", param.Request.UserAgent()),
			slog.String("request_id", requestID),
			slog.Time("timestamp", param.TimeStamp),
		)

		// Return empty string since we're using slog for actual logging
		return ""
	})
}

// LoggerWithConfig creates a structured logging middleware with custom configuration
func LoggerWithConfig(logger *slog.Logger, config gin.LoggerConfig) gin.HandlerFunc {
	if config.Formatter == nil {
		config.Formatter = func(param gin.LogFormatterParams) string {
			requestID := ""
			if param.Request != nil {
				if id := GetRequestID(param.Request.Context()); id != "" {
					requestID = id
				}
			}

			// Skip logging for certain paths if configured
			if config.SkipPaths != nil {
				for _, path := range config.SkipPaths {
					if param.Path == path {
						return ""
					}
				}
			}

			logger.Info("HTTP request",
				slog.String("method", param.Method),
				slog.String("path", param.Path),
				slog.Int("status", param.StatusCode),
				slog.Duration("latency", param.Latency),
				slog.String("ip", param.ClientIP),
				slog.String("user_agent", param.Request.UserAgent()),
				slog.String("request_id", requestID),
				slog.Time("timestamp", param.TimeStamp),
			)

			return ""
		}
	}

	return gin.LoggerWithConfig(config)
}
