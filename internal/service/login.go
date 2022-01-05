package service

import "github.com/adrianolmedo/go-restapi-practice/internal/storage"

type LoginService interface {
	Execute(email, password string) error
}

type loginService struct {
	repository storage.LoginRepository
}

func NewLoginService(r storage.LoginRepository) LoginService {
	return &loginService{r}
}

func (s loginService) Execute(email, password string) error {
	if err := validateEmail(email); err != nil {
		return err
	}

	return s.repository.UserByLogin(email, password)
}
