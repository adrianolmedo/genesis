package postgres

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
	stmt, err := tx.Prepare("INSERT INTO invoice_items (invoice_header_id, product_id) VALUES ($1, $2) RETURNING id, created_at")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range items {
		err = stmt.QueryRow(headerID, item.ProductID).Scan(&item.ID, &item.CreatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r InvoiceItemRepository) DeleteAll() error {
	stmt, err := r.db.Prepare("TRUNCATE TABLE invoice_items RESTART IDENTITY CASCADE")
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
