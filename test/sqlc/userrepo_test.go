package sqlc

import (
	"testing"

	"github.com/adrianolmedo/genesis/test"
	"github.com/adrianolmedo/genesis/user"
)

// TestCreateUser go test -v -run '^TestCreateUser' -tags integration -args -database-url postgres://user:password@host:port/dbname?sslmode=disable
func TestCreateUser(t *testing.T) {
	t.Cleanup(func() {
		cleanUsersData(t)
	})
	ctx := test.Ctx(t)
	db := openDB(ctx, t)
	defer db.Close()

	input := &user.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "example@gmail.com",
		Password:  "1234567a",
	}
	u := user.NewRepo(db)
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

// cleanUsersData removes all user data from the database.
func cleanUsersData(t *testing.T) {
	ctx := test.Ctx(t)
	db := openDB(ctx, t)
	defer db.Close()

	u := user.NewRepo(db)
	err := u.DeleteAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
