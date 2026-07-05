package user

import (
	userv1 "grpc-sandbox/gen/grpc/user/v1"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateParamsFromProto and patchParamsFromProto compose the generated
// createParamsFromProto with the extra inputs (id, FieldMask) that the mapper
// does not model.
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
