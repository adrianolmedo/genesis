package billing

import (
	"context"
	"fmt"

	"github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/pgsql/sqlc/dbgen"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pborman/uuid"
)

type Repo struct {
	db *pgxpool.Pool
	q  *dbgen.Queries // for non-tx operations
}

// NewInvoice creates a new Invoice repository instance.
// Since [InvoiceHader], [InvoiceItem] (or any other repository in question)
// are closely related to Invoice, they are created as part of the Invoice structure.
func NewRepo(db *pgxpool.Pool) *Repo {
	q := dbgen.New(db)
	return &Repo{
		db: db,
		q:  q,
	}
}

// CreateInvoice creates a new invoice with its header and items.
func (i Repo) CreateInvoice(ctx context.Context, inv *Invoice) error {
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
	if err := i.CreateHeader(ctx, tx, inv.Header); err != nil {
		return fmt.Errorf("invoice header: %w", err)
	}

	// Create invoice items
	if err := i.CreateItem(ctx, tx, inv.Header.ID, inv.Items); err != nil {
		return fmt.Errorf("invoice items: %w", err)
	}

	return tx.Commit(ctx)
}

// CreateHeader creates a new invoice header in the database.
func (r Repo) CreateHeader(ctx context.Context, tx pgx.Tx, m *InvoiceHeader) error {
	m.UUID = genesis.NextUUID()

	row, err := r.q.WithTx(tx).InvoiceHeaderCreate(ctx, dbgen.InvoiceHeaderCreateParams{
		Uuid:     uuid.Parse(m.UUID),
		ClientID: m.ClientID,
	})
	if err != nil {
		return err
	}

	m.ID = row.ID
	return nil
}

// CreateItem creates items associated with a header and product for the invoice.
func (r Repo) CreateItem(ctx context.Context, tx pgx.Tx, headerID int64, items ItemList) error {
	for _, item := range items {
		row, err := r.q.WithTx(tx).InvoiceItemCreate(ctx, dbgen.InvoiceItemCreateParams{
			InvoiceHeaderID: headerID,
			ProductID:       item.ProductID,
		})
		if err != nil {
			return err
		}

		item.ID = row.ID
		item.CreatedAt = row.CreatedAt

	}
	return nil
}
