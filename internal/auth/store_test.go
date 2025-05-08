package auth

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestTokenStore_IsValid(t *testing.T) {
	tmpFile := "test_token_store.tokens"
	os.WriteFile(tmpFile, []byte("token1\ntoken2\ntoken3\n"), 0644)
	defer os.Remove(tmpFile)

	store := NewTokenStore(tmpFile)
	if err := store.Load(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		token   string
		isValid bool
	}{
		{"token1", true},
		{"tokenX", false},
	}

	for _, tt := range tests {
		if got := store.IsValid(tt.token); got != tt.isValid {
			t.Errorf("IsValid(%q) = %v; want %v", tt.token, got, tt.isValid)
		}
	}
}

func TestPeriodicReload(t *testing.T) {
	tmpFile := "test_reload.tokens"
	os.WriteFile(tmpFile, []byte("token-a\n"), 0644)
	defer os.Remove(tmpFile)

	store := NewTokenStore(tmpFile)
	store.Load()

	go PeriodicReload(store, 100*time.Millisecond)
	time.Sleep(150 * time.Millisecond)
	os.WriteFile(tmpFile, []byte("token-b\n"), 0644)

	time.Sleep(300 * time.Millisecond)
	if !store.IsValid("token-b") {
		t.Fatal("token-b should be valid after reload")
	}
}

func TestReloadOnSignal(t *testing.T) {
	tmpFile := "test_signal_reload.tokens"
	os.WriteFile(tmpFile, []byte("token-x\n"), 0644)
	defer os.Remove(tmpFile)

	store := NewTokenStore(tmpFile)
	store.Load()
	go ReloadOnSignal(store)

	os.WriteFile(tmpFile, []byte("token-y\n"), 0644)

	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("failed to find process: %v", err)
	}
	if err := process.Signal(syscall.SIGHUP); err != nil {
		t.Fatalf("failed to send SIGHUP: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	if !store.IsValid("token-y") {
		t.Fatal("token-y should be valid after signal reload")
	}
}
