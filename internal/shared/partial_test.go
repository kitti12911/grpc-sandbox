package shared

import (
	"testing"

	userv1 "grpc-sandbox/gen/grpc/user/v1"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestValidateMaskAllowsNestedLeafFields(t *testing.T) {
	err := ValidateMask(&fieldmaskpb.FieldMask{Paths: []string{
		"email",
		"profile.first_name",
		"profile.address.country_code",
	}}, userMessage(), map[string]bool{
		"id":         true,
		"created_at": true,
		"updated_at": true,
	})
	if err != nil {
		t.Fatalf("ValidateMask() error = %v", err)
	}
}

func TestValidateMaskRejectsUnsupportedPaths(t *testing.T) {
	tests := []string{
		"id",
		"created_at",
		"updated_at",
		"profile",
		"profile.address",
		"profile.address.unknown",
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			err := ValidateMask(&fieldmaskpb.FieldMask{Paths: []string{path}}, userMessage(), map[string]bool{
				"id":         true,
				"created_at": true,
				"updated_at": true,
			})
			if err == nil {
				t.Fatal("ValidateMask() error = nil")
			}
		})
	}
}

func TestExtractChangesSupportsNestedFields(t *testing.T) {
	displayName := "Display Name"
	firstName := "First"
	line2 := "Floor 2"
	countryCode := "TH"

	changes := ExtractChanges(&fieldmaskpb.FieldMask{Paths: []string{
		"display_name",
		"status",
		"profile.first_name",
		"profile.address.line2",
		"profile.address.country_code",
	}}, &userv1.User{
		DisplayName: &displayName,
		Status:      userv1.UserStatus_USER_STATUS_ACTIVE,
		Profile: &userv1.UserProfile{
			FirstName: &firstName,
			Address: &userv1.UserAddress{
				Line2:       &line2,
				CountryCode: &countryCode,
			},
		},
	})

	assertChange(t, changes, "display_name", displayName)
	assertChange(t, changes, "status", protoreflect.EnumNumber(userv1.UserStatus_USER_STATUS_ACTIVE))
	assertChange(t, changes, "profile.first_name", firstName)
	assertChange(t, changes, "profile.address.line2", line2)
	assertChange(t, changes, "profile.address.country_code", countryCode)
}

func TestExtractChangesCanClearOptionalLeaf(t *testing.T) {
	changes := ExtractChanges(&fieldmaskpb.FieldMask{Paths: []string{
		"profile.address.line2",
	}}, &userv1.User{
		Profile: &userv1.UserProfile{
			Address: &userv1.UserAddress{},
		},
	})

	if got, ok := changes["profile.address.line2"]; !ok || got != nil {
		t.Fatalf("changes[profile.address.line2] = %#v, %t; want nil, true", got, ok)
	}
}

func TestExtractNestedChangesReturnsRootFields(t *testing.T) {
	changes := map[string]any{
		"email":                        "a@example.com",
		"username":                     "alice",
		"profile.first_name":           "Alice",
		"profile.address.country_code": "TH",
	}

	got := ExtractNestedChanges(changes, map[string]string{
		"email":    "email",
		"username": "username",
	}, RootNestedName)

	assertChange(t, got, "email", "a@example.com")
	assertChange(t, got, "username", "alice")
	if len(got) != 2 {
		t.Fatalf("ExtractNestedChanges() len = %d, want 2", len(got))
	}
}

func TestExtractNestedChangesReturnsDirectNestedFields(t *testing.T) {
	changes := map[string]any{
		"email":                        "a@example.com",
		"profile.first_name":           "Alice",
		"profile.last_name":            "Example",
		"profile.address.country_code": "TH",
	}

	got := ExtractNestedChanges(changes, map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
	}, "profile")

	assertChange(t, got, "first_name", "Alice")
	assertChange(t, got, "last_name", "Example")
	if len(got) != 2 {
		t.Fatalf("ExtractNestedChanges() len = %d, want 2", len(got))
	}
}

func TestExtractNestedChangesReturnsDeepNestedFields(t *testing.T) {
	changes := map[string]any{
		"profile.first_name":           "Alice",
		"profile.address.city":         "Bangkok",
		"profile.address.country_code": "TH",
	}

	got := ExtractNestedChanges(changes, map[string]string{
		"city":         "city",
		"country_code": "country_code",
	}, "profile.address")

	assertChange(t, got, "city", "Bangkok")
	assertChange(t, got, "country_code", "TH")
	if len(got) != 2 {
		t.Fatalf("ExtractNestedChanges() len = %d, want 2", len(got))
	}
}

func TestExtractNestedChangesCanRenameFields(t *testing.T) {
	changes := map[string]any{
		"profile.phone_number": "+66123456789",
	}

	got := ExtractNestedChanges(changes, map[string]string{
		"phone_number": "phone",
	}, "profile")

	assertChange(t, got, "phone", "+66123456789")
}

func TestValidateMaskDoesNotRequireParentMessageValues(t *testing.T) {
	err := ValidateMask(&fieldmaskpb.FieldMask{Paths: []string{
		"profile.address.line2",
	}}, &userv1.User{}, nil)
	if err != nil {
		t.Fatalf("ValidateMask() error = %v", err)
	}
}

func userMessage() *userv1.User {
	return &userv1.User{
		Profile: &userv1.UserProfile{
			Address: &userv1.UserAddress{},
		},
	}
}

func assertChange(t *testing.T, changes map[string]any, key string, want any) {
	t.Helper()

	got, ok := changes[key]
	if !ok {
		t.Fatalf("changes[%s] is missing", key)
	}
	if got != want {
		t.Fatalf("changes[%s] = %#v, want %#v", key, got, want)
	}
}
