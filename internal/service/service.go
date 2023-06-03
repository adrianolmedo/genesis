package service

import (
	"fmt"

	"github.com/adrianolmedo/genesis/internal/storage"
)

type Service struct {
	Storage storage.Storage
	User    UserService
	Login   LoginService
	Store   StoreService
	Billing BillingService
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
