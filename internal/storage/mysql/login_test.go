//go:build integration
// +build integration

package mysql_test

import (
	"strings"
	"testing"

	"github.com/adrianolmedo/genesis/internal/domain"
	"github.com/adrianolmedo/genesis/internal/storage/mysql"
)

func TestUserByLogin(t *testing.T) {
	t.Cleanup(func() {
		cleanUsersData(t)
	})

	db := openDB(t)
	defer closeDB(t, db)
	insertUsersData(t, db)

	tt := []struct {
		name           string
		input          domain.UserLoginForm
		errExpected    bool
		wantErrContain string
	}{
		{
			name: "successful",
			input: domain.UserLoginForm{
				Email:    "example@gmail.com",
				Password: "1234567a",
			},
			errExpected:    false,
			wantErrContain: "",
		},
		{
			name: "user-not-found",
			input: domain.UserLoginForm{
				Email:    "example@hotmail.com",
				Password: "1234567a",
			},
			errExpected:    true,
			wantErrContain: "user not found",
		},
	}

	for _, tc := range tt {
		err := mysql.NewLoginRepository(db).UserByLogin(tc.input.Email, tc.input.Password)
		if (err != nil) != tc.errExpected {
			t.Errorf("%s: UserByLogin(%s, %s): unexpected error status: %v",
				tc.name, tc.input.Email, tc.input.Password, err)
		}

		if err != nil && !strings.Contains(err.Error(), tc.wantErrContain) {
			t.Fatalf("want error string %q to contain %q", err.Error(), tc.wantErrContain)
		}
	}
}
