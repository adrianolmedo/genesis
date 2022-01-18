package mock

import (
	"errors"
	"time"

	"github.com/adrianolmedo/go-restapi-practice/internal/domain"
)

type UserRepositoryOk struct{}

func (UserRepositoryOk) Create(*domain.User) error {
	return nil
}

func (UserRepositoryOk) ByID(int64) (*domain.User, error) {
	return &domain.User{
		ID:        1,
		UUID:      "a25sH4gd6s2",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "example@gmail.com",
		Password:  "1234567a",
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{},
		DeletedAt: time.Time{},
	}, nil
}

func (UserRepositoryOk) Update(domain.User) error {
	return nil
}

func (UserRepositoryOk) All() ([]*domain.User, error) {
	users := []*domain.User{
		{
			ID:        1,
			UUID:      "a25sH4gd6s2",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "example@gmail.com",
			Password:  "1234567a",
			CreatedAt: time.Now(),
			UpdatedAt: time.Time{},
			DeletedAt: time.Time{},
		},
		{
			ID:        2,
			UUID:      "a25sH4gd6s1",
			FirstName: "Jane",
			LastName:  "Roe",
			Email:     "qwerty@hotmail.com",
			Password:  "1234567b",
			CreatedAt: time.Now(),
			UpdatedAt: time.Time{},
			DeletedAt: time.Time{},
		},
	}
	return users, nil
}

func (UserRepositoryOk) Delete(int64) error {
	return nil
}

// ---

type UserRepositoryError struct{}

func (UserRepositoryError) Create(*domain.User) error {
	return errors.New("mock error")
}

func (UserRepositoryError) ByID(int64) (*domain.User, error) {
	return &domain.User{}, errors.New("mock error")
}

func (UserRepositoryError) Update(domain.User) error {
	return errors.New("mock error")
}

func (UserRepositoryError) All() ([]*domain.User, error) {
	users := make([]*domain.User, 0)
	return users, errors.New("mock error")
}

func (UserRepositoryError) Delete(int64) error {
	return errors.New("mock error")
}
