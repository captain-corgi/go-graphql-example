package logging

import (
	"context"
	"strings"
	"testing"
)

func TestWithRequestID(t *testing.T) {
	ctx := context.Background()
	requestID := "test-request-123"

	newCtx := WithRequestID(ctx, requestID)

	if retrievedID := GetRequestID(newCtx); retrievedID != requestID {
		t.Errorf("Expected request ID %s, got %s", requestID, retrievedID)
	}
}

func TestGetRequestID_Empty(t *testing.T) {
	ctx := context.Background()

	if requestID := GetRequestID(ctx); requestID != "" {
		t.Errorf("Expected empty request ID, got %s", requestID)
	}
}

func TestWithCorrelationID(t *testing.T) {
	ctx := context.Background()
	correlationID := "test-correlation-456"

	newCtx := WithCorrelationID(ctx, correlationID)

	if retrievedID := GetCorrelationID(newCtx); retrievedID != correlationID {
		t.Errorf("Expected correlation ID %s, got %s", correlationID, retrievedID)
	}
}

func TestGetCorrelationID_Empty(t *testing.T) {
	ctx := context.Background()

	if correlationID := GetCorrelationID(ctx); correlationID != "" {
		t.Errorf("Expected empty correlation ID, got %s", correlationID)
	}
}

func TestGenerateRequestID(t *testing.T) {
	id1 := GenerateRequestID()
	id2 := GenerateRequestID()

	// Check that IDs are different
	if id1 == id2 {
		t.Error("Generated request IDs should be unique")
	}

	// Check that IDs have the correct prefix
	if !strings.HasPrefix(id1, "req-") {
		t.Errorf("Request ID should start with 'req-', got %s", id1)
	}
	if !strings.HasPrefix(id2, "req-") {
		t.Errorf("Request ID should start with 'req-', got %s", id2)
	}

	// Check that IDs have reasonable length (prefix + 16 hex chars)
	expectedLength := len("req-") + 16
	if len(id1) != expectedLength {
		t.Errorf("Expected request ID length %d, got %d", expectedLength, len(id1))
	}
}

func TestGenerateCorrelationID(t *testing.T) {
	id1 := GenerateCorrelationID()
	id2 := GenerateCorrelationID()

	// Check that IDs are different
	if id1 == id2 {
		t.Error("Generated correlation IDs should be unique")
	}

	// Check that IDs have the correct prefix
	if !strings.HasPrefix(id1, "corr-") {
		t.Errorf("Correlation ID should start with 'corr-', got %s", id1)
	}
	if !strings.HasPrefix(id2, "corr-") {
		t.Errorf("Correlation ID should start with 'corr-', got %s", id2)
	}

	// Check that IDs have reasonable length (prefix + 32 hex chars)
	expectedLength := len("corr-") + 32
	if len(id1) != expectedLength {
		t.Errorf("Expected correlation ID length %d, got %d", expectedLength, len(id1))
	}
}

func TestContextChaining(t *testing.T) {
	ctx := context.Background()
	requestID := "req-123"
	correlationID := "corr-456"

	// Chain context additions
	ctx = WithRequestID(ctx, requestID)
	ctx = WithCorrelationID(ctx, correlationID)

	// Verify both values are preserved
	if retrievedReqID := GetRequestID(ctx); retrievedReqID != requestID {
		t.Errorf("Expected request ID %s, got %s", requestID, retrievedReqID)
	}
	if retrievedCorrID := GetCorrelationID(ctx); retrievedCorrID != correlationID {
		t.Errorf("Expected correlation ID %s, got %s", correlationID, retrievedCorrID)
	}
}

func TestContextOverwrite(t *testing.T) {
	ctx := context.Background()

	// Set initial request ID
	ctx = WithRequestID(ctx, "req-first")

	// Overwrite with new request ID
	ctx = WithRequestID(ctx, "req-second")

	// Should have the new value
	if requestID := GetRequestID(ctx); requestID != "req-second" {
		t.Errorf("Expected request ID 'req-second', got %s", requestID)
	}
}
