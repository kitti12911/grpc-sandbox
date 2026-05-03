package user

import (
	"grpc-sandbox/internal/database"

	orm "github.com/kitti12911/lib-orm"
)

type GetByIDParams struct {
	ID string `validate:"required,uuid"`
}

type ListParams struct {
	Limit   int
	Offset  int
	Filters []orm.Filter
	OrderBy []orm.OrderBy
}

type ListResult struct {
	Users []database.User
	Total int64
}
