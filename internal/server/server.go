package server

import (
	"ecom_test/internal/config"
	"ecom_test/pkg/middlewarex"
	"net/http"
	"time"
)

func NewServer(cfg config.Config, service TaskService) *http.Server {
	mux := http.NewServeMux()
	handler := NewTaskHandler(service)
	handler.RegisterRoutes(mux)

	var wrappedMux http.Handler = middlewarex.Logger(mux)

	return &http.Server{
		Addr:         cfg.Addr,
		Handler:      wrappedMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}
