package mysql_test

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/adrianolmedo/go-restapi/internal/domain"
	"github.com/adrianolmedo/go-restapi/internal/storage/mysql"
)

func TestCreateUser(t *testing.T) {
	t.Cleanup(func() {
		cleanUsersData(t)
	})

	db := openDB(t)
	defer closeDB(t, db)
	r := mysql.NewUserRepository(db)

	input := &domain.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "example@gmail.com",
		Password:  "1234567a",
	}

	if err := r.Create(input); err != nil {
		t.Fatal(err)
	}

	user, err := r.ByID(input.ID)
	if err != nil {
		t.Fatal(err)
	}

	if user.CreatedAt.IsZero() {
		t.Error("expected created at")
	}

	if !user.UpdatedAt.IsZero() {
		t.Error("unexpected updated at")
	}

	if !user.DeletedAt.IsZero() {
		t.Error("unexpected deleted at")
	}
}

func TestUserByID(t *testing.T) {
	t.Cleanup(func() {
		cleanUsersData(t)
	})

	db := openDB(t)
	defer closeDB(t, db)
	insertUsersData(t, db)

	tt := []struct {
		name           string
		input          int64
		wantEmail      string
		errExpected    bool
		wantErrContain string
	}{
		{
			name:           "ok",
			input:          1,
			wantEmail:      "example@gmail.com",
			errExpected:    false,
			wantErrContain: "",
		},
		{
			name:           "user-not-found",
			input:          3,
			wantEmail:      "qwerty@hotmail.com",
			errExpected:    true,
			wantErrContain: "user not found",
		},
	}

	for _, tc := range tt {
		got, err := mysql.NewUserRepository(db).ByID(tc.input)
		if (err != nil) != tc.errExpected {
			t.Fatalf("%s: ByID(%d): unexpected error status: %v", tc.name, tc.input, err)
		}

		if err != nil && !strings.Contains(err.Error(), tc.wantErrContain) {
			t.Fatalf("want error string %q to contain %q", err.Error(), tc.wantErrContain)
		}

		if !tc.errExpected && tc.wantEmail != got.Email {
			t.Fatalf("%s: ByID(%d): want %s, got %s", tc.name, tc.input, tc.wantEmail, got.Email)
		}
	}
}

func cleanUsersData(t *testing.T) {
	db := openDB(t)
	defer closeDB(t, db)

	err := mysql.NewUserRepository(db).DeleteAll()
	if err != nil {
		t.Fatal(err)
	}
}

func insertUsersData(t *testing.T, db *sql.DB) {
	//db := openDB(t)
	//defer closeDB(t, db)

	if err := mysql.NewUserRepository(db).Create(&domain.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "example@gmail.com",
		Password:  "1234567a",
	}); err != nil {
		t.Fatal(err)
	}

	if err := mysql.NewUserRepository(db).Create(&domain.User{
		FirstName: "Jane",
		LastName:  "Roe",
		Email:     "qwerty@hotmail.com",
		Password:  "1234567b",
	}); err != nil {
		t.Fatal(err)
	}
}
