package config

import (
	"time"

	"github.com/kitti12911/lib-monitor/profiling"
	"github.com/kitti12911/lib-monitor/tracing"
	liborm "github.com/kitti12911/lib-orm"
	"github.com/kitti12911/lib-util/v2/logger"
)

type Config struct {
	Service   Service          `mapstructure:"service" validate:"required"`
	Logging   logger.Config    `mapstructure:"logging"`
	Tracing   tracing.Config   `mapstructure:"tracing"`
	Profiling profiling.Config `mapstructure:"profiling"`
	Database  liborm.Config    `mapstructure:"database" validate:"required"`
}

type Service struct {
	Name            string        `mapstructure:"name"             env:"SERVICE_NAME"      validate:"required"`
	Port            int           `mapstructure:"port"             env:"PORT"              validate:"required,gte=1,lte=65535"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT"`
}
