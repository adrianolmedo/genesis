package mysql_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/adrianolmedo/go-restapi-practice/internal/domain"
	"github.com/adrianolmedo/go-restapi-practice/internal/storage/mysql"
)

func TestCreateUser(t *testing.T) {
	t.Cleanup(func() {
		cleanUsersData(t)
	})

	db := openDB(t)
	defer closeDB(t, db)
	r := mysql.NewUserRepository(db)

	input := &domain.User{
		UUID:      domain.NextUserID(),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "example@gmail.com",
		Password:  "1234567a",
		CreatedAt: time.Now(),
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
		name        string
		input       int64
		wantEmail   string
		errExpected bool
	}{
		{
			name:        "ok",
			input:       1,
			wantEmail:   "example@gmail.com",
			errExpected: false,
		},
		{
			name:        "not-found",
			input:       3,
			wantEmail:   "qwerty@hotmail.com",
			errExpected: true,
		},
	}

	for _, tc := range tt {
		got, err := mysql.NewUserRepository(db).ByID(tc.input)
		errReceived := err != nil

		if errReceived != tc.errExpected {
			t.Fatalf("%s: ByID(%d): unexpected error status: %v", tc.name, tc.input, err)
		}

		if !tc.errExpected && tc.wantEmail != got.Email {
			t.Errorf("%s: ByID(%d): want %s, got %s", tc.name, tc.input, tc.wantEmail, got.Email)
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
		UUID:      domain.NextUserID(),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "example@gmail.com",
		Password:  "1234567a",
		CreatedAt: time.Now(),
	}); err != nil {
		t.Fatal(err)
	}

	if err := mysql.NewUserRepository(db).Create(&domain.User{
		UUID:      domain.NextUserID(),
		FirstName: "Jane",
		LastName:  "Roe",
		Email:     "qwerty@hotmail.com",
		Password:  "1234567b",
		CreatedAt: time.Now(),
	}); err != nil {
		t.Fatal(err)
	}
}
