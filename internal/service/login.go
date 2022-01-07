package service

import "github.com/adrianolmedo/go-restapi-practice/internal/repository"

type LoginService interface {
	Execute(email, password string) error
}

type loginService struct {
	repo repository.LoginRepository
}

func NewLoginService(repo repository.LoginRepository) LoginService {
	return &loginService{repo}
}

func (ls loginService) Execute(email, password string) error {
	if err := validateEmail(email); err != nil {
		return err
	}

	return ls.repo.UserByLogin(email, password)
}
