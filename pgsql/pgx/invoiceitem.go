package pgx

import (
	"context"
	"fmt"

	domain "github.com/adrianolmedo/genesis"

	"github.com/jackc/pgx/v5"
)

// InvoiceItem repository.
type InvoiceItem struct {
	conn *pgx.Conn
}

// Create create item asociated to a header and product for the invoice.
func (InvoiceItem) Create(tx pgx.Tx, headerID int64, items domain.ItemList) error {
	for _, item := range items {
		err := tx.QueryRow(context.Background(),
			`INSERT INTO "invoice_item" (invoice_header_id, product_id) VALUES ($1, $2) RETURNING id, created_at`,
			headerID, item.ProductID).Scan(&item.ID, &item.CreatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAll delete all invoice items.
func (ii InvoiceItem) DeleteAll() error {
	_, err := ii.conn.Exec(context.Background(), `TRUNCATE TABLE "invoice_item" RESTART IDENTITY`)
	if err != nil {
		return fmt.Errorf("can't truncate table: %v", err)
	}

	return nil
}
