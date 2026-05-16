package user

import (
	"testing"
	"time"

	userv1 "grpc-sandbox/gen/grpc/user/v1"
	"grpc-sandbox/internal/database"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestToProtoUserStatus(t *testing.T) {
	t.Parallel()

	tests := map[string]userv1.UserStatus{
		"active":   userv1.UserStatus_USER_STATUS_ACTIVE,
		"disabled": userv1.UserStatus_USER_STATUS_DISABLED,
		"pending":  userv1.UserStatus_USER_STATUS_PENDING,
		"":         userv1.UserStatus_USER_STATUS_UNSPECIFIED,
		"unknown":  userv1.UserStatus_USER_STATUS_UNSPECIFIED,
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			assert.Equal(t, want, toProtoUserStatus(input))
		})
	}
}

func TestUserStatusFromProto(t *testing.T) {
	t.Parallel()

	tests := map[userv1.UserStatus]string{
		userv1.UserStatus_USER_STATUS_ACTIVE:      "active",
		userv1.UserStatus_USER_STATUS_DISABLED:    "disabled",
		userv1.UserStatus_USER_STATUS_PENDING:     "pending",
		userv1.UserStatus_USER_STATUS_UNSPECIFIED: "",
	}

	for input, want := range tests {
		t.Run(input.String(), func(t *testing.T) {
			assert.Equal(t, want, userStatusFromProto(input))
		})
	}
}

func TestToProtoUserNil(t *testing.T) {
	t.Parallel()
	assert.Nil(t, toProtoUser(nil))
	assert.Nil(t, toProtoUserProfile(nil))
	assert.Nil(t, toProtoUserAddress(nil))
}

func TestToProtoUserCopiesFields(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	user := &database.User{
		ID:          "u1",
		Email:       "kit@example.com",
		Username:    "kit",
		DisplayName: new("Kit"),
		Status:      "active",
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Profile: &database.UserProfile{
			FirstName:   new("Kit"),
			LastName:    new("Last"),
			PhoneNumber: new("0800"),
			Address: &database.UserAddress{
				Line1:       new("1 Road"),
				Line2:       new("Apt 2"),
				City:        new("Bangkok"),
				State:       new("BKK"),
				PostalCode:  new("10110"),
				CountryCode: new("TH"),
			},
		},
	}

	got := toProtoUser(user)

	assert.Equal(t, "u1", got.GetId())
	assert.Equal(t, "kit@example.com", got.GetEmail())
	assert.Equal(t, "kit", got.GetUsername())
	assert.Equal(t, "Kit", got.GetDisplayName())
	assert.Equal(t, userv1.UserStatus_USER_STATUS_ACTIVE, got.GetStatus())
	assert.Equal(t, createdAt, got.GetCreatedAt().AsTime())
	assert.Equal(t, updatedAt, got.GetUpdatedAt().AsTime())

	profile := got.GetProfile()
	assert.Equal(t, "Kit", profile.GetFirstName())
	assert.Equal(t, "Last", profile.GetLastName())
	assert.Equal(t, "0800", profile.GetPhoneNumber())

	address := profile.GetAddress()
	assert.Equal(t, "1 Road", address.GetLine1())
	assert.Equal(t, "Apt 2", address.GetLine2())
	assert.Equal(t, "Bangkok", address.GetCity())
	assert.Equal(t, "BKK", address.GetState())
	assert.Equal(t, "10110", address.GetPostalCode())
	assert.Equal(t, "TH", address.GetCountryCode())
}

func TestCreateParamsFromProtoNil(t *testing.T) {
	t.Parallel()
	assert.Equal(t, CreateParams{}, createParamsFromProto(nil))
	assert.Nil(t, createProfileParamsFromProto(nil))
	assert.Nil(t, createAddressParamsFromProto(nil))
}

func TestCreateParamsFromProtoCopiesNestedFields(t *testing.T) {
	t.Parallel()

	user := &userv1.User{
		Email:       "kit@example.com",
		Username:    "kit",
		DisplayName: new("Kit"),
		Status:      userv1.UserStatus_USER_STATUS_ACTIVE,
		Profile: &userv1.UserProfile{
			FirstName:   new("Kit"),
			LastName:    new("Last"),
			PhoneNumber: new("0800"),
			Address: &userv1.UserAddress{
				Line1:       new("1 Road"),
				City:        new("Bangkok"),
				CountryCode: new("TH"),
			},
		},
	}

	got := createParamsFromProto(user)
	assert.Equal(t, "kit@example.com", got.Email)
	assert.Equal(t, "kit", got.Username)
	assert.Equal(t, "Kit", *got.DisplayName)
	assert.Equal(t, "active", got.Status)

	require := got.Profile
	assert.Equal(t, "Kit", *require.FirstName)
	assert.Equal(t, "Last", *require.LastName)
	assert.Equal(t, "0800", *require.PhoneNumber)

	addr := require.Address
	assert.Equal(t, "1 Road", *addr.Line1)
	assert.Equal(t, "Bangkok", *addr.City)
	assert.Equal(t, "TH", *addr.CountryCode)
}

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
