package mock

import "errors"

type LoginRepositoryOk struct{}

func (LoginRepositoryOk) UserByLogin(email, password string) error {
	return nil
}

type LoginRepositoryError struct{}

func (LoginRepositoryError) UserByLogin(email, password string) error {
	return errors.New("mock error")
}
