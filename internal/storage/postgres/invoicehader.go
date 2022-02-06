package postgres

import (
	"database/sql"

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

func (InvoiceHeaderRepository) CreateTx(*sql.Tx, *domain.InvoiceHeader) error {
	return nil
}
