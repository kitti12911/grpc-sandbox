package user

import (
	userv1 "grpc-sandbox/gen/grpc/user/v1"
	"grpc-sandbox/internal/database"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoUser(user *database.User) *userv1.User {
	if user == nil {
		return nil
	}

	return &userv1.User{
		Id:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Status:      toProtoUserStatus(user.Status),
		Profile:     toProtoUserProfile(user.Profile),
		CreatedAt:   timestamppb.New(user.CreatedAt),
		UpdatedAt:   timestamppb.New(user.UpdatedAt),
	}
}

func toProtoUserProfile(profile *database.UserProfile) *userv1.UserProfile {
	if profile == nil {
		return nil
	}

	return &userv1.UserProfile{
		FirstName:   profile.FirstName,
		LastName:    profile.LastName,
		PhoneNumber: profile.PhoneNumber,
		Address:     toProtoUserAddress(profile.Address),
	}
}

func toProtoUserAddress(address *database.UserAddress) *userv1.UserAddress {
	if address == nil {
		return nil
	}

	return &userv1.UserAddress{
		Line1:       address.Line1,
		Line2:       address.Line2,
		City:        address.City,
		State:       address.State,
		PostalCode:  address.PostalCode,
		CountryCode: address.CountryCode,
	}
}

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
