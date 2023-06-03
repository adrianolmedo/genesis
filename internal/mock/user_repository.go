package mock

import (
	"errors"
	"time"

	"github.com/adrianolmedo/genesis/internal/domain"
)

type UserRepositoryOk struct{}

func (UserRepositoryOk) Create(*domain.User) error {
	return nil
}

func (UserRepositoryOk) ByID(id int64) (*domain.User, error) {
	if id == 1 {
		return &domain.User{
			ID:        1,
			UUID:      "7f5bb73c-d893-4de1-9dc9-beef376e3fbb",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "example@gmail.com",
			Password:  "1234567a",
			CreatedAt: time.Now(),
			UpdatedAt: time.Time{},
			DeletedAt: time.Time{},
		}, nil
	}

	return &domain.User{}, domain.ErrUserNotFound
}

func (UserRepositoryOk) Update(domain.User) error {
	return nil
}

func (UserRepositoryOk) All() ([]*domain.User, error) {
	users := []*domain.User{
		{
			ID:        1,
			UUID:      "7f5bb73c-d893-4de1-9dc9-beef376e3fbb",
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
			UUID:      "230dd357-e116-4baa-be61-c802ad968f3c",
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
