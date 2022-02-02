package postgres_test

import (
	"testing"

	"github.com/adrianolmedo/go-restapi/internal/domain"
	"github.com/adrianolmedo/go-restapi/internal/storage/postgres"
)

func TestUserByLogin(t *testing.T) {
	t.Cleanup(func() {
		cleanUsersData(t)
	})

	db := openDB(t)
	defer closeDB(t, db)
	insertUsersData(t, db)

	tt := []struct {
		name        string
		input       domain.UserLoginForm
		errExpected bool
	}{
		{
			name: "successful",
			input: domain.UserLoginForm{
				Email:    "example@gmail.com",
				Password: "1234567a",
			},
			errExpected: false,
		},
		{
			name: "user-not-found",
			input: domain.UserLoginForm{
				Email:    "example@hotmail.com",
				Password: "1234567a",
			},
			errExpected: true,
		},
	}

	for _, tc := range tt {
		err := postgres.NewLoginRepository(db).UserByLogin(tc.input.Email, tc.input.Password)
		if (err != nil) != tc.errExpected {
			t.Errorf("%s: UserByLogin(%s, %s): unexpected error status: %v",
				tc.name, tc.input.Email, tc.input.Password, err)
		}
	}
}
