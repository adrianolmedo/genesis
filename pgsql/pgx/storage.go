package pgx

import (
	"fmt"

	"github.com/adrianolmedo/genesis"
)

// Storage represents all repositories.
type Storage struct {
	User     User
	Product  Product
	Customer Customer
	Invoice  Invoice
}

// NewStorage start postgres database connection, build the Storage and return it
// its pointer.
func NewStorage(cfg genesis.Config) (*Storage, error) {
	db, err := newDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres: %v", err)
	}

	return &Storage{
		User:     User{conn: db},
		Product:  Product{conn: db},
		Customer: Customer{conn: db},
		Invoice: Invoice{
			conn:   db,
			header: InvoiceHeader{conn: db},
			items:  InvoiceItem{conn: db},
		},
	}, nil
}
