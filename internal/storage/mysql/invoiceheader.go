package mysql

import (
	"database/sql"
	"fmt"

	"github.com/adrianolmedo/genesis/internal/domain"
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
	stmt, err := tx.Prepare("INSERT INTO invoice_headers(client_id) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(header.ClientID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	header.ID = id // Important!
	return nil
}

func (r InvoiceHeaderRepository) Delete(id int64) error {
	stmt, err := r.db.Prepare("DELETE FROM invoice_headers WHERE id = ?")
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
		return domain.ErrInvoiceHeaderNotFound
	}
	return nil
}

func (r InvoiceHeaderRepository) Reset() error {
	stmt, err := r.db.Prepare("ALTER TABLE invoice_headers AUTO_INCREMENT = 1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("can't alter table: %v", err)
	}

	return nil
}
