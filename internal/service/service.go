package service

import (
	"fmt"

	"github.com/adrianolmedo/genesis/internal/storage"
)

type Service struct {
	Storage storage.Storage
	User    *userService
	Login   *loginService
	Store   *storeService
	Billing *billingService
}

func New(s storage.Storage) (*Service, error) {
	r, err := s.ProvideRepository()
	if err != nil {
		return nil, fmt.Errorf("error from storage: %v", err)
	}

	return &Service{
		User:    NewUserService(r.User),
		Login:   NewLoginService(r.Login),
		Store:   NewStoreService(r.Product),
		Billing: NewBillingService(r.Invoice),
	}, nil
}
