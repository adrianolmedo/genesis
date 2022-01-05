package service

import "github.com/adrianolmedo/go-restapi-practice/internal/storage"

type Services struct {
	UserService  UserService
	LoginService LoginService
}

func NewServices(r storage.Repositories) *Services {
	return &Services{
		UserService:  NewUserService(r.UserRepository),
		LoginService: NewLoginService(r.LoginRepository),
	}
}
