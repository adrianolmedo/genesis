package postgres_test

import (
	"testing"

	"github.com/adrianolmedo/go-restapi/internal/domain"
	"github.com/adrianolmedo/go-restapi/internal/storage/postgres"
)

func TestCreateTxInvoiceHeader(t *testing.T) {
	t.Cleanup(func() {
		cleanInvoiceHeadersData(t)
	})

	db := openDB(t)
	defer closeDB(t, db)

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	ih := postgres.NewInvoiceHeaderRepository(db)

	input := &domain.InvoiceHeader{
		ClientID: 1,
	}

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
}

func cleanInvoiceHeadersData(t *testing.T) {
	db := openDB(t)
	defer closeDB(t, db)

	err := postgres.NewInvoiceHeaderRepository(db).DeleteAll()
	if err != nil {
		t.Fatal(err)
	}
}
