package main

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/daxroc/authlite/internal/auth"
)

func TestRunServer_ValidatesTokenWithCustomSensitiveHeaders(t *testing.T) {
	tmpFile := "test_server.tokens"
	os.WriteFile(tmpFile, []byte("tok123\n"), 0644)
	defer os.Remove(tmpFile)

	cfg := auth.Config{
		TokenFile:        tmpFile,
		ReloadInterval:   time.Minute,
		ListenAddr:       ":18081",
		TokenHeader:      "Authorization",
		LogFormat:        "logfmt",
		SensitiveHeaders: []string{"authorization"},
	}

	go func() {
		_ = RunServer(cfg)
	}()

	time.Sleep(300 * time.Millisecond)

	req, _ := http.NewRequest("GET", "http://localhost:18081/auth", nil)
	req.Header.Set("X-Request-ID", "test-req-id")
	req.Header.Set("Authorization", "Bearer tok123")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}
