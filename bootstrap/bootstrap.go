package bootstrap

import (
	"github.com/adrianolmedo/genesis/billing"
	storage "github.com/adrianolmedo/genesis/pgsql/sqlc"
	"github.com/adrianolmedo/genesis/store"
	"github.com/adrianolmedo/genesis/user"
)

// Bootstrap holds services and their dependencies.
type Bootstrap struct {
	User    *user.Service
	Store   *store.Service
	Billing *billing.Service
}

// New returns a new Bootstrap instance with initialized services.
func New(s *storage.Storage) *Bootstrap {
	return &Bootstrap{
		User:    user.NewService(s.User),
		Store:   store.NewService(s.Product, s.Customer),
		Billing: billing.NewService(s.Invoice),
	}
}
