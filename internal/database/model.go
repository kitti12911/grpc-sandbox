package database

import (
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID          string       `bun:"id,pk,type:uuid,default:uuidv7()"`
	Email       string       `bun:"email,notnull"`
	Username    string       `bun:"username,notnull"`
	DisplayName *string      `bun:"display_name"`
	Status      string       `bun:"status,notnull"`
	CreatedAt   time.Time    `bun:"created_at,notnull,default:now()"`
	UpdatedAt   time.Time    `bun:"updated_at,notnull,default:now()"`
	DeletedAt   sql.NullTime `bun:"deleted_at,soft_delete,nullzero"`

	Profile *UserProfile `bun:"rel:has-one,join:id=user_id"`
}

type UserProfile struct {
	bun.BaseModel `bun:"table:user_profiles,alias:up"`

	ID          string    `bun:"id,pk,type:uuid,default:uuidv7()"`
	UserID      string    `bun:"user_id,type:uuid,notnull"`
	FirstName   *string   `bun:"first_name"`
	LastName    *string   `bun:"last_name"`
	PhoneNumber *string   `bun:"phone_number"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:now()"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:now()"`

	User    *User        `bun:"rel:belongs-to,join:user_id=id"`
	Address *UserAddress `bun:"rel:has-one,join:id=user_profile_id"`
}

type UserAddress struct {
	bun.BaseModel `bun:"table:user_addresses,alias:ua"`

	ID            string    `bun:"id,pk,type:uuid,default:uuidv7()"`
	UserProfileID string    `bun:"user_profile_id,type:uuid,notnull"`
	Line1         *string   `bun:"line1"`
	Line2         *string   `bun:"line2"`
	City          *string   `bun:"city"`
	State         *string   `bun:"state"`
	PostalCode    *string   `bun:"postal_code"`
	CountryCode   *string   `bun:"country_code"`
	CreatedAt     time.Time `bun:"created_at,notnull,default:now()"`
	UpdatedAt     time.Time `bun:"updated_at,notnull,default:now()"`

	Profile *UserProfile `bun:"rel:belongs-to,join:user_profile_id=id"`
}

func models() []any {
	return []any{
		(*User)(nil),
		(*UserProfile)(nil),
		(*UserAddress)(nil),
	}
}
