package user_test

import (
	"testing"

	"github.com/adrianolmedo/go-restapi/user"
)

func TestCheckEmptyFields(t *testing.T) {
	tt := []struct {
		name        string
		user        user.User
		errExpected bool
	}{
		{
			name:        "empty-struct",
			user:        user.User{},
			errExpected: true,
		},
		{
			name:        "empty-fields",
			user:        user.User{FirstName: "", LastName: "", Email: "", Password: ""},
			errExpected: true,
		},
		{
			name: "filled-fields",
			user: user.User{
				FirstName: "Adrián",
				LastName:  "Olmedo",
				Email:     "aol.ve@aol.com",
				Password:  "1234567@"},
			errExpected: false,
		},
	}

	for _, tc := range tt {
		err := tc.user.CheckEmptyFields()
		errReceived := err != nil

		if errReceived != tc.errExpected {
			t.Fatalf("%s: validateUser: unexpected error status: %v", tc.name, err)
		}
	}
}
