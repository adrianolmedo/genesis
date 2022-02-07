package postgres

import (
	"database/sql"
	"fmt"

	"github.com/adrianolmedo/go-restapi/internal/domain"
)

type InvoiceHeaderRepository struct {
	db *sql.DB
}

func NewInvoiceHeaderRepository(db *sql.DB) *InvoiceHeaderRepository {
	return &InvoiceHeaderRepository{
		db: db,
	}
}

func (InvoiceHeaderRepository) CreateTx(tx *sql.Tx, header *domain.InvoiceHeader) error {
	stmt, err := tx.Prepare("INSERT INTO invoice_headers(client_id) VALUES ($1) RETURNING id, created_at")
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRow(header.ClientID).Scan(&header.ID, &header.CreatedAt)
}

func (r InvoiceHeaderRepository) DeleteAll() error {
	stmt, err := r.db.Prepare("TRUNCATE TABLE invoice_headers RESTART IDENTITY CASCADE")
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
