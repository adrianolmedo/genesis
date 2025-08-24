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
func NewRepo(db *pgxpool.Pool) *Repo {
	q := dbgen.New(db)
	return &Repo{
		db: db,
		q:  q,
	}
}

// CreateInvoice creates a new invoice with its header and items.
func (r Repo) CreateInvoice(ctx context.Context, inv *Invoice) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Create invoice header
	if err := r.CreateHeader(ctx, tx, inv.Header); err != nil {
		return fmt.Errorf("invoice header: %w", err)
	}

	// Create invoice items
	if err := r.CreateItem(ctx, tx, inv.Header.ID, inv.Items); err != nil {
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

// DeleteAll delete all invoice headers (permanantly).
func (r Repo) DeleteAll(ctx context.Context) error {
	err := r.q.InvoiceHeaderDeleteAll(ctx)
	if err != nil {
		return fmt.Errorf("can't truncate table: %v", err)
	}
	return nil
}

// DeleteAllItems deletes all invoice items.
// This is used for testing purposes to reset the state of the invoice items table.
func (r Repo) DeleteAllItems(ctx context.Context) error {
	err := r.q.InvoiceItemDeleteAll(ctx)
	if err != nil {
		return fmt.Errorf("can't truncate table: %v", err)
	}
	return nil
}
