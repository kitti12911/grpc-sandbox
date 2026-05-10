package worker

import (
	"context"
	"log/slog"

	"github.com/kitti12911/lib-util/v3/apperror"
	"github.com/kitti12911/lib-util/v3/validator"
)

type publisher interface {
	Publish(ctx context.Context, topic string, payload any, metadata ...map[string]string) error
}

type Service struct {
	publisher publisher
	topic     string
	validator *validator.Validator
}

func NewService(publisher publisher, topic string) *Service {
	return &Service{
		publisher: publisher,
		topic:     topic,
		validator: validator.New(),
	}
}

func (s *Service) Submit(ctx context.Context, params SubmitParams) (string, error) {
	if err := s.validator.Validate(params); err != nil {
		slog.WarnContext(ctx, "invalid submit worker job request", "error", err)
		return "", apperror.InvalidInput("invalid submit worker job request", err)
	}

	job := Job(params)
	if err := s.publisher.Publish(ctx, s.topic, job); err != nil {
		slog.ErrorContext(ctx, "failed to publish worker job", "error", err, "job_id", params.ID, "job_type", params.Type)
		return "", apperror.Internal("failed to submit worker job", err)
	}

	return params.ID, nil
}
