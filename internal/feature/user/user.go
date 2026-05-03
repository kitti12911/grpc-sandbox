package user

type GetByIDParams struct {
	ID string `validate:"required,uuid"`
}
