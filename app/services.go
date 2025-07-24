package app

import storage "github.com/adrianolmedo/genesis/pgsql/sqlc"

// Services represents all services layers.
type Services struct {
	User    userService
	Store   storeService
	Billing billingService
}

// NewServices return a pointer of Services.
func NewServices(s *storage.Storage) *Services {
	return &Services{
		User: userService{repo: s.User},
		Store: storeService{
			repoProduct:  s.Product,
			repoCustomer: s.Customer,
		},
		Billing: billingService{repo: s.Invoice},
	}
}
