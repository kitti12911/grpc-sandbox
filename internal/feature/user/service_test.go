package user

import (
	"context"
	"errors"
	"testing"

	"github.com/kitti12911/lib-util/v3/apperror"

	"grpc-sandbox/internal/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubUserRepository struct {
	getByIDFunc                func(ctx context.Context, id string) (*database.User, error)
	listFunc                   func(ctx context.Context, params ListParams) (*ListResult, error)
	createUserFunc             func(ctx context.Context, params CreateParams) (*database.User, error)
	createProfileFunc          func(ctx context.Context, userID string, params CreateProfileParams) (*database.UserProfile, error)
	createAddressFunc          func(ctx context.Context, userProfileID string, params CreateAddressParams) (*database.UserAddress, error)
	getProfileIDByUserIDFunc   func(ctx context.Context, userID string) (string, error)
	updateUserFunc             func(ctx context.Context, params UpdateParams) (int64, error)
	updateProfileByUserIDFunc  func(ctx context.Context, userID string, params CreateProfileParams) (int64, error)
	updateAddressByProfileFunc func(ctx context.Context, profileID string, params CreateAddressParams) (int64, error)
	patchUserFunc              func(ctx context.Context, userID string, fields map[string]any) (int64, error)
	patchProfileByUserIDFunc   func(ctx context.Context, userID string, fields map[string]any) (int64, error)
	patchAddressByProfileFunc  func(ctx context.Context, profileID string, fields map[string]any) (int64, error)
	deleteUserFunc             func(ctx context.Context, userID string) (int64, error)
	deleteAddressByProfileFunc func(ctx context.Context, profileID string) error
	deleteAddressByUserIDFunc  func(ctx context.Context, userID string) error
	deleteProfileByUserIDFunc  func(ctx context.Context, userID string) error
}

func (r stubUserRepository) GetByID(ctx context.Context, id string) (*database.User, error) {
	if r.getByIDFunc != nil {
		return r.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (r stubUserRepository) CreateUser(ctx context.Context, params CreateParams) (*database.User, error) {
	if r.createUserFunc == nil {
		return &database.User{}, nil
	}
	return r.createUserFunc(ctx, params)
}

func (r stubUserRepository) CreateProfile(
	ctx context.Context,
	userID string,
	params CreateProfileParams,
) (*database.UserProfile, error) {
	if r.createProfileFunc != nil {
		return r.createProfileFunc(ctx, userID, params)
	}
	return &database.UserProfile{}, nil
}

func (r stubUserRepository) CreateAddress(
	ctx context.Context,
	userProfileID string,
	params CreateAddressParams,
) (*database.UserAddress, error) {
	if r.createAddressFunc != nil {
		return r.createAddressFunc(ctx, userProfileID, params)
	}
	return &database.UserAddress{}, nil
}

func (r stubUserRepository) GetProfileIDByUserID(ctx context.Context, userID string) (string, error) {
	if r.getProfileIDByUserIDFunc != nil {
		return r.getProfileIDByUserIDFunc(ctx, userID)
	}
	return "", nil
}

func (r stubUserRepository) UpdateUser(ctx context.Context, params UpdateParams) (int64, error) {
	if r.updateUserFunc == nil {
		return 0, nil
	}
	return r.updateUserFunc(ctx, params)
}

func (r stubUserRepository) PatchUser(ctx context.Context, userID string, fields map[string]any) (int64, error) {
	if r.patchUserFunc == nil {
		return 0, nil
	}
	return r.patchUserFunc(ctx, userID, fields)
}

func (r stubUserRepository) UpdateProfileByUserID(
	ctx context.Context,
	userID string,
	params CreateProfileParams,
) (int64, error) {
	if r.updateProfileByUserIDFunc != nil {
		return r.updateProfileByUserIDFunc(ctx, userID, params)
	}
	return 0, nil
}

func (r stubUserRepository) PatchProfileByUserID(
	ctx context.Context,
	userID string,
	fields map[string]any,
) (int64, error) {
	if r.patchProfileByUserIDFunc != nil {
		return r.patchProfileByUserIDFunc(ctx, userID, fields)
	}
	return 0, nil
}

func (r stubUserRepository) UpdateAddressByProfileID(
	ctx context.Context,
	profileID string,
	params CreateAddressParams,
) (int64, error) {
	if r.updateAddressByProfileFunc != nil {
		return r.updateAddressByProfileFunc(ctx, profileID, params)
	}
	return 0, nil
}

func (r stubUserRepository) PatchAddressByProfileID(
	ctx context.Context,
	profileID string,
	fields map[string]any,
) (int64, error) {
	if r.patchAddressByProfileFunc != nil {
		return r.patchAddressByProfileFunc(ctx, profileID, fields)
	}
	return 0, nil
}

func (r stubUserRepository) List(ctx context.Context, params ListParams) (*ListResult, error) {
	if r.listFunc != nil {
		return r.listFunc(ctx, params)
	}
	return nil, nil
}

func (r stubUserRepository) DeleteAddressByProfileID(ctx context.Context, profileID string) error {
	if r.deleteAddressByProfileFunc != nil {
		return r.deleteAddressByProfileFunc(ctx, profileID)
	}
	return nil
}

func (r stubUserRepository) DeleteAddressByUserID(ctx context.Context, userID string) error {
	if r.deleteAddressByUserIDFunc != nil {
		return r.deleteAddressByUserIDFunc(ctx, userID)
	}
	return nil
}

func (r stubUserRepository) DeleteProfileByUserID(ctx context.Context, userID string) error {
	if r.deleteProfileByUserIDFunc != nil {
		return r.deleteProfileByUserIDFunc(ctx, userID)
	}
	return nil
}

func (r stubUserRepository) DeleteUser(ctx context.Context, userID string) (int64, error) {
	if r.deleteUserFunc == nil {
		return 0, nil
	}
	return r.deleteUserFunc(ctx, userID)
}

type stubTransactionProvider struct{}

func (stubTransactionProvider) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func requireAppError(t *testing.T, err error) *apperror.Error {
	t.Helper()

	var appErr *apperror.Error
	require.True(t, errors.As(err, &appErr))

	return appErr
}

func TestServiceCreate(t *testing.T) {
	service := NewService(stubUserRepository{
		createUserFunc: func(_ context.Context, params CreateParams) (*database.User, error) {
			assert.Equal(t, "kit@example.com", params.Email)
			return &database.User{ID: "0198f8f0-0000-7000-8000-000000000999"}, nil
		},
	}, stubTransactionProvider{})

	id, err := service.Create(context.Background(), CreateParams{
		Email:    "kit@example.com",
		Username: "kit",
		Status:   "pending",
	})

	require.NoError(t, err)
	assert.Equal(t, "0198f8f0-0000-7000-8000-000000000999", id)
}

func TestServiceCreateValidatesRequest(t *testing.T) {
	called := false
	service := NewService(stubUserRepository{
		createUserFunc: func(context.Context, CreateParams) (*database.User, error) {
			called = true
			return nil, nil
		},
	}, stubTransactionProvider{})

	_, err := service.Create(context.Background(), CreateParams{})

	require.Error(t, err)
	appErr := requireAppError(t, err)
	assert.Equal(t, apperror.CodeInvalidInput, appErr.Code())
	assert.False(t, called)
}

func TestServiceCreatePassesThroughAppError(t *testing.T) {
	want := apperror.AlreadyExist("user already exists", nil)
	service := NewService(stubUserRepository{
		createUserFunc: func(context.Context, CreateParams) (*database.User, error) {
			return nil, want
		},
	}, stubTransactionProvider{})

	_, err := service.Create(context.Background(), CreateParams{
		Email:    "kit@example.com",
		Username: "kit",
		Status:   "pending",
	})

	require.ErrorIs(t, err, want)
}

func TestServiceUpdate(t *testing.T) {
	service := NewService(stubUserRepository{
		updateUserFunc: func(_ context.Context, params UpdateParams) (int64, error) {
			assert.Equal(t, "0198f8f0-0000-7000-8000-000000000999", params.ID)
			assert.Equal(t, "kit@example.com", params.Email)
			return 1, nil
		},
	}, stubTransactionProvider{})

	affectedRows, err := service.Update(context.Background(), UpdateParams{
		ID:       "0198f8f0-0000-7000-8000-000000000999",
		Email:    "kit@example.com",
		Username: "kit",
		Status:   "active",
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1), affectedRows)
}

func TestServiceUpdateValidatesRequest(t *testing.T) {
	called := false
	service := NewService(stubUserRepository{
		updateUserFunc: func(context.Context, UpdateParams) (int64, error) {
			called = true
			return 0, nil
		},
	}, stubTransactionProvider{})

	_, err := service.Update(context.Background(), UpdateParams{})

	require.Error(t, err)
	appErr := requireAppError(t, err)
	assert.Equal(t, apperror.CodeInvalidInput, appErr.Code())
	assert.False(t, called)
}

func TestServiceUpdateReturnsNotFound(t *testing.T) {
	service := NewService(stubUserRepository{
		updateUserFunc: func(context.Context, UpdateParams) (int64, error) {
			return 0, nil
		},
	}, stubTransactionProvider{})

	_, err := service.Update(context.Background(), UpdateParams{
		ID:       "0198f8f0-0000-7000-8000-000000000999",
		Email:    "kit@example.com",
		Username: "kit",
		Status:   "active",
	})

	require.Error(t, err)
	appErr := requireAppError(t, err)
	assert.Equal(t, apperror.CodeNotFound, appErr.Code())
}

func TestServicePatch(t *testing.T) {
	service := NewService(stubUserRepository{
		patchUserFunc: func(_ context.Context, userID string, fields map[string]any) (int64, error) {
			assert.Equal(t, "0198f8f0-0000-7000-8000-000000000999", userID)
			assert.Equal(t, "patched@example.com", fields["email"])
			assert.Equal(t, "active", fields["status"])
			return 1, nil
		},
	}, stubTransactionProvider{})

	affectedRows, err := service.Patch(context.Background(), PatchParams{
		ID:     "0198f8f0-0000-7000-8000-000000000999",
		User:   CreateParams{Email: "patched@example.com", Status: "active"},
		Fields: []string{"email", "status"},
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1), affectedRows)
}

func TestServicePatchCreatesMissingProfileAndAddress(t *testing.T) {
	service := NewService(stubUserRepository{
		patchUserFunc: func(context.Context, string, map[string]any) (int64, error) {
			return 1, nil
		},
		getProfileIDByUserIDFunc: func(context.Context, string) (string, error) {
			return "", nil
		},
		createProfileFunc: func(_ context.Context, userID string, _ CreateProfileParams) (*database.UserProfile, error) {
			assert.Equal(t, "0198f8f0-0000-7000-8000-000000000999", userID)
			return &database.UserProfile{ID: "0198f8f0-0000-7000-8000-000000000aaa"}, nil
		},
		patchAddressByProfileFunc: func(_ context.Context, profileID string, fields map[string]any) (int64, error) {
			assert.Equal(t, "0198f8f0-0000-7000-8000-000000000aaa", profileID)
			city, ok := fields["city"].(*string)
			require.True(t, ok)
			assert.Equal(t, "Chiang Mai", *city)
			return 0, nil
		},
		createAddressFunc: func(_ context.Context, profileID string, params CreateAddressParams) (*database.UserAddress, error) {
			assert.Equal(t, "0198f8f0-0000-7000-8000-000000000aaa", profileID)
			assert.Equal(t, "Chiang Mai", *params.City)
			return &database.UserAddress{}, nil
		},
	}, stubTransactionProvider{})

	affectedRows, err := service.Patch(context.Background(), PatchParams{
		ID: "0198f8f0-0000-7000-8000-000000000999",
		User: CreateParams{
			Profile: &CreateProfileParams{
				Address: &CreateAddressParams{
					City: new("Chiang Mai"),
				},
			},
		},
		Fields: []string{"profile.address.city"},
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1), affectedRows)
}

func TestServicePatchValidatesRequest(t *testing.T) {
	called := false
	service := NewService(stubUserRepository{
		patchUserFunc: func(context.Context, string, map[string]any) (int64, error) {
			called = true
			return 0, nil
		},
	}, stubTransactionProvider{})

	_, err := service.Patch(context.Background(), PatchParams{
		ID: "0198f8f0-0000-7000-8000-000000000999",
	})

	require.Error(t, err)
	appErr := requireAppError(t, err)
	assert.Equal(t, apperror.CodeInvalidInput, appErr.Code())
	assert.False(t, called)
}

func TestServicePatchReturnsNotFound(t *testing.T) {
	service := NewService(stubUserRepository{
		patchUserFunc: func(context.Context, string, map[string]any) (int64, error) {
			return 0, nil
		},
	}, stubTransactionProvider{})

	_, err := service.Patch(context.Background(), PatchParams{
		ID:     "0198f8f0-0000-7000-8000-000000000999",
		User:   CreateParams{Email: "patched@example.com"},
		Fields: []string{"email"},
	})

	require.Error(t, err)
	appErr := requireAppError(t, err)
	assert.Equal(t, apperror.CodeNotFound, appErr.Code())
}

func TestServiceDelete(t *testing.T) {
	service := NewService(stubUserRepository{
		deleteUserFunc: func(_ context.Context, userID string) (int64, error) {
			assert.Equal(t, "0198f8f0-0000-7000-8000-000000000999", userID)
			return 1, nil
		},
	}, stubTransactionProvider{})

	affectedRows, err := service.Delete(context.Background(), DeleteParams{
		ID: "0198f8f0-0000-7000-8000-000000000999",
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1), affectedRows)
}

func TestServiceDeleteValidatesRequest(t *testing.T) {
	called := false
	service := NewService(stubUserRepository{
		deleteUserFunc: func(context.Context, string) (int64, error) {
			called = true
			return 0, nil
		},
	}, stubTransactionProvider{})

	_, err := service.Delete(context.Background(), DeleteParams{})

	require.Error(t, err)
	appErr := requireAppError(t, err)
	assert.Equal(t, apperror.CodeInvalidInput, appErr.Code())
	assert.False(t, called)
}

func TestServiceDeleteReturnsNotFound(t *testing.T) {
	service := NewService(stubUserRepository{
		deleteUserFunc: func(context.Context, string) (int64, error) {
			return 0, nil
		},
	}, stubTransactionProvider{})

	_, err := service.Delete(context.Background(), DeleteParams{
		ID: "0198f8f0-0000-7000-8000-000000000999",
	})

	require.Error(t, err)
	appErr := requireAppError(t, err)
	assert.Equal(t, apperror.CodeNotFound, appErr.Code())
}

func TestServiceGetByID(t *testing.T) {
	want := &database.User{ID: "0198f8f0-0000-7000-8000-000000000999"}
	service := NewService(stubUserRepository{
		getByIDFunc: func(_ context.Context, id string) (*database.User, error) {
			assert.Equal(t, "0198f8f0-0000-7000-8000-000000000999", id)
			return want, nil
		},
	}, stubTransactionProvider{})

	got, err := service.GetByID(context.Background(), GetByIDParams{ID: "0198f8f0-0000-7000-8000-000000000999"})
	require.NoError(t, err)
	assert.Same(t, want, got)
}

func TestServiceGetByIDValidatesRequest(t *testing.T) {
	service := NewService(stubUserRepository{}, stubTransactionProvider{})

	_, err := service.GetByID(context.Background(), GetByIDParams{ID: "not-a-uuid"})
	require.Error(t, err)
	assert.Equal(t, apperror.CodeInvalidInput, requireAppError(t, err).Code())
}

func TestServiceGetByIDReturnsNotFound(t *testing.T) {
	service := NewService(stubUserRepository{
		getByIDFunc: func(context.Context, string) (*database.User, error) {
			return nil, nil
		},
	}, stubTransactionProvider{})

	_, err := service.GetByID(context.Background(), GetByIDParams{ID: "0198f8f0-0000-7000-8000-000000000999"})
	require.Error(t, err)
	assert.Equal(t, apperror.CodeNotFound, requireAppError(t, err).Code())
}

func TestServiceGetByIDWrapsRepoError(t *testing.T) {
	service := NewService(stubUserRepository{
		getByIDFunc: func(context.Context, string) (*database.User, error) {
			return nil, errors.New("boom")
		},
	}, stubTransactionProvider{})

	_, err := service.GetByID(context.Background(), GetByIDParams{ID: "0198f8f0-0000-7000-8000-000000000999"})
	require.Error(t, err)
	assert.Equal(t, apperror.CodeInternal, requireAppError(t, err).Code())
}

func TestServiceList(t *testing.T) {
	want := &ListResult{Users: []database.User{{ID: "u1"}}}
	service := NewService(stubUserRepository{
		listFunc: func(context.Context, ListParams) (*ListResult, error) {
			return want, nil
		},
	}, stubTransactionProvider{})

	got, err := service.List(context.Background(), ListParams{})
	require.NoError(t, err)
	assert.Same(t, want, got)
}

func TestServiceListPassesThroughAppError(t *testing.T) {
	want := apperror.InvalidInput("invalid filters", nil)
	service := NewService(stubUserRepository{
		listFunc: func(context.Context, ListParams) (*ListResult, error) {
			return nil, want
		},
	}, stubTransactionProvider{})

	_, err := service.List(context.Background(), ListParams{})
	require.ErrorIs(t, err, want)
}

func TestServiceListWrapsRepoError(t *testing.T) {
	service := NewService(stubUserRepository{
		listFunc: func(context.Context, ListParams) (*ListResult, error) {
			return nil, errors.New("boom")
		},
	}, stubTransactionProvider{})

	_, err := service.List(context.Background(), ListParams{})
	require.Error(t, err)
	assert.Equal(t, apperror.CodeInternal, requireAppError(t, err).Code())
}

func TestServiceUpdateReplacesExistingProfileAndAddress(t *testing.T) {
	updatedProfile := false
	updatedAddress := false
	service := NewService(stubUserRepository{
		updateUserFunc: func(context.Context, UpdateParams) (int64, error) {
			return 1, nil
		},
		getProfileIDByUserIDFunc: func(context.Context, string) (string, error) {
			return "0198f8f0-0000-7000-8000-000000000aaa", nil
		},
		updateProfileByUserIDFunc: func(_ context.Context, userID string, _ CreateProfileParams) (int64, error) {
			assert.Equal(t, "0198f8f0-0000-7000-8000-000000000999", userID)
			updatedProfile = true
			return 1, nil
		},
		updateAddressByProfileFunc: func(_ context.Context, profileID string, _ CreateAddressParams) (int64, error) {
			assert.Equal(t, "0198f8f0-0000-7000-8000-000000000aaa", profileID)
			updatedAddress = true
			return 1, nil
		},
	}, stubTransactionProvider{})

	_, err := service.Update(context.Background(), UpdateParams{
		ID:       "0198f8f0-0000-7000-8000-000000000999",
		Email:    "kit@example.com",
		Username: "kit",
		Status:   "active",
		Profile: &CreateProfileParams{
			Address: &CreateAddressParams{City: new("Bangkok")},
		},
	})
	require.NoError(t, err)
	assert.True(t, updatedProfile)
	assert.True(t, updatedAddress)
}

func TestServiceUpdateCreatesProfileWhenMissing(t *testing.T) {
	createdProfile := false
	createdAddress := false
	service := NewService(stubUserRepository{
		updateUserFunc: func(context.Context, UpdateParams) (int64, error) {
			return 1, nil
		},
		getProfileIDByUserIDFunc: func(context.Context, string) (string, error) {
			return "", nil
		},
		createProfileFunc: func(context.Context, string, CreateProfileParams) (*database.UserProfile, error) {
			createdProfile = true
			return &database.UserProfile{ID: "0198f8f0-0000-7000-8000-000000000aaa"}, nil
		},
		updateAddressByProfileFunc: func(context.Context, string, CreateAddressParams) (int64, error) {
			return 0, nil
		},
		createAddressFunc: func(context.Context, string, CreateAddressParams) (*database.UserAddress, error) {
			createdAddress = true
			return &database.UserAddress{}, nil
		},
	}, stubTransactionProvider{})

	_, err := service.Update(context.Background(), UpdateParams{
		ID:       "0198f8f0-0000-7000-8000-000000000999",
		Email:    "kit@example.com",
		Username: "kit",
		Status:   "active",
		Profile: &CreateProfileParams{
			Address: &CreateAddressParams{City: new("Bangkok")},
		},
	})
	require.NoError(t, err)
	assert.True(t, createdProfile)
	assert.True(t, createdAddress)
}

func TestServiceUpdateDeletesProfileWhenAbsent(t *testing.T) {
	deletedAddress := false
	deletedProfile := false
	service := NewService(stubUserRepository{
		updateUserFunc: func(context.Context, UpdateParams) (int64, error) {
			return 1, nil
		},
		deleteAddressByUserIDFunc: func(context.Context, string) error {
			deletedAddress = true
			return nil
		},
		deleteProfileByUserIDFunc: func(context.Context, string) error {
			deletedProfile = true
			return nil
		},
	}, stubTransactionProvider{})

	_, err := service.Update(context.Background(), UpdateParams{
		ID:       "0198f8f0-0000-7000-8000-000000000999",
		Email:    "kit@example.com",
		Username: "kit",
		Status:   "active",
	})
	require.NoError(t, err)
	assert.True(t, deletedAddress)
	assert.True(t, deletedProfile)
}

func TestServiceUpdateDeletesAddressWhenAbsentFromExistingProfile(t *testing.T) {
	deletedAddress := false
	service := NewService(stubUserRepository{
		updateUserFunc: func(context.Context, UpdateParams) (int64, error) {
			return 1, nil
		},
		getProfileIDByUserIDFunc: func(context.Context, string) (string, error) {
			return "0198f8f0-0000-7000-8000-000000000aaa", nil
		},
		updateProfileByUserIDFunc: func(context.Context, string, CreateProfileParams) (int64, error) {
			return 1, nil
		},
		deleteAddressByProfileFunc: func(_ context.Context, profileID string) error {
			assert.Equal(t, "0198f8f0-0000-7000-8000-000000000aaa", profileID)
			deletedAddress = true
			return nil
		},
	}, stubTransactionProvider{})

	_, err := service.Update(context.Background(), UpdateParams{
		ID:       "0198f8f0-0000-7000-8000-000000000999",
		Email:    "kit@example.com",
		Username: "kit",
		Status:   "active",
		Profile:  &CreateProfileParams{FirstName: new("Kit")},
	})
	require.NoError(t, err)
	assert.True(t, deletedAddress)
}

func TestServiceCreateWithProfileAndAddress(t *testing.T) {
	createdProfile := false
	createdAddress := false
	service := NewService(stubUserRepository{
		createUserFunc: func(context.Context, CreateParams) (*database.User, error) {
			return &database.User{ID: "0198f8f0-0000-7000-8000-000000000999"}, nil
		},
		createProfileFunc: func(context.Context, string, CreateProfileParams) (*database.UserProfile, error) {
			createdProfile = true
			return &database.UserProfile{ID: "0198f8f0-0000-7000-8000-000000000aaa"}, nil
		},
		createAddressFunc: func(context.Context, string, CreateAddressParams) (*database.UserAddress, error) {
			createdAddress = true
			return &database.UserAddress{}, nil
		},
	}, stubTransactionProvider{})

	_, err := service.Create(context.Background(), CreateParams{
		Email:    "kit@example.com",
		Username: "kit",
		Status:   "active",
		Profile: &CreateProfileParams{
			FirstName: new("Kit"),
			Address:   &CreateAddressParams{City: new("Bangkok")},
		},
	})
	require.NoError(t, err)
	assert.True(t, createdProfile)
	assert.True(t, createdAddress)
}

func TestServiceCreateWrapsNonAppError(t *testing.T) {
	service := NewService(stubUserRepository{
		createUserFunc: func(context.Context, CreateParams) (*database.User, error) {
			return nil, errors.New("boom")
		},
	}, stubTransactionProvider{})

	_, err := service.Create(context.Background(), CreateParams{
		Email: "kit@example.com", Username: "kit", Status: "active",
	})
	require.Error(t, err)
	assert.Equal(t, apperror.CodeInternal, requireAppError(t, err).Code())
}

func TestServiceUpdateWrapsNonAppError(t *testing.T) {
	service := NewService(stubUserRepository{
		updateUserFunc: func(context.Context, UpdateParams) (int64, error) {
			return 0, errors.New("boom")
		},
	}, stubTransactionProvider{})

	_, err := service.Update(context.Background(), UpdateParams{
		ID: "0198f8f0-0000-7000-8000-000000000999", Email: "k@example.com", Username: "k", Status: "active",
	})
	require.Error(t, err)
	assert.Equal(t, apperror.CodeInternal, requireAppError(t, err).Code())
}

func TestServicePatchWrapsNonAppError(t *testing.T) {
	service := NewService(stubUserRepository{
		patchUserFunc: func(context.Context, string, map[string]any) (int64, error) {
			return 0, errors.New("boom")
		},
	}, stubTransactionProvider{})

	_, err := service.Patch(context.Background(), PatchParams{
		ID:     "0198f8f0-0000-7000-8000-000000000999",
		User:   CreateParams{Email: "kit@example.com"},
		Fields: []string{"email"},
	})
	require.Error(t, err)
	assert.Equal(t, apperror.CodeInternal, requireAppError(t, err).Code())
}

func TestServiceDeleteWrapsNonAppError(t *testing.T) {
	service := NewService(stubUserRepository{
		deleteUserFunc: func(context.Context, string) (int64, error) {
			return 0, errors.New("boom")
		},
	}, stubTransactionProvider{})

	_, err := service.Delete(context.Background(), DeleteParams{ID: "0198f8f0-0000-7000-8000-000000000999"})
	require.Error(t, err)
	assert.Equal(t, apperror.CodeInternal, requireAppError(t, err).Code())
}

func TestServicePatchRejectsInvalidPatchValues(t *testing.T) {
	cases := map[string]PatchParams{
		"bad email": {
			ID:     "0198f8f0-0000-7000-8000-000000000999",
			User:   CreateParams{Email: "not-an-email"},
			Fields: []string{"email"},
		},
		"empty username": {
			ID:     "0198f8f0-0000-7000-8000-000000000999",
			User:   CreateParams{Username: ""},
			Fields: []string{"username"},
		},
		"bad status": {
			ID:     "0198f8f0-0000-7000-8000-000000000999",
			User:   CreateParams{Status: "weird"},
			Fields: []string{"status"},
		},
	}
	service := NewService(stubUserRepository{}, stubTransactionProvider{})
	for name, params := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := service.Patch(context.Background(), params)
			require.Error(t, err)
			assert.Equal(t, apperror.CodeInvalidInput, requireAppError(t, err).Code())
		})
	}
}
