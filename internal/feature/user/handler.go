package user

import (
	"context"

	userv1 "grpc-sandbox/gen/grpc/user/v1"
	"grpc-sandbox/internal/database"
)

type userService interface {
	GetByID(ctx context.Context, params GetByIDParams) (*database.User, error)
}

type Handler struct {
	userv1.UnimplementedUserServiceServer
	userService userService
}

func NewHandler(userService userService) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	result, err := h.userService.GetByID(ctx, GetByIDParams{
		ID: req.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &userv1.GetUserResponse{
		User: toProtoUser(result),
	}, nil
}
