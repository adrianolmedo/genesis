package sqlc

import (
	"context"
	"fmt"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/pgsql/sqlc/dbgen"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Invoice repository.
type Invoice struct {
	db     *pgxpool.Pool
	q      *dbgen.Queries // for non-tx operations
	header *InvoiceHeader
	items  *InvoiceItem
}

// NewInvoice creates a new Invoice repository instance.
// Since [InvoiceHader], [InvoiceItem] (or any other repository in question)
// are closely related to Invoice, they are created as part of the Invoice structure.
func NewInvoice(db *pgxpool.Pool) *Invoice {
	q := dbgen.New(db)
	return &Invoice{
		db:     db,
		q:      q,
		header: NewInvoiceHeader(db),
		items:  NewInvoiceItem(db),
	}
}

func (i Invoice) Create(ctx context.Context, inv *domain.Invoice) error {
	tx, err := i.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Create invoice header
	if err := i.header.Create(ctx, tx, inv.Header); err != nil {
		return fmt.Errorf("invoice header: %w", err)
	}

	// Create invoice items
	if err := i.items.Create(ctx, tx, inv.Header.ID, inv.Items); err != nil {
		return fmt.Errorf("invoice items: %w", err)
	}

	return tx.Commit(ctx)
}
