package service

import (
	"fmt"

	"github.com/adrianolmedo/go-restapi-practice/internal/storage"
)

type Service struct {
	Storage storage.Storage
	User    UserService
	Login   LoginService
}

func New(s storage.Storage) (*Service, error) {
	r, err := s.ProvideRepository()
	if err != nil {
		return nil, fmt.Errorf("error from storage: %v", err)
	}

	return &Service{
		User:  NewUserService(r.User),
		Login: NewLoginService(r.Login),
	}, nil
}
