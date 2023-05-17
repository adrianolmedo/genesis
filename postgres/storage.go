package postgres

import (
	"fmt"

	"github.com/adrianolmedo/go-restapi/config"
)

type Storage struct {
	User    User
	Product Product
	Invoice Invoice
}

func NewStorage(dbcfg config.DB) (*Storage, error) {
	if dbcfg.Engine == "" {
		return nil, fmt.Errorf("database engine not selected")
	}

	if dbcfg.Engine == "postgres" {
		db, err := newDB(dbcfg)
		if err != nil {
			return nil, fmt.Errorf("postgres: %v", err)
		}

		return &Storage{
			User:    User{db: db},
			Product: Product{db: db},
			Invoice: Invoice{
				db:     db,
				header: InvoiceHeader{db: db},
				items:  InvoiceItem{db: db},
			},
		}, nil
	}

	return nil, fmt.Errorf("database engine '%s' not implemented", dbcfg.Engine)
}
