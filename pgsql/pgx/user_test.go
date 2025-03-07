//go:build integration
// +build integration

package pgx

import (
	"testing"

	domain "github.com/adrianolmedo/genesis"

	"github.com/jackc/pgx/v5"
)

func TestCreateUser(t *testing.T) {
	t.Cleanup(func() {
		cleanUsersData(t)
	})

	conn := openDB(t)
	defer closeDB(t, conn)

	u := User{conn: conn}

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

	if !got.UpdatedAt.IsZero() {
		t.Error("unexpected updated at")
	}

	if !got.DeletedAt.IsZero() {
		t.Error("unexpected deleted at")
	}
}

func cleanUsersData(t *testing.T) {
	conn := openDB(t)
	defer closeDB(t, conn)

	u := User{conn}
	err := u.DeleteAll()
	if err != nil {
		t.Fatal(err)
	}
}

func insertUsersData(t *testing.T, conn *pgx.Conn) {
	//conn := openDB(t)
	//defer closeDB(t, conn)

	u := User{conn: conn}

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
