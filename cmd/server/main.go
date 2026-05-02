package main

import (
	"context"
	"grpc-sandbox/internal/config"
	"log/slog"
	"os"
	"time"

	"github.com/kitti12911/lib-monitor/profiling"
	"github.com/kitti12911/lib-monitor/tracing"
	orm "github.com/kitti12911/lib-orm"
	libconfig "github.com/kitti12911/lib-util/v2/config"
	"github.com/kitti12911/lib-util/v2/logger"

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
	db, err := orm.New(
		ctx,
		cfg.Database,
		orm.WithApplicationName(cfg.Service.Name),
		orm.WithTracing(cfg.Tracing.Enabled),
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.ErrorContext(ctx, "failed to close database", "error", err)
		}
	}()
}
