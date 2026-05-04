package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kitti12911/lib-util/v3/apperror"

	fieldmap "grpc-sandbox/gen/database"
	"grpc-sandbox/internal/database"

	orm "github.com/kitti12911/lib-orm/v2"
	"github.com/uptrace/bun/driver/pgdriver"
)

var (
	userPatchColumns    = orm.WritableColumns(fieldmap.UserRootFields, "id", "created_at", "updated_at", "deleted_at")
	profilePatchColumns = orm.WritableColumns(fieldmap.UserProfileFields, "id", "user_id", "created_at", "updated_at")
	addressPatchColumns = orm.WritableColumns(fieldmap.UserAddressFields, "id", "user_profile_id", "created_at", "updated_at")
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

func (r *repository) GetProfileIDByUserID(ctx context.Context, userID string) (string, error) {
	var profileID string
	if err := r.db.IDB(ctx).NewSelect().
		Model((*database.UserProfile)(nil)).
		Column("id").
		Where("user_id = ?", userID).
		Scan(ctx, &profileID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return profileID, nil
}

func (r *repository) UpdateUser(ctx context.Context, params UpdateParams) (int64, error) {
	result, err := r.db.IDB(ctx).NewUpdate().
		Model((*database.User)(nil)).
		Set("email = ?", params.Email).
		Set("username = ?", params.Username).
		Set("display_name = ?", params.DisplayName).
		Set("status = ?", params.Status).
		Set("updated_at = now()").
		Where("id = ?", params.ID).
		Exec(ctx)
	if err != nil {
		if isUniqueViolation(err) {
			return 0, apperror.AlreadyExist("user already exists", err)
		}
		return 0, err
	}

	return result.RowsAffected()
}

func (r *repository) PatchUser(ctx context.Context, userID string, fields map[string]any) (int64, error) {
	query := r.db.IDB(ctx).NewUpdate().
		Model((*database.User)(nil)).
		Set("updated_at = now()").
		Where("id = ?", userID)

	if err := orm.ApplyPatchFields(query, fields, userPatchColumns); err != nil {
		return 0, err
	}

	result, err := query.Exec(ctx)
	if err != nil {
		if isUniqueViolation(err) {
			return 0, apperror.AlreadyExist("user already exists", err)
		}
		return 0, err
	}

	return result.RowsAffected()
}

func (r *repository) UpdateProfileByUserID(
	ctx context.Context,
	userID string,
	params CreateProfileParams,
) (int64, error) {
	result, err := r.db.IDB(ctx).NewUpdate().
		Model((*database.UserProfile)(nil)).
		Set("first_name = ?", params.FirstName).
		Set("last_name = ?", params.LastName).
		Set("phone_number = ?", params.PhoneNumber).
		Set("updated_at = now()").
		Where("user_id = ?", userID).
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (r *repository) PatchProfileByUserID(
	ctx context.Context,
	userID string,
	fields map[string]any,
) (int64, error) {
	query := r.db.IDB(ctx).NewUpdate().
		Model((*database.UserProfile)(nil)).
		Set("updated_at = now()").
		Where("user_id = ?", userID)

	if err := orm.ApplyPatchFields(query, fields, profilePatchColumns); err != nil {
		return 0, err
	}

	result, err := query.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (r *repository) UpdateAddressByProfileID(
	ctx context.Context,
	userProfileID string,
	params CreateAddressParams,
) (int64, error) {
	result, err := r.db.IDB(ctx).NewUpdate().
		Model((*database.UserAddress)(nil)).
		Set("line1 = ?", params.Line1).
		Set("line2 = ?", params.Line2).
		Set("city = ?", params.City).
		Set("state = ?", params.State).
		Set("postal_code = ?", params.PostalCode).
		Set("country_code = ?", params.CountryCode).
		Set("updated_at = now()").
		Where("user_profile_id = ?", userProfileID).
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (r *repository) PatchAddressByProfileID(
	ctx context.Context,
	userProfileID string,
	fields map[string]any,
) (int64, error) {
	query := r.db.IDB(ctx).NewUpdate().
		Model((*database.UserAddress)(nil)).
		Set("updated_at = now()").
		Where("user_profile_id = ?", userProfileID)

	if err := orm.ApplyPatchFields(query, fields, addressPatchColumns); err != nil {
		return 0, err
	}

	result, err := query.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
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

func (r *repository) DeleteAddressByProfileID(ctx context.Context, userProfileID string) error {
	_, err := r.db.IDB(ctx).NewDelete().
		Model((*database.UserAddress)(nil)).
		Where("user_profile_id = ?", userProfileID).
		Exec(ctx)
	return err
}

func (r *repository) DeleteAddressByUserID(ctx context.Context, userID string) error {
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
		Model(&database.User{}).
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
