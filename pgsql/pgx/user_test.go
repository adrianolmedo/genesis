//go:build integration
// +build integration

package pgx

import (
	"context"
	"testing"

	domain "github.com/adrianolmedo/genesis"

	"github.com/jackc/pgx/v5"
)

func TestCreateUser(t *testing.T) {
	t.Cleanup(func() {
		cleanUsersData(t)
	})

	ctx := testCtx(t)
	conn := openDB(ctx, t)
	defer closeDB(ctx, t, conn)

	u := User{conn: conn}

	input := &domain.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "example@gmail.com",
		Password:  "1234567a",
	}

	if err := u.Create(ctx, input); err != nil {
		t.Fatal(err)
	}

	got, err := u.ByID(ctx, 1)
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

	ctx := testCtx(t)
	conn := openDB(ctx, t)
	defer closeDB(ctx, t, conn)
	insertUsersData(ctx, t, conn)

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

	u := User{conn: conn}

	for _, tc := range tt {
		err := u.ByLogin(ctx, tc.inputEmail, tc.inputPass)
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

	ctx := testCtx(t)
	conn := openDB(ctx, t)
	defer closeDB(ctx, t, conn)
	insertUsersData(ctx, t, conn)

	input := domain.User{
		ID:        1,
		FirstName: "Adrián",
		LastName:  "Olmedo",
		Email:     "example@gmail.com",
		Password:  "1234567a",
	}

	u := User{conn: conn}

	if err := u.Update(ctx, input); err != nil {
		t.Fatal(err)
	}

	got, err := u.ByID(ctx, input.ID)
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
	ctx := testCtx(t)
	conn := openDB(ctx, t)
	defer closeDB(ctx, t, conn)

	u := User{conn}
	err := u.DeleteAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func insertUsersData(ctx context.Context, t *testing.T, conn *pgx.Conn) {
	t.Helper()

	u := User{conn: conn}

	if err := u.Create(ctx, &domain.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "example@gmail.com",
		Password:  "1234567a",
	}); err != nil {
		t.Fatal(err)
	}

	if err := u.Create(ctx, &domain.User{
		FirstName: "Jane",
		LastName:  "Roe",
		Email:     "qwerty@hotmail.com",
		Password:  "1234567b",
	}); err != nil {
		t.Fatal(err)
	}
}
