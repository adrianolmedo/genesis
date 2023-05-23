package gorestapi

import "github.com/adrianolmedo/go-restapi/postgres"

type Services struct {
	User    userService
	Store   storeService
	Billing billingService
}

func NewServices(strg *postgres.Storage) *Services {
	return &Services{
		User:    userService{repo: strg.User},
		Store:   storeService{repo: strg.Product},
		Billing: billingService{repo: strg.Invoice},
	}
}
