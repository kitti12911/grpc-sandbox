package mockuser

import (
	"context"
	"testing"

	userv1 "grpc-sandbox/gen/grpc/user/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerGetMockUser(t *testing.T) {
	handler := NewHandler()

	resp, err := handler.GetMockUser(context.Background(), &userv1.GetMockUserRequest{
		Id: "mock-123",
	})

	require.NoError(t, err)
	require.NotNil(t, resp.GetUser())
	assert.Equal(t, "mock-123", resp.GetUser().GetId())
	assert.Equal(t, "Mock User", resp.GetUser().GetName())
}

func TestHandlerGetMockUserDefaultsID(t *testing.T) {
	handler := NewHandler()

	resp, err := handler.GetMockUser(context.Background(), &userv1.GetMockUserRequest{})

	require.NoError(t, err)
	require.NotNil(t, resp.GetUser())
	assert.Equal(t, "mock-user", resp.GetUser().GetId())
	assert.Equal(t, "Mock User", resp.GetUser().GetName())
}
