package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/daxroc/authlite/internal/logging"
	"github.com/daxroc/authlite/internal/metrics"
)

func NewHandler(store *TokenStore, headerName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		raw := r.Header.Get(headerName)
		reqID := r.Context().Value(requestIDKey)

		status := "unauthorized"
		var token string
		if strings.EqualFold(headerName, "Authorization") && strings.HasPrefix(raw, "Bearer ") {
			token = strings.TrimPrefix(raw, "Bearer ")
		} else if !strings.EqualFold(headerName, "Authorization") {
			token = strings.TrimSpace(raw)
		}

		valid := token != "" && store.IsValid(token)
		if valid {
			status = "ok"
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		metrics.AuthRequests.WithLabelValues(status).Inc()
		metrics.AuthDuration.WithLabelValues(status).Observe(time.Since(start).Seconds())

		logging.Info("token_validation",
			"request_id", reqID,
			"token_provided", token != "",
			"token_valid", valid,
			"status", status,
		)
	})
}
