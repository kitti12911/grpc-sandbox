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

func (r stubUserRepository) List(context.Context, ListParams) (*ListResult, error) {
	return nil, nil
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
