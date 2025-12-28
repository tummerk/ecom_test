package application

import (
	"context"
	"ecom_test/internal/config"
	"ecom_test/internal/domain/service"
	"ecom_test/internal/infrastructure/persistance"
	"ecom_test/internal/server"
	"ecom_test/pkg/application/modules"
	"ecom_test/pkg/contextx"
	"log"
	"log/slog"
	"os/signal"
	"syscall"
)

func Run(cfg config.Config) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	ctx = contextx.WithLogger(ctx, slog.Default())

	repository := persistance.NewTaskRepository()
	service := service.NewTaskService(repository)
	server := server.NewServer(cfg, service)

	httpModule := modules.HTTPServer{ShutdownTimeout: cfg.ShutdownTimeout}

	logger(ctx).Info("start http server")
	if err := httpModule.Run(ctx, server); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
	logger(ctx).Info("application stopped successfully")
}
