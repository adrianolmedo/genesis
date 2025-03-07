package app

import strg "github.com/adrianolmedo/genesis/pgsql/pq"

// Services represents all services layers.
type Services struct {
	User    userService
	Store   storeService
	Billing billingService
}

// NewServices return a pointer of Services.
func NewServices(strg *strg.Storage) *Services {
	return &Services{
		User: userService{repo: strg.User},
		Store: storeService{
			repoProduct:  strg.Product,
			repoCustomer: strg.Customer,
		},
		Billing: billingService{repo: strg.Invoice},
	}
}
