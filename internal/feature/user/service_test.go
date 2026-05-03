package user

import (
	"context"
	"testing"

	"grpc-sandbox/internal/apperror"
	"grpc-sandbox/internal/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubUserRepository struct {
	createUserFunc func(ctx context.Context, params CreateParams) (*database.User, error)
	updateUserFunc func(ctx context.Context, params UpdateParams) (int64, error)
	deleteUserFunc func(ctx context.Context, userID string) (int64, error)
}

func (r stubUserRepository) GetByID(context.Context, string) (*database.User, error) {
	return nil, nil
}

func (r stubUserRepository) CreateUser(ctx context.Context, params CreateParams) (*database.User, error) {
	if r.createUserFunc == nil {
		return &database.User{}, nil
	}
	return r.createUserFunc(ctx, params)
}

func (r stubUserRepository) CreateProfile(
	context.Context,
	string,
	CreateProfileParams,
) (*database.UserProfile, error) {
	return &database.UserProfile{}, nil
}

func (r stubUserRepository) CreateAddress(
	context.Context,
	string,
	CreateAddressParams,
) (*database.UserAddress, error) {
	return &database.UserAddress{}, nil
}

func (r stubUserRepository) GetProfileIDByUserID(context.Context, string) (string, error) {
	return "", nil
}

func (r stubUserRepository) UpdateUser(ctx context.Context, params UpdateParams) (int64, error) {
	if r.updateUserFunc == nil {
		return 0, nil
	}
	return r.updateUserFunc(ctx, params)
}

func (r stubUserRepository) UpdateProfileByUserID(
	context.Context,
	string,
	CreateProfileParams,
) (int64, error) {
	return 0, nil
}

func (r stubUserRepository) UpdateAddressByProfileID(
	context.Context,
	string,
	CreateAddressParams,
) (int64, error) {
	return 0, nil
}

func (r stubUserRepository) List(context.Context, ListParams) (*ListResult, error) {
	return nil, nil
}

func (r stubUserRepository) DeleteAddressByProfileID(context.Context, string) error {
	return nil
}

func (r stubUserRepository) DeleteAddressByUserID(context.Context, string) error {
	return nil
}

func (r stubUserRepository) DeleteProfileByUserID(context.Context, string) error {
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
	appErr, ok := err.(*apperror.Error)
	require.True(t, ok)
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
	appErr, ok := err.(*apperror.Error)
	require.True(t, ok)
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
	appErr, ok := err.(*apperror.Error)
	require.True(t, ok)
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
	appErr, ok := err.(*apperror.Error)
	require.True(t, ok)
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
	appErr, ok := err.(*apperror.Error)
	require.True(t, ok)
	assert.Equal(t, apperror.CodeNotFound, appErr.Code())
}
