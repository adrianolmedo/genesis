package user

import (
	"testing"
)

func TestValidEmail(t *testing.T) {
	user := new(User)

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
		user.Email = tc.email

		err := user.validateEmail()
		errReceived := err != nil

		if errReceived != tc.errExpected {
			t.Fatalf("%s: validateUser: unexpected error status: %v", tc.name, err)
		}
	}
}

func TestValidateNoEmptyFields(t *testing.T) {
	tt := []struct {
		name        string
		user        *User
		errExpected bool
	}{
		/*{
			name:        "nill-struct",
			user:        nil,
			errExpected: true,
		},*/
		{
			name:        "empty-struct",
			user:        &User{},
			errExpected: true,
		},
		{
			name:        "empty-fields",
			user:        &User{FirstName: "", LastName: "", Email: "", Password: ""},
			errExpected: true,
		},
		{
			name:        "filled-fields",
			user:        &User{FirstName: "Adrián", LastName: "Olmedo", Email: "aol.ve@aol.com", Password: "1234567@"},
			errExpected: false,
		},
	}

	for _, tc := range tt {
		err := tc.user.validateNoEmptyFields()
		errReceived := err != nil

		if errReceived != tc.errExpected {
			t.Fatalf("%s: validateUser: unexpected error status: %v", tc.name, err)
		}
	}
}
