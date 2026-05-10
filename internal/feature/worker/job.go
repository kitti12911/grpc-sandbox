package worker

import "encoding/json"

type Job struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type SubmitParams struct {
	ID      string          `validate:"required"`
	Type    string          `validate:"required"`
	Payload json.RawMessage `validate:"-"`
}
