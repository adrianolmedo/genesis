package mock

import (
	"github.com/adrianolmedo/genesis/internal/domain"
)

type LoginRepositoryOk struct{}

func (LoginRepositoryOk) UserByLogin(email, password string) error {
	return nil
}

type LoginRepositoryError struct{}

func (LoginRepositoryError) UserByLogin(email, password string) error {
	return domain.ErrUserNotFound
}
