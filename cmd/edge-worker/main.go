package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"industrial-edge-protocol/internal/config"
	"industrial-edge-protocol/internal/service"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	app := service.NewApplication(logger, cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger.Info("edge worker started", "interval", cfg.WorkerInterval)
	ticker := time.NewTicker(cfg.WorkerInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := app.RunCycle(ctx); err != nil {
				logger.Error("worker cycle failed", "error", err)
			}
		}
	}
}
