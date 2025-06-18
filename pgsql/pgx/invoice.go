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
func (i Invoice) Create(inv *domain.Invoice) error {
	tx, err := i.conn.Begin(context.Background())
	if err != nil {
		return err
	}

	err = i.header.Create(tx, inv.Header)
	if err != nil {
		tx.Rollback(context.Background())
		return fmt.Errorf("invoice header: %w", err)
	}

	err = i.items.Create(tx, inv.Header.ID, inv.Items)
	if err != nil {
		tx.Rollback(context.Background())
		return fmt.Errorf("invoice items: %w", err)
	}

	return tx.Commit(context.Background())
}
