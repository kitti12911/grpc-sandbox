package user

import (
	"grpc-sandbox/internal/database"

	orm "github.com/kitti12911/lib-orm/v3"
)

type GetByIDParams struct {
	ID string `validate:"required,uuid"`
}

type DeleteParams struct {
	ID string `validate:"required,uuid"`
}

type CreateParams struct {
	Email       string               `field:"email" validate:"required,email"`
	Username    string               `field:"username" validate:"required"`
	DisplayName *string              `field:"display_name"`
	Status      string               `field:"status" validate:"required,oneof=active disabled pending"`
	Profile     *CreateProfileParams `field:"profile"`
}

type CreateProfileParams struct {
	FirstName   *string              `field:"first_name"`
	LastName    *string              `field:"last_name"`
	PhoneNumber *string              `field:"phone_number"`
	Address     *CreateAddressParams `field:"address"`
}

type CreateAddressParams struct {
	Line1       *string `field:"line1"`
	Line2       *string `field:"line2"`
	City        *string `field:"city"`
	State       *string `field:"state"`
	PostalCode  *string `field:"postal_code"`
	CountryCode *string `field:"country_code"`
}

type UpdateParams struct {
	ID          string `validate:"required,uuid"`
	Email       string `validate:"required,email"`
	Username    string `validate:"required"`
	DisplayName *string
	Status      string `validate:"required,oneof=active disabled pending"`
	Profile     *CreateProfileParams
}

type PatchParams struct {
	ID     string       `validate:"required,uuid"`
	User   CreateParams `validate:"-"`
	Fields []string
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
