package mysql_test

import (
	"database/sql"
	"testing"

	"github.com/adrianolmedo/go-restapi/internal/domain"
	"github.com/adrianolmedo/go-restapi/internal/storage/mysql"
)

func TestCreateTxInvoiceHeader(t *testing.T) {
	db := openDB(t)
	defer closeDB(t, db)

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	input := &domain.InvoiceHeader{
		ClientID: 1,
	}

	ih := mysql.NewInvoiceHeaderRepository(db)
	if err := ih.CreateTx(tx, input); err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	if !(input.ID > 0) {
		t.Fatal("invoice header not created")
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	cleanInvoiceHeadersData(t, db, input.ID)
}

func cleanInvoiceHeadersData(t *testing.T, db *sql.DB, headerID int64) {
	ih := mysql.NewInvoiceHeaderRepository(db)

	if err := ih.Delete(headerID); err != nil {
		t.Fatal(err)
	}

	if err := ih.Reset(); err != nil {
		t.Fatal(err)
	}
}
