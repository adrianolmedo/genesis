package mysql

import (
	"database/sql"
	"fmt"

	"github.com/adrianolmedo/go-restapi/internal/domain"
)

type InvoiceRepository struct {
	db     *sql.DB
	header *InvoiceHeaderRepository
	items  *InvoiceItemRepository
}

func NewInvoiceRepository(db *sql.DB, header *InvoiceHeaderRepository, items *InvoiceItemRepository) *InvoiceRepository {
	return &InvoiceRepository{
		db:     db,
		header: header,
		items:  items,
	}
}

func (r InvoiceRepository) Create(invoice *domain.Invoice) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	err = r.header.CreateTx(tx, invoice.Header)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("invoice header %s", err)
	}

	err = r.items.CreateTx(tx, invoice.Header.ID, invoice.Items)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("invoice items %s", err)
	}

	return tx.Commit()
}
