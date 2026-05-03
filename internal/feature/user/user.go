package user

import (
	"grpc-sandbox/internal/database"

	orm "github.com/kitti12911/lib-orm/v2"
)

type GetByIDParams struct {
	ID string `validate:"required,uuid"`
}

type DeleteParams struct {
	ID string `validate:"required,uuid"`
}

type CreateParams struct {
	Email       string `validate:"required,email"`
	Username    string `validate:"required"`
	DisplayName *string
	Status      string `validate:"required,oneof=active disabled pending"`
	Profile     *CreateProfileParams
}

type CreateProfileParams struct {
	FirstName   *string
	LastName    *string
	PhoneNumber *string
	Address     *CreateAddressParams
}

type CreateAddressParams struct {
	Line1       *string
	Line2       *string
	City        *string
	State       *string
	PostalCode  *string
	CountryCode *string
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
