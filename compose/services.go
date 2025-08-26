package compose

import (
	"github.com/adrianolmedo/genesis/billing"
	storage "github.com/adrianolmedo/genesis/pgsql/sqlc"
	"github.com/adrianolmedo/genesis/store"
	"github.com/adrianolmedo/genesis/user"
)

// Services holds services and their dependencies.
type Services struct {
	User    *user.Service
	Store   *store.Service
	Billing *billing.Service
}

// NewServices returns a new Services instance with initialized services.
func NewServices(s *storage.Storage) *Services {
	return &Services{
		User:    user.NewService(s.User),
		Store:   store.NewService(s.Product, s.Customer),
		Billing: billing.NewService(s.Invoice),
	}
}
