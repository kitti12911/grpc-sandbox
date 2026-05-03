package user

import (
	"context"
	"database/sql"
	"errors"

	"grpc-sandbox/internal/database"

	"github.com/uptrace/bun"
)

type repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetByID(ctx context.Context, id string) (*database.User, error) {
	user := new(database.User)
	if err := r.db.NewSelect().
		Model(user).
		Relation("Profile").
		Relation("Profile.Address").
		Where("u.id = ?", id).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
