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

func (r *repository) List(ctx context.Context, params ListParams) (*ListResult, error) {
	users := make([]database.User, 0)

	query := r.db.NewSelect().
		Model(&users).
		Relation("Profile").
		Relation("Profile.Address")

	if err := applyFilters(query, params.Filters); err != nil {
		return nil, err
	}

	if err := applyOrderBy(query, params.OrderBy); err != nil {
		return nil, err
	}

	total, err := query.
		Limit(params.Limit).
		Offset(params.Offset).
		ScanAndCount(ctx)
	if err != nil {
		return nil, err
	}

	return &ListResult{
		Users: users,
		Total: int64(total),
	}, nil
}
