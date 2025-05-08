package auth

import (
	"flag"
	"strings"
	"time"
)

type Config struct {
	TokenFile        string
	ReloadInterval   time.Duration
	ListenAddr       string
	TokenHeader      string
	LogFormat        string
	SensitiveHeaders []string
}

func LoadConfig() Config {
	var rawHeaders string

	cfg := Config{}
	flag.StringVar(&cfg.TokenFile, "token-file", "/secrets/tokens.txt", "Path to token file")
	flag.DurationVar(&cfg.ReloadInterval, "reload-interval", 5*time.Minute, "Token reload interval")
	flag.StringVar(&cfg.ListenAddr, "listen", ":8080", "HTTP listen address")
	flag.StringVar(&cfg.TokenHeader, "token-header", "Authorization", "Header to extract bearer token from")
	flag.StringVar(&cfg.LogFormat, "log-format", "logfmt", "Log format: logfmt or json")
	flag.StringVar(&rawHeaders, "sensitive-headers", "authorization,x-auth-token,x-api-key", "Comma-separated list of sensitive headers to redact")
	flag.Parse()

	cfg.SensitiveHeaders = parseCSV(rawHeaders)
	return cfg
}

func parseCSV(raw string) []string {
	fields := strings.Split(raw, ",")
	for i := range fields {
		fields[i] = strings.TrimSpace(strings.ToLower(fields[i]))
	}
	return fields
}
