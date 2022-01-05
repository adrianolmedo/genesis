package domain

import (
	"testing"
)

func TestCheckEmptyFields(t *testing.T) {
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
		err := tc.user.CheckEmptyFields()
		errReceived := err != nil

		if errReceived != tc.errExpected {
			t.Fatalf("%s: validateUser: unexpected error status: %v", tc.name, err)
		}
	}
}
