package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/kitti12911/lib-util/v3/apperror"

	"grpc-sandbox/internal/database"

	"github.com/kitti12911/lib-util/v3/validator"
)

type patchData struct {
	userFields    map[string]any
	profileFields map[string]any
	addressFields map[string]any
	profile       CreateProfileParams
	address       CreateAddressParams
}

type userRepository interface {
	GetByID(ctx context.Context, id string) (*database.User, error)
	CreateUser(ctx context.Context, params CreateParams) (*database.User, error)
	CreateProfile(ctx context.Context, userID string, params CreateProfileParams) (*database.UserProfile, error)
	CreateAddress(ctx context.Context, userProfileID string, params CreateAddressParams) (*database.UserAddress, error)
	GetProfileIDByUserID(ctx context.Context, userID string) (string, error)
	UpdateUser(ctx context.Context, params UpdateParams) (int64, error)
	PatchUser(ctx context.Context, userID string, fields map[string]any) (int64, error)
	UpdateProfileByUserID(ctx context.Context, userID string, params CreateProfileParams) (int64, error)
	PatchProfileByUserID(ctx context.Context, userID string, fields map[string]any) (int64, error)
	UpdateAddressByProfileID(ctx context.Context, userProfileID string, params CreateAddressParams) (int64, error)
	PatchAddressByProfileID(ctx context.Context, userProfileID string, fields map[string]any) (int64, error)
	List(ctx context.Context, params ListParams) (*ListResult, error)
	DeleteAddressByProfileID(ctx context.Context, userProfileID string) error
	DeleteAddressByUserID(ctx context.Context, userID string) error
	DeleteProfileByUserID(ctx context.Context, userID string) error
	DeleteUser(ctx context.Context, userID string) (int64, error)
}

type databaseProvider interface {
	Transaction(ctx context.Context, fn func(context.Context) error) error
}

type Service struct {
	userRepository userRepository
	db             databaseProvider
	validator      *validator.Validator
}

func NewService(userRepository userRepository, db databaseProvider) *Service {
	return &Service{
		userRepository: userRepository,
		db:             db,
		validator:      validator.New(),
	}
}

func (s *Service) GetByID(ctx context.Context, params GetByIDParams) (*database.User, error) {
	if err := s.validator.Validate(params); err != nil {
		slog.WarnContext(ctx, "invalid get user request", "error", err)
		return nil, apperror.InvalidInput("id must be a valid uuid", err)
	}

	user, err := s.userRepository.GetByID(ctx, params.ID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get user", "error", err, "id", params.ID)
		return nil, apperror.Internal("failed to get user", err)
	}

	if user == nil {
		slog.WarnContext(ctx, "user not found", "id", params.ID)
		return nil, apperror.NotFound("user not found", nil)
	}

	return user, nil
}

func (s *Service) Create(ctx context.Context, params CreateParams) (string, error) {
	if err := s.validator.Validate(params); err != nil {
		slog.WarnContext(ctx, "invalid create user request", "error", err)
		return "", apperror.InvalidInput("invalid create user request", err)
	}

	var user *database.User
	if err := s.db.Transaction(ctx, func(ctx context.Context) error {
		var err error
		user, err = s.userRepository.CreateUser(ctx, params)
		if err != nil {
			return err
		}

		if params.Profile == nil {
			return nil
		}

		profile, err := s.userRepository.CreateProfile(ctx, user.ID, *params.Profile)
		if err != nil {
			return err
		}

		if params.Profile.Address == nil {
			return nil
		}

		_, err = s.userRepository.CreateAddress(ctx, profile.ID, *params.Profile.Address)
		return err
	}); err != nil {
		if err, ok := errors.AsType[*apperror.Error](err); ok {
			return "", err
		}

		slog.ErrorContext(ctx, "failed to create user", "error", err)
		return "", apperror.Internal("failed to create user", err)
	}

	return user.ID, nil
}

func (s *Service) List(ctx context.Context, params ListParams) (*ListResult, error) {
	result, err := s.userRepository.List(ctx, params)
	if err != nil {
		if err, ok := errors.AsType[*apperror.Error](err); ok {
			return nil, err
		}

		slog.ErrorContext(ctx, "failed to list users", "error", err)
		return nil, apperror.Internal("failed to list users", err)
	}

	return result, nil
}

func (s *Service) Update(ctx context.Context, params UpdateParams) (int64, error) {
	if err := s.validator.Validate(params); err != nil {
		slog.WarnContext(ctx, "invalid update user request", "error", err)
		return 0, apperror.InvalidInput("invalid update user request", err)
	}

	var affectedRows int64
	if err := s.db.Transaction(ctx, func(ctx context.Context) error {
		var err error
		affectedRows, err = s.userRepository.UpdateUser(ctx, params)
		if err != nil || affectedRows == 0 {
			return err
		}

		return s.replaceProfile(ctx, params.ID, params.Profile)
	}); err != nil {
		if err, ok := errors.AsType[*apperror.Error](err); ok {
			return 0, err
		}

		slog.ErrorContext(ctx, "failed to update user", "error", err, "id", params.ID)
		return 0, apperror.Internal("failed to update user", err)
	}

	if affectedRows == 0 {
		slog.WarnContext(ctx, "user not found", "id", params.ID)
		return 0, apperror.NotFound("user not found", nil)
	}

	return affectedRows, nil
}

func (s *Service) Patch(ctx context.Context, params PatchParams) (int64, error) {
	if err := s.validator.Validate(params); err != nil {
		slog.WarnContext(ctx, "invalid patch user request", "error", err)
		return 0, apperror.InvalidInput("invalid patch user request", err)
	}

	data := patchFields(params)
	if err := s.validatePatch(data); err != nil {
		slog.WarnContext(ctx, "invalid patch user request", "error", err)
		return 0, apperror.InvalidInput("invalid patch user request", err)
	}

	var affectedRows int64
	if err := s.db.Transaction(ctx, func(ctx context.Context) error {
		var err error
		affectedRows, err = s.userRepository.PatchUser(ctx, params.ID, data.userFields)
		if err != nil || affectedRows == 0 {
			return err
		}

		if len(data.profileFields) > 0 || len(data.addressFields) > 0 {
			return s.patchProfile(ctx, params.ID, data)
		}
		return nil
	}); err != nil {
		if err, ok := errors.AsType[*apperror.Error](err); ok {
			return 0, err
		}

		slog.ErrorContext(ctx, "failed to patch user", "error", err, "id", params.ID)
		return 0, apperror.Internal("failed to patch user", err)
	}

	if affectedRows == 0 {
		slog.WarnContext(ctx, "user not found", "id", params.ID)
		return 0, apperror.NotFound("user not found", nil)
	}

	return affectedRows, nil
}

func (s *Service) validatePatch(data patchData) error {
	if len(data.userFields)+len(data.profileFields)+len(data.addressFields) == 0 {
		return errors.New("update_mask has no supported fields")
	}
	if email, ok := data.userFields["email"]; ok {
		if err := s.validator.Engine().Var(email, "required,email"); err != nil {
			return err
		}
	}
	if username, ok := data.userFields["username"]; ok {
		if err := s.validator.Engine().Var(username, "required"); err != nil {
			return err
		}
	}
	if status, ok := data.userFields["status"]; ok {
		if err := s.validator.Engine().Var(status, "required,oneof=active disabled pending"); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) patchProfile(ctx context.Context, userID string, data patchData) error {
	profileID, err := s.userRepository.GetProfileIDByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if profileID == "" {
		var profile *database.UserProfile
		profile, err = s.userRepository.CreateProfile(ctx, userID, data.profile)
		if err != nil {
			return err
		}
		profileID = profile.ID
	} else if len(data.profileFields) > 0 {
		if _, patchErr := s.userRepository.PatchProfileByUserID(ctx, userID, data.profileFields); patchErr != nil {
			return patchErr
		}
	}

	if len(data.addressFields) == 0 {
		return nil
	}

	affectedRows, err := s.userRepository.PatchAddressByProfileID(ctx, profileID, data.addressFields)
	if err != nil {
		return err
	}
	if affectedRows > 0 {
		return nil
	}

	_, err = s.userRepository.CreateAddress(ctx, profileID, data.address)
	return err
}

func (s *Service) replaceProfile(ctx context.Context, userID string, params *CreateProfileParams) error {
	if params == nil {
		if err := s.userRepository.DeleteAddressByUserID(ctx, userID); err != nil {
			return err
		}
		return s.userRepository.DeleteProfileByUserID(ctx, userID)
	}

	profileID, err := s.userRepository.GetProfileIDByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if profileID == "" {
		var profile *database.UserProfile
		profile, err = s.userRepository.CreateProfile(ctx, userID, *params)
		if err != nil {
			return err
		}
		profileID = profile.ID
	} else if _, err := s.userRepository.UpdateProfileByUserID(ctx, userID, *params); err != nil {
		return err
	}

	return s.replaceAddress(ctx, profileID, params.Address)
}

func (s *Service) replaceAddress(ctx context.Context, profileID string, params *CreateAddressParams) error {
	if params == nil {
		return s.userRepository.DeleteAddressByProfileID(ctx, profileID)
	}

	affectedRows, err := s.userRepository.UpdateAddressByProfileID(ctx, profileID, *params)
	if err != nil {
		return err
	}
	if affectedRows > 0 {
		return nil
	}

	_, err = s.userRepository.CreateAddress(ctx, profileID, *params)
	return err
}

func (s *Service) Delete(ctx context.Context, params DeleteParams) (int64, error) {
	if err := s.validator.Validate(params); err != nil {
		slog.WarnContext(ctx, "invalid delete user request", "error", err)
		return 0, apperror.InvalidInput("id must be a valid uuid", err)
	}

	var affectedRows int64
	if err := s.db.Transaction(ctx, func(ctx context.Context) error {
		var err error
		affectedRows, err = s.userRepository.DeleteUser(ctx, params.ID)
		return err
	}); err != nil {
		if err, ok := errors.AsType[*apperror.Error](err); ok {
			return 0, err
		}

		slog.ErrorContext(ctx, "failed to delete user", "error", err, "id", params.ID)
		return 0, apperror.Internal("failed to delete user", err)
	}

	if affectedRows == 0 {
		slog.WarnContext(ctx, "user not found", "id", params.ID)
		return 0, apperror.NotFound("user not found", nil)
	}

	return affectedRows, nil
}
