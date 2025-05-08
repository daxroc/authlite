package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/daxroc/authlite/internal/logging"
	"github.com/google/uuid"
)

type contextKey string

const requestIDKey = contextKey("requestID")

func LoggingMiddleware(sensitive []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String() // or use time-based fallback
		}
		rw.Header().Set("X-Request-ID", reqID)

		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		next.ServeHTTP(rw, r.WithContext(ctx))

		logging.Info("http_request",
			"request_id", reqID,
			"remote", r.RemoteAddr,
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.statusCode,
			"duration", time.Since(start).Seconds(),
			"headers", redactHeaders(r.Header, sensitive),
		)
	})
}

func redactHeaders(h http.Header, sensitive []string) map[string]string {
	redacted := make(map[string]string)
	for k, v := range h {
		lk := strings.ToLower(k)
		if contains(lk, sensitive) {
			redacted[k] = "***REDACTED***"
		} else {
			redacted[k] = strings.Join(v, ",")
		}
	}
	return redacted
}

func contains(item string, list []string) bool {
	for _, val := range list {
		if val == item {
			return true
		}
	}
	return false
}

// Capture response status
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
