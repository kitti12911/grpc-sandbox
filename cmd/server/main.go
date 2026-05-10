package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"grpc-sandbox/internal/config"
	"grpc-sandbox/internal/database"
	"grpc-sandbox/internal/feature/user"
	"grpc-sandbox/internal/feature/worker"
	"grpc-sandbox/internal/server"

	async "github.com/kitti12911/lib-async"
	"github.com/kitti12911/lib-monitor/profiling"
	"github.com/kitti12911/lib-monitor/tracing"
	libconfig "github.com/kitti12911/lib-util/v3/config"
	"github.com/kitti12911/lib-util/v3/logger"

	"github.com/dromara/carbon/v2"
)

func main() {
	os.Exit(run())
}

func run() int {
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
		return 1
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
		return 1
	}
	defer func() {
		if shutdownErr := profiling.Shutdown(profiler); shutdownErr != nil {
			slog.ErrorContext(ctx, "failed to stop profiling", "error", shutdownErr)
		}
	}()

	// Init tracing
	tp, err := tracing.NewFromConfig(ctx, cfg.Service.Name, cfg.Tracing)
	if err != nil {
		slog.ErrorContext(ctx, "failed to init tracing", "error", err)
		return 1
	}
	defer func() {
		if shutdownErr := tracing.Shutdown(ctx, tp); shutdownErr != nil {
			slog.ErrorContext(ctx, "failed to stop tracing", "error", shutdownErr)
		}
	}()

	// Init async bus
	bus, err := async.NewNATS(cfg.NATS, nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to nats", "error", err)
		return 1
	}
	defer func() {
		if closeErr := bus.Close(); closeErr != nil {
			slog.ErrorContext(ctx, "failed to close nats bus", "error", closeErr)
		}
	}()

	// Init database
	db, err := database.New(ctx, cfg)
	if err != nil {
		slog.ErrorContext(ctx, "failed to init database", "error", err)
		return 1
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			slog.ErrorContext(ctx, "failed to close database", "error", closeErr)
		}
	}()

	// Init repositories, services, and handlers
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository, db)
	userHandler := user.NewHandler(userService)
	workerService := worker.NewService(bus, cfg.Worker.Topic)
	workerHandler := worker.NewHandler(workerService)

	// Start gRPC server
	srv, err := server.NewGRPCServer(ctx, cfg.Service.Port, userHandler, workerHandler)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create gRPC server", "error", err)
		return 1
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- srv.Start()
	}()

	slog.InfoContext(ctx, "gRPC server started", "port", cfg.Service.Port)

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	case err := <-serverErr:
		slog.ErrorContext(ctx, "gRPC server error", "error", err)
		return 1
	}

	slog.InfoContext(ctx, "shutting down gRPC server")

	shutdownCtx, cancel := context.WithTimeout(ctx, cfg.Service.ShutdownTimeout)
	defer cancel()

	srv.Stop(shutdownCtx)

	slog.InfoContext(ctx, "server stopped")

	return 0
}
