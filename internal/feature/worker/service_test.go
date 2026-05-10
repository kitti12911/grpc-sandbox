package worker

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/kitti12911/lib-util/v3/apperror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubPublisher struct {
	publishFunc func(ctx context.Context, topic string, payload any, metadata ...map[string]string) error
}

func (p stubPublisher) Publish(ctx context.Context, topic string, payload any, metadata ...map[string]string) error {
	if p.publishFunc == nil {
		return nil
	}
	return p.publishFunc(ctx, topic, payload, metadata...)
}

func TestServiceSubmitPublishesJob(t *testing.T) {
	expectedPayload := json.RawMessage(`{"message":"hello"}`)
	called := false
	service := NewService(stubPublisher{
		publishFunc: func(_ context.Context, topic string, payload any, _ ...map[string]string) error {
			called = true
			assert.Equal(t, "worker-sandbox-jobs", topic)
			assert.Equal(t, Job{
				ID:      "job-1",
				Type:    "debug.print",
				Payload: expectedPayload,
			}, payload)
			return nil
		},
	}, "worker-sandbox-jobs")

	id, err := service.Submit(context.Background(), SubmitParams{
		ID:      "job-1",
		Type:    "debug.print",
		Payload: expectedPayload,
	})

	require.NoError(t, err)
	assert.Equal(t, "job-1", id)
	assert.True(t, called)
}

func TestServiceSubmitValidatesRequest(t *testing.T) {
	called := false
	service := NewService(stubPublisher{
		publishFunc: func(context.Context, string, any, ...map[string]string) error {
			called = true
			return nil
		},
	}, "worker-sandbox-jobs")

	_, err := service.Submit(context.Background(), SubmitParams{})

	require.Error(t, err)
	appErr := requireAppError(t, err)
	assert.Equal(t, apperror.CodeInvalidInput, appErr.Code())
	assert.False(t, called)
}

func TestServiceSubmitReturnsInternalErrorWhenPublishFails(t *testing.T) {
	expected := errors.New("publish failed")
	service := NewService(stubPublisher{
		publishFunc: func(context.Context, string, any, ...map[string]string) error {
			return expected
		},
	}, "worker-sandbox-jobs")

	_, err := service.Submit(context.Background(), SubmitParams{
		ID:   "job-1",
		Type: "debug.print",
	})

	require.Error(t, err)
	appErr := requireAppError(t, err)
	assert.Equal(t, apperror.CodeInternal, appErr.Code())
	assert.ErrorIs(t, err, expected)
}

func requireAppError(t *testing.T, err error) *apperror.Error {
	t.Helper()

	var appErr *apperror.Error
	require.True(t, errors.As(err, &appErr))

	return appErr
}
