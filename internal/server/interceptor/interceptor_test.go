package interceptor

import (
	"context"
	"testing"

	"grpc-sandbox/internal/apperror"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrorHandlerAddsTraceIDToAppError(t *testing.T) {
	traceID := trace.TraceID{1, 2, 3}
	ctx := trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  trace.SpanID{4, 5, 6},
	}))

	_, err := ErrorHandler()(ctx, nil, &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/Create",
	}, func(context.Context, any) (any, error) {
		return nil, apperror.InvalidInput("invalid input", nil)
	})

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Equal(t, "invalid input (trace_id="+traceID.String()+")", st.Message())
}

func TestMessageWithTraceID(t *testing.T) {
	assert.Equal(t, "failed", messageWithTraceID("failed", ""))
	assert.Equal(t, "failed (trace_id=abc)", messageWithTraceID("failed", "abc"))
}
