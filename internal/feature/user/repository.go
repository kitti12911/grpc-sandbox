package user

import (
	"context"
	"database/sql"
	"errors"

	"grpc-sandbox/internal/apperror"
	"grpc-sandbox/internal/database"

	orm "github.com/kitti12911/lib-orm/v2"
	"github.com/uptrace/bun/driver/pgdriver"
)

type repository struct {
	db *orm.DB
}

func NewRepository(db *orm.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetByID(ctx context.Context, id string) (*database.User, error) {
	user := new(database.User)
	if err := r.db.IDB(ctx).NewSelect().
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

func (r *repository) CreateUser(ctx context.Context, params CreateParams) (*database.User, error) {
	user := &database.User{
		Email:       params.Email,
		Username:    params.Username,
		DisplayName: params.DisplayName,
		Status:      params.Status,
	}

	if err := r.db.IDB(ctx).NewInsert().
		Model(user).
		Returning("id").
		Scan(ctx); err != nil {
		if isUniqueViolation(err) {
			return nil, apperror.AlreadyExist("user already exists", err)
		}
		return nil, err
	}

	return user, nil
}

func (r *repository) CreateProfile(
	ctx context.Context,
	userID string,
	params CreateProfileParams,
) (*database.UserProfile, error) {
	profile := &database.UserProfile{
		UserID:      userID,
		FirstName:   params.FirstName,
		LastName:    params.LastName,
		PhoneNumber: params.PhoneNumber,
	}

	if err := r.db.IDB(ctx).NewInsert().
		Model(profile).
		Returning("id").
		Scan(ctx); err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *repository) CreateAddress(
	ctx context.Context,
	userProfileID string,
	params CreateAddressParams,
) (*database.UserAddress, error) {
	address := &database.UserAddress{
		UserProfileID: userProfileID,
		Line1:         params.Line1,
		Line2:         params.Line2,
		City:          params.City,
		State:         params.State,
		PostalCode:    params.PostalCode,
		CountryCode:   params.CountryCode,
	}

	if err := r.db.IDB(ctx).NewInsert().
		Model(address).
		Returning("id").
		Scan(ctx); err != nil {
		return nil, err
	}

	return address, nil
}

func (r *repository) List(ctx context.Context, params ListParams) (*ListResult, error) {
	users := make([]database.User, 0)

	query := r.db.IDB(ctx).NewSelect().
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

func (r *repository) DeleteAddressesByUserID(ctx context.Context, userID string) error {
	_, err := r.db.IDB(ctx).NewDelete().
		Model((*database.UserAddress)(nil)).
		Where("user_profile_id IN (SELECT id FROM user_profiles WHERE user_id = ?)", userID).
		Exec(ctx)
	return err
}

func (r *repository) DeleteProfileByUserID(ctx context.Context, userID string) error {
	_, err := r.db.IDB(ctx).NewDelete().
		Model((*database.UserProfile)(nil)).
		Where("user_id = ?", userID).
		Exec(ctx)
	return err
}

func (r *repository) DeleteUser(ctx context.Context, userID string) (int64, error) {
	result, err := r.db.IDB(ctx).NewDelete().
		Model((*database.User)(nil)).
		Where("id = ?", userID).
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func isUniqueViolation(err error) bool {
	var pgErr pgdriver.Error
	return errors.As(err, &pgErr) && pgErr.Field('C') == "23505"
}
