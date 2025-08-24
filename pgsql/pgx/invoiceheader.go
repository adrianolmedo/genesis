package pgx

import (
	"context"
	"fmt"

	"github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/billing"

	"github.com/jackc/pgx/v5"
)

// InvoiceHeader repository.
type InvoiceHeader struct {
	conn *pgx.Conn
}

func (InvoiceHeader) Create(ctx context.Context, tx pgx.Tx, m *billing.InvoiceHeader) error {
	m.UUID = genesis.NextUUID()

	err := tx.QueryRow(ctx,
		`INSERT INTO "invoice_header" (uuid, client_id) VALUES ($1, $2) RETURNING id, created_at`,
		m.UUID, m.ClientID).Scan(&m.ID, &m.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (ih InvoiceHeader) DeleteAll(ctx context.Context) error {
	_, err := ih.conn.Exec(ctx, `TRUNCATE TABLE "invoice_header" RESTART IDENTITY CASCADE`)
	if err != nil {
		return fmt.Errorf("can't truncate table: %v", err)
	}

	return nil
}
