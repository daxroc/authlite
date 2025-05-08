package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestStore() *TokenStore {
	store := &TokenStore{
		tokens: map[string]struct{}{
			"token-valid": {},
		},
	}
	return store
}

func TestHandler_AuthorizationHeader(t *testing.T) {
	store := setupTestStore()
	handler := NewHandler(store, "Authorization")

	req := httptest.NewRequest("GET", "/auth", nil)
	req.Header.Set("Authorization", "Bearer token-valid")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", w.Code)
	}
}

func TestHandler_MissingHeader(t *testing.T) {
	store := setupTestStore()
	handler := NewHandler(store, "Authorization")

	req := httptest.NewRequest("GET", "/auth", nil)
	// No Authorization header

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized, got %d", w.Code)
	}
}

func TestHandler_CustomHeader(t *testing.T) {
	store := setupTestStore()
	handler := NewHandler(store, "X-Auth-Token")

	req := httptest.NewRequest("GET", "/auth", nil)
	req.Header.Set("X-Auth-Token", "token-valid")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK with custom header, got %d", w.Code)
	}
}

func TestHandler_InvalidToken(t *testing.T) {
	store := setupTestStore()
	handler := NewHandler(store, "Authorization")

	req := httptest.NewRequest("GET", "/auth", nil)
	req.Header.Set("Authorization", "Bearer wrong-token")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for invalid token, got %d", w.Code)
	}
}
