package user

import (
	fieldmap "grpc-sandbox/gen/database"
	"grpc-sandbox/internal/apperror"

	orm "github.com/kitti12911/lib-orm"
	"github.com/uptrace/bun"
)

func applyFilters(query *bun.SelectQuery, filters []orm.Filter) error {
	if err := orm.ApplyFilters(query, filters, fieldmap.UserColumns); err != nil {
		return apperror.InvalidInput("invalid filters", err)
	}

	return nil
}

func applyOrderBy(query *bun.SelectQuery, orderBy []orm.OrderBy) error {
	if err := orm.ApplyOrderBy(query, orderBy, fieldmap.UserColumns); err != nil {
		return apperror.InvalidInput("invalid order by", err)
	}

	return nil
}
