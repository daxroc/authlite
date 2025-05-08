package auth

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRedactHeaders(t *testing.T) {
	h := http.Header{
		"Authorization": {"Bearer abc"},
		"X-Api-Key":     {"123"},
		"X-Trace-Id":    {"trace123"},
	}

	sensitive := []string{"authorization", "x-api-key"}
	redacted := redactHeaders(h, sensitive)

	expected := map[string]string{
		"Authorization": "***REDACTED***",
		"X-Api-Key":     "***REDACTED***",
		"X-Trace-Id":    "trace123",
	}

	if !reflect.DeepEqual(redacted, expected) {
		t.Errorf("unexpected redaction: got %+v, want %+v", redacted, expected)
	}
}

func TestLoggingMiddleware_LogsRedactedHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("X-Api-Key", "secret")

	rw := httptest.NewRecorder()

	h := LoggingMiddleware([]string{"authorization", "x-api-key"}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rw.Code)
	}
}
