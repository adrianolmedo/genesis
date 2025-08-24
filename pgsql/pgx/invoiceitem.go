package pgx

import (
	"context"
	"fmt"

	"github.com/adrianolmedo/genesis/billing"
	"github.com/jackc/pgx/v5"
)

// InvoiceItem repository.
type InvoiceItem struct {
	conn *pgx.Conn
}

// Create create item asociated to a header and product for the invoice.
func (InvoiceItem) Create(ctx context.Context, tx pgx.Tx, headerID int64, items billing.ItemList) error {
	for _, item := range items {
		err := tx.QueryRow(ctx,
			`INSERT INTO "invoice_item" (invoice_header_id, product_id) VALUES ($1, $2) RETURNING id, created_at`,
			headerID, item.ProductID).Scan(&item.ID, &item.CreatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAll delete all invoice items.
func (ii InvoiceItem) DeleteAll(ctx context.Context) error {
	_, err := ii.conn.Exec(ctx, `TRUNCATE TABLE "invoice_item" RESTART IDENTITY`)
	if err != nil {
		return fmt.Errorf("can't truncate table: %v", err)
	}

	return nil
}
