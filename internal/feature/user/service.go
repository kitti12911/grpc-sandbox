package user

import (
	"context"
	"log/slog"

	"grpc-sandbox/internal/apperror"
	"grpc-sandbox/internal/database"

	"github.com/go-playground/validator/v10"
)

type userRepository interface {
	GetByID(ctx context.Context, id string) (*database.User, error)
}

type Service struct {
	userRepository userRepository
	validator      *validator.Validate
}

func NewService(userRepository userRepository) *Service {
	return &Service{
		userRepository: userRepository,
		validator:      validator.New(),
	}
}

func (s *Service) GetByID(ctx context.Context, params GetByIDParams) (*database.User, error) {
	if err := s.validator.Struct(params); err != nil {
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
