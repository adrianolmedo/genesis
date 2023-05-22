package gorestapi

import "github.com/adrianolmedo/go-restapi/postgres"

type Services struct {
	User    UserService
	Store   StoreService
	Billing BillingService
}

func NewServices(strg *postgres.Storage) *Services {
	return &Services{
		User:    UserService{repo: strg.User},
		Store:   StoreService{repo: strg.Product},
		Billing: BillingService{repo: strg.Invoice},
	}
}
