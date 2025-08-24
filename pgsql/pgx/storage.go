package pgx

import (
	"context"
	"fmt"

	"github.com/adrianolmedo/genesis"
)

// Storage represents all repositories.
type Storage struct {
	User     UserRepo
	Product  Product
	Customer Customer
	Invoice  Invoice
}

// NewStorage start postgres database connection, build the Storage and return it
// its pointer.
func NewStorage(ctx context.Context, cfg genesis.Config) (*Storage, error) {
	db, err := newDB(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres: %v", err)
	}

	return &Storage{
		User:     UserRepo{conn: db},
		Product:  Product{conn: db},
		Customer: Customer{conn: db},
		Invoice: Invoice{
			conn:   db,
			header: InvoiceHeader{conn: db},
			items:  InvoiceItem{conn: db},
		},
	}, nil
}
