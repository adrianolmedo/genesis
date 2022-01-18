package service

import (
	"fmt"

	"github.com/adrianolmedo/go-restapi-practice/internal/storage"
)

type Service struct {
	Repo  storage.RepoProvider
	User  UserService
	Login LoginService
}

func New(p storage.RepoProvider) (*Service, error) {
	r, err := p.ProvideRepo()
	if err != nil {
		return nil, fmt.Errorf("error from storage: %v", err)
	}

	return &Service{
		User:  NewUserService(r.User),
		Login: NewLoginService(r.Login),
	}, nil
}
