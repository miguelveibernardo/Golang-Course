package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-cli/api"
)

// Ensuring that the middlewar creates a traceID, adds it to the context, and sets it in the response header
func TestTraceMiddleware_AddsTraceID(t *testing.T) {
	var capturedTraceID string

	//Dummy handler to read TraceID context
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedTraceID = api.GetTraceID(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler := api.TraceMiddleware(next)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	//Checks if the trace ID was added to the context
	if capturedTraceID == "" || capturedTraceID == "no-trace" {
		t.Errorf("Expected valid trace ID, got %q", capturedTraceID)
	}

	//checks if TraceID header is set in the response
	headerTraceID := w.Header().Get("X-Trace-ID")
	if headerTraceID == "" {
		t.Error("Expected X-Trace-ID header to be set, but is missing")
	}

	//Ensures header and context mach
	if headerTraceID != capturedTraceID {
		t.Errorf("TraceID mismatch: header=%q, context=%q", headerTraceID, capturedTraceID)
	}
}

// Validates that the requests reah the next handler
func TestTraceMiddleware_CallNextHandler(t *testing.T) {
	called := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler := api.TraceMiddleware(next)
	handler.ServeHTTP(w, req)
	if !called {
		t.Error("Expected next handler to be called, but it was not")
	}
}
