//go:build integration
// +build integration

package pq

import (
	"database/sql"
	"testing"

	domain "github.com/adrianolmedo/genesis"
)

func TestCreateUser(t *testing.T) {
	t.Cleanup(func() {
		cleanUsersData(t)
	})

	db := openDB(t)
	defer closeDB(t, db)

	u := User{db: db}

	input := &domain.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "example@gmail.com",
		Password:  "1234567a",
	}

	if err := u.Create(input); err != nil {
		t.Fatal(err)
	}

	got, err := u.ByID(1)
	if err != nil {
		t.Fatal(err)
	}

	if got.CreatedAt.IsZero() {
		t.Error("expected created at")
	}

	if got.UpdatedAt != nil {
		t.Error("unexpected updated at")
	}

	if got.DeletedAt != nil {
		t.Error("unexpected deleted at")
	}
}

func TestUserByLogin(t *testing.T) {
	t.Cleanup(func() {
		cleanUsersData(t)
	})

	db := openDB(t)
	defer closeDB(t, db)
	insertUsersData(t, db)

	tt := []struct {
		name        string
		inputEmail  string
		inputPass   string
		errExpected bool
	}{
		{
			name:        "successful",
			inputEmail:  "example@gmail.com",
			inputPass:   "1234567a",
			errExpected: false,
		},
		{
			name:        "user-not-found",
			inputEmail:  "example@hotmail.com",
			inputPass:   "1234567a",
			errExpected: true,
		},
	}

	u := User{db: db}

	for _, tc := range tt {
		err := u.ByLogin(tc.inputEmail, tc.inputPass)
		if (err != nil) != tc.errExpected {
			t.Errorf("%s: ByLogin(%s, %s): unexpected error status: %v",
				tc.name, tc.inputEmail, tc.inputPass, err)
		}
	}
}

func TestUpdateUser(t *testing.T) {
	t.Cleanup(func() {
		cleanUsersData(t)
	})

	db := openDB(t)
	defer closeDB(t, db)
	insertUsersData(t, db)

	input := domain.User{
		ID:        1,
		FirstName: "Adrián",
		LastName:  "Olmedo",
		Email:     "example@gmail.com",
		Password:  "1234567a",
	}

	u := User{db: db}

	if err := u.Update(input); err != nil {
		t.Fatal(err)
	}

	got, err := u.ByID(input.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.FirstName != input.FirstName {
		t.Errorf("FirstName: want %s, got %s", input.FirstName, got.FirstName)
	}

	if got.LastName != input.LastName {
		t.Errorf("LastName: want %s, got %s", input.LastName, got.LastName)
	}

	if got.CreatedAt.IsZero() {
		t.Error("expected created at")
	}

	if got.UpdatedAt.IsZero() {
		t.Error("expected updated at")
	}

	if got.DeletedAt != nil {
		t.Error("unexpected deleted at")
	}
}

func cleanUsersData(t *testing.T) {
	db := openDB(t)
	defer closeDB(t, db)

	u := User{db}
	err := u.DeleteAll()
	if err != nil {
		t.Fatal(err)
	}
}

func insertUsersData(t *testing.T, db *sql.DB) {
	//db := openDB(t)
	//defer closeDB(t, db)

	u := User{db: db}

	if err := u.Create(&domain.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "example@gmail.com",
		Password:  "1234567a",
	}); err != nil {
		t.Fatal(err)
	}

	if err := u.Create(&domain.User{
		FirstName: "Jane",
		LastName:  "Roe",
		Email:     "qwerty@hotmail.com",
		Password:  "1234567b",
	}); err != nil {
		t.Fatal(err)
	}
}
