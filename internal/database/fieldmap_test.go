package database

import (
	fieldmap "grpc-sandbox/gen/database"
	"reflect"
	"testing"
)

func TestGeneratedUserFieldMaps(t *testing.T) {
	tests := []struct {
		name string
		got  map[string]string
		want map[string]string
	}{
		{
			name: "root",
			got:  fieldmap.UserRootFields,
			want: map[string]string{
				"email":        "email",
				"username":     "username",
				"display_name": "display_name",
				"status":       "status",
			},
		},
		{
			name: "profile",
			got:  fieldmap.UserProfileFields,
			want: map[string]string{
				"first_name":   "first_name",
				"last_name":    "last_name",
				"phone_number": "phone_number",
			},
		},
		{
			name: "address",
			got:  fieldmap.UserAddressFields,
			want: map[string]string{
				"line1":        "line1",
				"line2":        "line2",
				"city":         "city",
				"state":        "state",
				"postal_code":  "postal_code",
				"country_code": "country_code",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Fatalf("generated map = %#v, want %#v", tt.got, tt.want)
			}
		})
	}
}
