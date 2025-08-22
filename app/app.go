package app

import (
	"github.com/adrianolmedo/genesis/billing"
	storage "github.com/adrianolmedo/genesis/pgsql/sqlc"
	"github.com/adrianolmedo/genesis/store"
	"github.com/adrianolmedo/genesis/user"
)

// App holds the services of the application.
type App struct {
	User    *user.Service
	Store   *store.Service
	Billing *billing.Service
}

// NewApp creates a new App with the given services.
func NewApp(s *storage.Storage) *App {
	return &App{
		User:    user.NewService(s.User),
		Store:   store.NewService(s.Product, s.Customer),
		Billing: billing.NewService(s.Invoice),
	}
}
