package mock

import "github.com/adrianolmedo/go-restapi/internal/storage"

type StorageOk struct{}

func (StorageOk) ProvideRepository() (*storage.Repository, error) {
	return &storage.Repository{
		User:  &UserRepositoryOk{},
		Login: &LoginRepositoryOk{},
	}, nil
}

type StorageError struct{}

func (StorageError) ProvideRepository() (*storage.Repository, error) {
	return &storage.Repository{
		User:  &UserRepositoryError{},
		Login: &LoginRepositoryError{},
	}, nil
}
