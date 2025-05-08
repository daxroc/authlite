package auth

import (
	"flag"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestLoadConfig_WithSensitiveHeaders(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{
		"authlite",
		"--token-file=/tmp/test.tokens",
		"--reload-interval=15s",
		"--listen=:9000",
		"--token-header=X-Auth-Token",
		"--sensitive-headers=secret,x-custom",
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cfg := LoadConfig()

	expected := Config{
		TokenFile:        "/tmp/test.tokens",
		ReloadInterval:   15 * time.Second,
		ListenAddr:       ":9000",
		TokenHeader:      "X-Auth-Token",
		LogFormat:        "logfmt",
		SensitiveHeaders: []string{"secret", "x-custom"},
	}

	if !reflect.DeepEqual(cfg, expected) {
		t.Errorf("unexpected config: got %+v, want %+v", cfg, expected)
	}
}
