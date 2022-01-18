package service

import "github.com/adrianolmedo/go-restapi-practice/internal/storage"

type LoginService interface {
	Execute(email, password string) error
}

type loginService struct {
	repo storage.LoginRepository
}

func NewLoginService(repo storage.LoginRepository) LoginService {
	return &loginService{repo}
}

func (ls loginService) Execute(email, password string) error {
	if err := validateEmail(email); err != nil {
		return err
	}

	return ls.repo.UserByLogin(email, password)
}
