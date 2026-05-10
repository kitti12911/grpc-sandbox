package worker

import (
	"context"
	"fmt"

	workerv1 "grpc-sandbox/gen/grpc/worker/v1"
)

type workerService interface {
	Submit(ctx context.Context, params SubmitParams) (string, error)
}

type Handler struct {
	workerv1.UnimplementedWorkerServiceServer
	workerService workerService
}

func NewHandler(workerService workerService) *Handler {
	return &Handler{workerService: workerService}
}

func (h *Handler) SubmitJob(
	ctx context.Context,
	req *workerv1.SubmitJobRequest,
) (*workerv1.SubmitJobResponse, error) {
	params, err := submitParamsFromProto(req.GetJob())
	if err != nil {
		return nil, err
	}

	id, err := h.workerService.Submit(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("submit worker job: %w", err)
	}

	return &workerv1.SubmitJobResponse{Id: id}, nil
}
