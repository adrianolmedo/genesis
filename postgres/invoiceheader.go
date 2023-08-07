package postgres

import (
	"database/sql"
	"fmt"

	domain "github.com/adrianolmedo/genesis"
)

// InvoiceHeader repository.
type InvoiceHeader struct {
	db *sql.DB
}

// Create insert the header invoice.
func (InvoiceHeader) Create(tx *sql.Tx, m *domain.InvoiceHeader) error {
	stmt, err := tx.Prepare(`INSERT INTO "invoice_header" (uuid, client_id) VALUES ($1, $2) RETURNING id, created_at`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	m.UUID = domain.NextUUID()

	return stmt.QueryRow(m.UUID, m.ClientID).Scan(&m.ID, &m.CreatedAt)
}

// DeleteAll delete all invoice header.
// TODO: Move to test.
func (ih InvoiceHeader) DeleteAll() error {
	stmt, err := ih.db.Prepare(`TRUNCATE TABLE "invoice_header" RESTART IDENTITY CASCADE`)
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
