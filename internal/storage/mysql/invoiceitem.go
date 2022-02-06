package mysql

import (
	"database/sql"
	"fmt"

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

func (InvoiceItemRepository) CreateTx(tx *sql.Tx, headerID int64, items domain.ItemList) error {
	stmt, err := tx.Prepare("INSERT INTO invoice_items(invoice_header_id, product_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range items {
		result, err := stmt.Exec(headerID, item.ProductID)
		if err != nil {
			return err
		}

		invoiceID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		item.ID = invoiceID
	}

	return nil
}

func (r InvoiceItemRepository) Delete(id int64) error {
	stmt, err := r.db.Prepare("DELETE FROM invoice_items WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrInvoiceItemNotFound
	}
	return nil
}

func (r InvoiceItemRepository) Reset() error {
	stmt, err := r.db.Prepare("ALTER TABLE invoice_items AUTO_INCREMENT = 1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("can't alter table: %v", err)
	}

	return nil
}
