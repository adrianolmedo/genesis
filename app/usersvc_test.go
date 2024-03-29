package app

import (
	"strings"
	"testing"

	domain "github.com/adrianolmedo/genesis"
)

func TestSignUp(t *testing.T) {
	tt := []struct {
		name           string
		input          *domain.User
		errExpected    bool
		wantErrContain string
	}{
		{
			name: "successful",
			input: &domain.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "example@gmail.com",
				Password:  "1234567",
			},
			errExpected:    false,
			wantErrContain: "",
		},
		{
			name: "empty-field",
			input: &domain.User{
				FirstName: "",
				LastName:  "Doe",
				Email:     "example@gmail.com",
				Password:  "1234567",
			},
			errExpected:    true,
			wantErrContain: "first name, email or password can't be empty",
		},
		{
			name: "bad-email",
			input: &domain.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "examplegmailcom",
				Password:  "1234567",
			},
			errExpected:    true,
			wantErrContain: "email not valid",
		},
	}

	for _, tc := range tt {
		err := signUp(tc.input)
		if (err != nil) != tc.errExpected {
			t.Fatalf("%s: unexpected error value %v", tc.name, err)
		}

		if err != nil && !strings.Contains(err.Error(), tc.wantErrContain) {
			t.Fatalf("want error string %q to contain %q", err.Error(), tc.wantErrContain)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	tt := []struct {
		name        string
		email       string
		errExpected bool
	}{
		{name: "typical-email", email: "aol.ve@aol.com", errExpected: false},
		{name: "not-dot-email", email: "aol.ve@aolcom", errExpected: false},
		{name: "not-@-email", email: "aol.veaolcom", errExpected: true},
	}

	u := new(domain.User)

	for _, tc := range tt {
		u.Email = tc.email

		err := validateEmail(u.Email)
		errReceived := err != nil

		if errReceived != tc.errExpected {
			t.Fatalf("%s: unexpected error value: %v", tc.name, err)
		}
	}
}
