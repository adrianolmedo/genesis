package user

import (
	"regexp"
	"testing"

	domain "github.com/adrianolmedo/genesis"
)

func TestUser(t *testing.T) {
	tt := []struct {
		name        string
		model       domain.User
		errExpected bool
	}{
		{
			name:        "empty-struct",
			model:       domain.User{},
			errExpected: true,
		},
		{
			name:        "empty-fields",
			model:       domain.User{FirstName: "", LastName: "", Email: "", Password: ""},
			errExpected: true,
		},
		{
			name: "filled-fields",
			model: domain.User{
				FirstName: "Adrián",
				LastName:  "Olmedo",
				Email:     "aol.ve@aol.com",
				Password:  "1234567@"},
			errExpected: false,
		},
	}

	for _, tc := range tt {
		err := tc.model.Validate()
		errReceived := err != nil

		if errReceived != tc.errExpected {
			t.Fatalf("%s: unexpected error value: %v", tc.name, err)
		}
	}
}

// TestNextUUID se segura que el campo UUID tenga un valor UUID válido.
func TestNextUUID(t *testing.T) {
	uuid := domain.NextUUID()

	if !isValidUUID(uuid) {
		t.Errorf("NextUUID() generate invalid UUID: %s", uuid)
	} else {
		t.Logf("%s: UUID valid! ", uuid)
	}
}

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
