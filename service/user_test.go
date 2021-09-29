package service

import (
	"testing"

	"go-rest-practice/model"
)

func TestValidEmail(t *testing.T) {
	tt := []struct {
		name        string
		email       string
		errExpected bool
	}{
		{name: "typical-email", email: "aol.ve@aol.com", errExpected: false},
		{name: "not-dot-email", email: "aol.ve@aolcom", errExpected: false},
		{name: "not-@-email", email: "aol.veaolcom", errExpected: true},
	}

	for _, tc := range tt {
		err := validateEmail(tc.email)
		errReceived := err != nil

		if errReceived != tc.errExpected {
			t.Fatalf("%s: validateUser: unexpected error status: %v", tc.name, err)
		}
	}
}

func TestValidateUser(t *testing.T) {
	tt := []struct {
		name        string
		user        *model.User
		errExpected bool
	}{
		{
			name:        "nill-struct",
			user:        nil,
			errExpected: true,
		},
		{
			name:        "empty-struct",
			user:        &model.User{},
			errExpected: true,
		},
		{
			name:        "empty-fields",
			user:        &model.User{FirstName: "", LastName: "", Email: "", Password: ""},
			errExpected: true,
		},
		{
			name:        "filled-fields",
			user:        &model.User{FirstName: "Adrián", LastName: "Olmedo", Email: "aol.ve@aol.com", Password: "1234567@"},
			errExpected: false,
		},
	}

	for _, tc := range tt {
		err := validateUser(tc.user)
		errReceived := err != nil

		if errReceived != tc.errExpected {
			t.Fatalf("%s: validateUser: unexpected error status: %v", tc.name, err)
		}
	}
}
