//go:build integration
// +build integration

package storage_test

import (
	"database/sql"
	"testing"

	"github.com/adrianolmedo/go-restapi/domain"
	"github.com/adrianolmedo/go-restapi/storage"
)

func TestCreateUser(t *testing.T) {
	t.Cleanup(func() {
		cleanUsersData(t)
	})

	db := openDB(t)
	defer closeDB(t, db)

	r := storage.NewUserRepository(db)

	input := &domain.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "example@gmail.com",
		Password:  "1234567a",
	}

	if err := r.Create(input); err != nil {
		t.Fatal(err)
	}

	user, err := r.ByID(1)
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

func cleanUsersData(t *testing.T) {
	db := openDB(t)
	defer closeDB(t, db)

	err := storage.NewUserRepository(db).DeleteAll()
	if err != nil {
		t.Fatal(err)
	}
}

func insertUsersData(t *testing.T, db *sql.DB) {
	//db := openDB(t)
	//defer closeDB(t, db)

	if err := storage.NewUserRepository(db).Create(&domain.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "example@gmail.com",
		Password:  "1234567a",
	}); err != nil {
		t.Fatal(err)
	}

	if err := storage.NewUserRepository(db).Create(&domain.User{
		FirstName: "Jane",
		LastName:  "Roe",
		Email:     "qwerty@hotmail.com",
		Password:  "1234567b",
	}); err != nil {
		t.Fatal(err)
	}
}