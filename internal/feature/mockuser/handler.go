package mockuser

import (
	"context"

	userv1 "grpc-sandbox/gen/grpc/user/v1"
)

type Handler struct {
	userv1.UnimplementedMockUserServiceServer
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetMockUser(
	_ context.Context,
	req *userv1.GetMockUserRequest,
) (*userv1.GetMockUserResponse, error) {
	id := req.GetId()
	if id == "" {
		id = "mock-user"
	}

	return &userv1.GetMockUserResponse{
		User: &userv1.MockUser{
			Id:   id,
			Name: "Mock User",
		},
	}, nil
}
