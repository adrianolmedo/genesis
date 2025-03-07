package pq

import (
	"database/sql"
	"fmt"

	domain "github.com/adrianolmedo/genesis"
)

// Invoice repository.
type Invoice struct {
	db     *sql.DB
	header InvoiceHeader
	items  InvoiceItem
}

// Create generate a full Invoice.
func (i Invoice) Create(inv *domain.Invoice) error {
	tx, err := i.db.Begin()
	if err != nil {
		return err
	}

	if err := i.header.Create(tx, inv.Header); err != nil {
		tx.Rollback()
		return fmt.Errorf("invoice header: %w", err)
	}
	//fmt.Printf("invoice created with id %d\n", inv.Header.ID)

	if err := i.items.Create(tx, inv.Header.ID, inv.Items); err != nil {
		tx.Rollback()
		return fmt.Errorf("invoice items: %w", err)
	}
	//fmt.Printf("items added %d\n", len(inv.Items))

	return tx.Commit()
}
