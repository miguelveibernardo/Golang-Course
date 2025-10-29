package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
)

type contextkey string

const traceIDKey contextkey = "traceID"

// creates a random 8-byte hex string
func generateTraceID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "unknown-trace"
	}
	return hex.EncodeToString(b)
}

// wraps an http.Handler to inject a TraceID into the request context
func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := generateTraceID()
		ctx := context.WithValue(r.Context(), traceIDKey, traceID)
		r = r.WithContext(ctx)

		slog.Info("Incoming request",
			"method", r.Method,
			"path", r.URL.Path,
			"trace_id", traceID,
		)
		w.Header().Set("X-Trace-ID", traceID)
		next.ServeHTTP(w, r)

		slog.Info("Request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"trace_id", traceID,
		)
	})
}

// extract the TraceID from context,when present
func GetTraceID(ctx context.Context) string {
	traceID, ok := ctx.Value(traceIDKey).(string)
	if !ok {
		return "no-trace"
	}
	return traceID
}
