package postgres

import (
	"database/sql"
	"fmt"

	"github.com/adrianolmedo/go-restapi/domain"
)

type Invoice struct {
	db     *sql.DB
	header InvoiceHeader
	items  InvoiceItem
}

func NewInvoice(
	db *sql.DB,
	header InvoiceHeader,
	items InvoiceItem) Invoice {
	return Invoice{
		db:     db,
		header: header,
		items:  items,
	}
}

func (r Invoice) Create(invoice *domain.Invoice) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if err := r.header.CreateTx(tx, invoice.Header); err != nil {
		tx.Rollback()
		return fmt.Errorf("invoice header: %w", err)
	}
	//fmt.Printf("invoice created with id %d\n", invoice.Header.ID)

	if err := r.items.CreateTx(tx, invoice.Header.ID, invoice.Items); err != nil {
		tx.Rollback()
		return fmt.Errorf("invoice items: %w", err)
	}
	//fmt.Printf("items added %d\n", len(invoice.Items))

	return tx.Commit()
}
