package worker

import (
	"encoding/json"
	"fmt"

	workerv1 "grpc-sandbox/gen/grpc/worker/v1"
)

func submitParamsFromProto(job *workerv1.WorkerJob) (SubmitParams, error) {
	if job == nil {
		return SubmitParams{}, nil
	}

	payload, err := payloadFromProto(job)
	if err != nil {
		return SubmitParams{}, err
	}

	return SubmitParams{
		ID:      job.GetId(),
		Type:    job.GetType(),
		Payload: payload,
	}, nil
}

func payloadFromProto(job *workerv1.WorkerJob) (json.RawMessage, error) {
	if job.GetPayload() == nil {
		return nil, nil
	}

	payload, err := job.GetPayload().MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("marshal worker job payload: %w", err)
	}
	return json.RawMessage(payload), nil
}
