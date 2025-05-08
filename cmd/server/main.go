package main

import (
	"log"
	"net/http"

	"github.com/daxroc/authlite/internal/auth"
	"github.com/daxroc/authlite/internal/logging"
	"github.com/daxroc/authlite/internal/metrics"
)

func main() {
	config := auth.LoadConfig()
	if err := RunServer(config); err != nil {
		log.Fatalf("[authlite] %v", err)
	}
}

func RunServer(config auth.Config) error {
	logging.SetFormat(logging.Format(config.LogFormat))
	store := auth.NewTokenStore(config.TokenFile)
	if err := store.Load(); err != nil {
		return err
	}
	go auth.PeriodicReload(store, config.ReloadInterval)
	go auth.ReloadOnSignal(store)
	metrics.Register()

	mux := http.NewServeMux()
	mux.Handle("/auth", auth.NewHandler(store, config.TokenHeader))
	mux.Handle("/metrics", metrics.Handler())

	wrapped := auth.LoggingMiddleware(config.SensitiveHeaders, mux)

	logging.Info("starting server", "addr", config.ListenAddr)
	return http.ListenAndServe(config.ListenAddr, wrapped)
}
