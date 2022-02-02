package service

import (
	"reflect"
	"strings"
	"testing"

	"github.com/adrianolmedo/go-restapi/internal/domain"
	"github.com/adrianolmedo/go-restapi/internal/mock"
	"github.com/adrianolmedo/go-restapi/internal/storage"
)

func TestSignUp(t *testing.T) {
	tt := []struct {
		name           string
		input          *domain.User
		mock           storage.UserRepository
		errExpected    bool
		wantErrContain string
	}{
		{
			name: "successful",
			input: &domain.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "example@gmail.com",
				Password:  "1234567",
			},
			mock:           mock.UserRepositoryOk{},
			errExpected:    false,
			wantErrContain: "",
		},
		{
			name: "empty-field",
			input: &domain.User{
				FirstName: "",
				LastName:  "Doe",
				Email:     "example@gmail.com",
				Password:  "1234567",
			},
			mock:           mock.UserRepositoryOk{},
			errExpected:    true,
			wantErrContain: "first name, email or password can't be empty",
		},
		{
			name: "bad-email",
			input: &domain.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "examplegmailcom",
				Password:  "1234567",
			},
			mock:           mock.UserRepositoryOk{},
			errExpected:    true,
			wantErrContain: "email not valid",
		},
		{
			name: "error-from-repository",
			input: &domain.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "example@gmail.com",
				Password:  "1234567",
			},
			mock:           mock.UserRepositoryError{},
			errExpected:    true,
			wantErrContain: "mock error",
		},
	}

	for _, tc := range tt {
		err := NewUserService(tc.mock).SignUp(tc.input)
		if (err != nil) != tc.errExpected {
			t.Fatalf("%s: unexpected error value %v", tc.name, err)
		}

		if err != nil && !strings.Contains(err.Error(), tc.wantErrContain) {
			t.Fatalf("want error string %q to contain %q", err.Error(), tc.wantErrContain)
		}
	}
}

func TestUserList(t *testing.T) {
	want := domain.UsersList{
		{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "example@gmail.com",
		},
		{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Roe",
			Email:     "qwerty@hotmail.com",
		},
	}

	errExpected := false
	users, err := NewUserService(mock.UserRepositoryOk{}).List()
	if (err != nil) != errExpected {
		t.Fatalf("unexpected error value %v", err)
	}

	got := make(domain.UsersList, 0, len(users))

	assemble := func(u *domain.User) domain.UserProfileDTO {
		return domain.UserProfileDTO{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
		}
	}

	for _, v := range users {
		got = append(got, assemble(v))
	}

	// Only for testing purposes, you may want to use reflect.DeepEqual.
	// It compares two elements of any type recursively.
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}
