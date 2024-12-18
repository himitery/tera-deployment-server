package main

import (
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"tera/deployment/pkg/config"
	"tera/deployment/pkg/logger"
	"time"
)

func main() {
	app := NewApp(config.Load())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		logger.Error("failed to start server", zap.Error(err))

		os.Exit(1)
	}

	<-app.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := app.Stop(shutdownCtx); err != nil {
		logger.Error("failed to stop server", zap.Error(err))
	}
}
