package postgres

import (
	"database/sql"
	"fmt"

	"github.com/adrianolmedo/go-restapi/domain"
)

type InvoiceHeader struct {
	db *sql.DB
}

func (InvoiceHeader) Create(tx *sql.Tx, m *domain.InvoiceHeader) error {
	stmt, err := tx.Prepare("INSERT INTO invoice_headers(client_id) VALUES ($1) RETURNING id, created_at")
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRow(m.ClientID).Scan(&m.ID, &m.CreatedAt)
}

func (ih InvoiceHeader) DeleteAll() error {
	stmt, err := ih.db.Prepare("TRUNCATE TABLE invoice_headers RESTART IDENTITY CASCADE")
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
