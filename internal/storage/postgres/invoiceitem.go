package postgres

import (
	"database/sql"

	"github.com/adrianolmedo/go-restapi/internal/domain"
)

type InvoiceItemRepository struct {
	db *sql.DB
}

func NewInvoiceItemRepository(db *sql.DB) *InvoiceItemRepository {
	return &InvoiceItemRepository{
		db: db,
	}
}

func (InvoiceItemRepository) CreateTx(*sql.Tx, int64, domain.InvoiceItem) error {
	return nil
}
