package user

import (
	"testing"

	userv1 "grpc-sandbox/gen/grpc/user/v1"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// Only the hand-written adapters in mapper.go are tested here; the generated
// mappers in mapper_generated.go are the generator's responsibility (covered by
// lib-orm/internal/mappergen) and are excluded from coverage in go-test.sh.

func TestUpdateParamsFromProto(t *testing.T) {
	t.Parallel()

	got := updateParamsFromProto("u1", nil)
	assert.Equal(t, UpdateParams{ID: "u1"}, got)

	user := &userv1.User{Email: "kit@example.com", Username: "kit", Status: userv1.UserStatus_USER_STATUS_PENDING}
	got = updateParamsFromProto("u1", user)
	assert.Equal(t, "u1", got.ID)
	assert.Equal(t, "kit@example.com", got.Email)
	assert.Equal(t, "kit", got.Username)
	assert.Equal(t, "pending", got.Status)
}

func TestPatchParamsFromProto(t *testing.T) {
	t.Parallel()

	user := &userv1.User{Email: "kit@example.com"}

	// Nil mask: fields slice is unset, params still carry user data.
	noMask := patchParamsFromProto("u1", user, nil)
	assert.Equal(t, "u1", noMask.ID)
	assert.Equal(t, "kit@example.com", noMask.User.Email)
	assert.Nil(t, noMask.Fields)

	withMask := patchParamsFromProto("u1", user, &fieldmaskpb.FieldMask{Paths: []string{"email", "username"}})
	assert.Equal(t, []string{"email", "username"}, withMask.Fields)
	assert.Equal(t, "kit@example.com", withMask.User.Email)
}
