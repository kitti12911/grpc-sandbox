package main

import (
	"context"
	"grpc-sandbox/internal/config"
	"grpc-sandbox/internal/database"
	"grpc-sandbox/internal/feature/user"
	"grpc-sandbox/internal/server"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kitti12911/lib-monitor/profiling"
	"github.com/kitti12911/lib-monitor/tracing"
	libconfig "github.com/kitti12911/lib-util/v3/config"
	"github.com/kitti12911/lib-util/v3/logger"

	"github.com/dromara/carbon/v2"
)

func main() {
	ctx := context.Background()

	// Set default time
	carbon.SetDefault(carbon.Default{
		Layout:       carbon.RFC3339Format,
		Timezone:     carbon.Bangkok,
		WeekStartsAt: carbon.Sunday,
		Locale:       "en",
	})

	// Load config
	cfg, err := libconfig.Load[config.Config]("config.yml")

	if err != nil {
		slog.ErrorContext(ctx, "failed to load config", "error", err)
		os.Exit(1)
	}

	if cfg.Service.ShutdownTimeout == 0 {
		cfg.Service.ShutdownTimeout = 10 * time.Second
	}

	// Init logger
	logger.NewFromConfig(cfg.Logging, cfg.Service.Name)

	// Init monitoring
	profiler, err := profiling.NewFromConfig(cfg.Service.Name, cfg.Profiling)
	if err != nil {
		slog.ErrorContext(ctx, "failed to init profiling", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := profiling.Shutdown(profiler); err != nil {
			slog.ErrorContext(ctx, "failed to stop profiling", "error", err)
		}
	}()

	// Init tracing
	tp, err := tracing.NewFromConfig(ctx, cfg.Service.Name, cfg.Tracing)
	if err != nil {
		slog.ErrorContext(ctx, "failed to init tracing", "error", err)
		os.Exit(1)
	}
	defer tracing.Shutdown(ctx, tp)

	// Init database
	db, err := database.New(ctx, cfg)
	if err != nil {
		slog.ErrorContext(ctx, "failed to init database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.ErrorContext(ctx, "failed to close database", "error", err)
		}
	}()

	// Init repositories, services, and handlers
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository, db)
	userHandler := user.NewHandler(userService)

	// Start gRPC server
	srv, err := server.NewGRPCServer(cfg.Service.Port, userHandler)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create gRPC server", "error", err)
		os.Exit(1)
	}

	go func() {
		if err := srv.Start(); err != nil {
			slog.ErrorContext(ctx, "gRPC server error", "error", err)
			os.Exit(1)
		}
	}()

	slog.InfoContext(ctx, "gRPC server started", "port", cfg.Service.Port)

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.InfoContext(ctx, "shutting down gRPC server")

	shutdownCtx, cancel := context.WithTimeout(ctx, cfg.Service.ShutdownTimeout)
	defer cancel()

	srv.Stop(shutdownCtx)

	slog.InfoContext(ctx, "server stopped")
}
