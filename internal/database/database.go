package database

import (
	"context"

	"grpc-sandbox/internal/config"
	"grpc-sandbox/internal/database/migrations"
	"grpc-sandbox/internal/database/seeders"

	orm "github.com/kitti12911/lib-orm/v2"
)

func New(ctx context.Context, cfg *config.Config) (*orm.DB, error) {
	db, err := orm.New(
		ctx,
		cfg.Database,
		orm.WithApplicationName(cfg.Service.Name),
		orm.WithModels(models()...),
		orm.WithTracing(cfg.Tracing.Enabled),
	)
	if err != nil {
		return nil, err
	}

	if err := orm.Init(ctx, db, cfg.Database, migrations.Migrations, seeders.Fixtures, seeders.Users); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
