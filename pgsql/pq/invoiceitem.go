package pq

import (
	"database/sql"
	"fmt"

	domain "github.com/adrianolmedo/genesis"
)

// InvoiceItem repository.
type InvoiceItem struct {
	db *sql.DB
}

// Create create item asociated to a header and product for the invoice.
func (InvoiceItem) Create(tx *sql.Tx, headerID int, items domain.ItemList) error {
	stmt, err := tx.Prepare(`INSERT INTO "invoice_item" (invoice_header_id, product_id) VALUES ($1, $2) RETURNING id, created_at`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range items {
		err = stmt.QueryRow(headerID, item.ProductID).Scan(&item.ID, &item.CreatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAll delete all invoice items.
// TODO: Test.
func (ii InvoiceItem) DeleteAll() error {
	stmt, err := ii.db.Prepare(`TRUNCATE TABLE "invoice_item" RESTART IDENTITY CASCADE`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("can't truncate table: %v", err)
	}

	return nil
}
