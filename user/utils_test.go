package user

import "testing"

func TestValidateEmail(t *testing.T) {
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

		err := validateEmail(user.Email)
		errReceived := err != nil
		if errReceived != tc.errExpected {
			t.Fatalf("%s: unexpected error value: %v", tc.name, err)
		}
	}
}
