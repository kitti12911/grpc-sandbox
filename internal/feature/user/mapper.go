package user

import (
	userv1 "grpc-sandbox/gen/grpc/user/v1"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// toProtoUserStatus and userStatusFromProto bridge the database `status`
// column (string) and the proto UserStatus enum. The generated mappers wire
// them in through `converters:` in protomapgen.yaml.
func toProtoUserStatus(status string) userv1.UserStatus {
	switch status {
	case "active":
		return userv1.UserStatus_USER_STATUS_ACTIVE
	case "disabled":
		return userv1.UserStatus_USER_STATUS_DISABLED
	case "pending":
		return userv1.UserStatus_USER_STATUS_PENDING
	default:
		return userv1.UserStatus_USER_STATUS_UNSPECIFIED
	}
}

func userStatusFromProto(status userv1.UserStatus) string {
	switch status {
	case userv1.UserStatus_USER_STATUS_ACTIVE:
		return "active"
	case userv1.UserStatus_USER_STATUS_DISABLED:
		return "disabled"
	case userv1.UserStatus_USER_STATUS_PENDING:
		return "pending"
	default:
		return ""
	}
}

// updateParamsFromProto and patchParamsFromProto compose the generated
// createParamsFromProto with the extra inputs (id, FieldMask) that
// mapgen proto does not currently model.
func updateParamsFromProto(id string, user *userv1.User) UpdateParams {
	params := createParamsFromProto(user)
	return UpdateParams{
		ID:          id,
		Email:       params.Email,
		Username:    params.Username,
		DisplayName: params.DisplayName,
		Status:      params.Status,
		Profile:     params.Profile,
	}
}

func patchParamsFromProto(id string, user *userv1.User, mask *fieldmaskpb.FieldMask) PatchParams {
	out := PatchParams{
		ID:   id,
		User: createParamsFromProto(user),
	}
	if mask != nil {
		out.Fields = mask.GetPaths()
	}
	return out
}
