package mock

import "github.com/adrianolmedo/go-restapi-practice/internal/storage"

type StorageOk struct{}

func (s StorageOk) ProvideRepo() (*storage.Repository, error) {
	return &storage.Repository{
		User:  &UserRepositoryOk{},
		Login: &LoginRepositoryOk{},
	}, nil
}
