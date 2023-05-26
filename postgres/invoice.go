package postgres

import (
	"database/sql"
	"fmt"

	domain "github.com/adrianolmedo/go-restapi"
)

type Invoice struct {
	db     *sql.DB
	header InvoiceHeader
	items  InvoiceItem
}

func (i Invoice) Create(invoice *domain.Invoice) error {
	tx, err := i.db.Begin()
	if err != nil {
		return err
	}

	if err := i.header.Create(tx, invoice.Header); err != nil {
		tx.Rollback()
		return fmt.Errorf("invoice header: %w", err)
	}
	//fmt.Printf("invoice created with id %d\n", invoice.Header.ID)

	if err := i.items.Create(tx, invoice.Header.ID, invoice.Items); err != nil {
		tx.Rollback()
		return fmt.Errorf("invoice items: %w", err)
	}
	//fmt.Printf("items added %d\n", len(invoice.Items))

	return tx.Commit()
}
