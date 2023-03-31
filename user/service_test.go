package user

import (
	"strings"
	"testing"
)

func TestSignUpService(t *testing.T) {
	tt := []struct {
		name           string
		input          *User
		errExpected    bool
		wantErrContain string
	}{
		{
			name: "successful",
			input: &User{
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
			input: &User{
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
			input: &User{
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
		err := signUpService(tc.input)
		if (err != nil) != tc.errExpected {
			t.Fatalf("%s: unexpected error value %v", tc.name, err)
		}

		if err != nil && !strings.Contains(err.Error(), tc.wantErrContain) {
			t.Fatalf("want error string %q to contain %q", err.Error(), tc.wantErrContain)
		}
	}
}
