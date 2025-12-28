package modules

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type HTTPServer struct {
	ShutdownTimeout time.Duration
}

func (h HTTPServer) Run(
	ctx context.Context,
	httpServer *http.Server,
) error {
	if httpServer == nil {
		return errors.New("http server is nil")
	}

	serverError := make(chan error, 1)

	go func() {
		logger(ctx).Info("http server started", slog.String("address", httpServer.Addr))

		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverError <- fmt.Errorf("httpServer.ListenAndServe: %w", err)
		}
		close(serverError)
	}()

	select {
	case <-ctx.Done():
		logger(ctx).Info("shutting down http server", slog.String("address", httpServer.Addr))

		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), h.ShutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("http server shutdown failed: %w", err)
		}

		logger(ctx).Info("http server stopped gracefully")
		return nil

	case err := <-serverError:
		return err
	}
}
