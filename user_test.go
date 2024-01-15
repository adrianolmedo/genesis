package genesis_test

import (
	"regexp"
	"testing"

	domain "github.com/adrianolmedo/genesis"
)

func TestCheckEmptyFields(t *testing.T) {
	tt := []struct {
		name        string
		user        domain.User
		errExpected bool
	}{
		{
			name:        "empty-struct",
			user:        domain.User{},
			errExpected: true,
		},
		{
			name:        "empty-fields",
			user:        domain.User{FirstName: "", LastName: "", Email: "", Password: ""},
			errExpected: true,
		},
		{
			name: "filled-fields",
			user: domain.User{
				FirstName: "Adrián",
				LastName:  "Olmedo",
				Email:     "aol.ve@aol.com",
				Password:  "1234567@"},
			errExpected: false,
		},
	}

	for _, tc := range tt {
		err := tc.user.Validate()
		errReceived := err != nil

		if errReceived != tc.errExpected {
			t.Fatalf("%s: validateUser: unexpected error status: %v", tc.name, err)
		}
	}
}

// TestNextUUID se segura que el campo UUID tenga un valor UUID válido.
func TestNextUUID(t *testing.T) {
	uuid := domain.NextUUID()

	if !isValidUUID(uuid) {
		t.Errorf("NextUUID() generate invalid UUID: %s", uuid)
	} else {
		t.Logf("%s: valid! ", uuid)
	}
}

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
