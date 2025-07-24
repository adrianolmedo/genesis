package pgx

import (
	"context"
	"fmt"

	domain "github.com/adrianolmedo/genesis"

	"github.com/jackc/pgx/v5"
)

// Invoice repository.
type Invoice struct {
	conn   *pgx.Conn
	header InvoiceHeader
	items  InvoiceItem
}

// Create generate a full Invoice.
func (i Invoice) Create(ctx context.Context, inv *domain.Invoice) error {
	tx, err := i.conn.Begin(ctx)
	if err != nil {
		return err
	}

	// Create invoice header
	err = i.header.Create(ctx, tx, inv.Header)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("invoice header: %w", err)
	}

	// Create invoice items
	err = i.items.Create(ctx, tx, inv.Header.ID, inv.Items)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("invoice items: %w", err)
	}

	return tx.Commit(ctx)
}
