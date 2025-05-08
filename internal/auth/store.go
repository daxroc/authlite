package auth

import (
	"bufio"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type TokenStore struct {
	mu     sync.RWMutex
	tokens map[string]struct{}
	path   string
}

func NewTokenStore(path string) *TokenStore {
	return &TokenStore{
		path:   path,
		tokens: make(map[string]struct{}),
	}
}

func (s *TokenStore) Load() error {
	file, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer file.Close()

	tokens := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			tokens[line] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	s.tokens = tokens
	s.mu.Unlock()

	log.Printf("[authlite] loaded %d tokens", len(tokens))
	return nil
}

func (s *TokenStore) IsValid(token string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.tokens[token]
	return exists
}

func PeriodicReload(s *TokenStore, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		log.Println("[authlite] periodic reload triggered")
		if err := s.Load(); err != nil {
			log.Printf("[authlite] error reloading tokens: %v", err)
		}
	}
}

func ReloadOnSignal(s *TokenStore) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP)
	for range ch {
		log.Println("[authlite] SIGHUP received, reloading tokens")
		if err := s.Load(); err != nil {
			log.Printf("[authlite] error reloading tokens: %v", err)
		}
	}
}
