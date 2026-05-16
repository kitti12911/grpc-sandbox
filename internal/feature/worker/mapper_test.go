package worker

import (
	"encoding/json"
	"testing"

	workerv1 "grpc-sandbox/gen/grpc/worker/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestSubmitParamsFromProtoNilJob(t *testing.T) {
	t.Parallel()
	got, err := submitParamsFromProto(nil)
	require.NoError(t, err)
	assert.Equal(t, SubmitParams{}, got)
}

func TestSubmitParamsFromProtoCopiesFields(t *testing.T) {
	t.Parallel()
	payload, err := structpb.NewStruct(map[string]any{"message": "hello"})
	require.NoError(t, err)

	got, err := submitParamsFromProto(&workerv1.WorkerJob{
		Id:      "job-1",
		Type:    "test",
		Payload: payload,
	})
	require.NoError(t, err)
	assert.Equal(t, "job-1", got.ID)
	assert.Equal(t, "test", got.Type)
	assert.JSONEq(t, `{"message":"hello"}`, string(got.Payload))
}

func TestPayloadFromProtoNilPayload(t *testing.T) {
	t.Parallel()
	got, err := payloadFromProto(&workerv1.WorkerJob{Id: "job-1"})
	require.NoError(t, err)
	assert.Equal(t, json.RawMessage(nil), got)
}

func TestPayloadFromProtoMarshalsStruct(t *testing.T) {
	t.Parallel()
	payload, err := structpb.NewStruct(map[string]any{"x": float64(1)})
	require.NoError(t, err)

	got, err := payloadFromProto(&workerv1.WorkerJob{Payload: payload})
	require.NoError(t, err)
	assert.JSONEq(t, `{"x":1}`, string(got))
}
