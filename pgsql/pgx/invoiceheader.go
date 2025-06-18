package pgx

import (
	"context"
	"fmt"

	domain "github.com/adrianolmedo/genesis"

	"github.com/jackc/pgx/v5"
)

// InvoiceHeader repository.
type InvoiceHeader struct {
	conn *pgx.Conn
}

func (InvoiceHeader) Create(tx pgx.Tx, m *domain.InvoiceHeader) error {
	m.UUID = domain.NextUUID()

	err := tx.QueryRow(context.Background(),
		`INSERT INTO "invoice_header" (uuid, client_id) VALUES ($1, $2) RETURNING id, created_at`,
		m.UUID, m.ClientID).Scan(&m.ID, &m.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (ih InvoiceHeader) DeleteAll() error {
	_, err := ih.conn.Exec(context.Background(), `TRUNCATE TABLE "invoice_header" RESTART IDENTITY CASCADE`)
	if err != nil {
		return fmt.Errorf("can't truncate table: %v", err)
	}

	return nil
}
