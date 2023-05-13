//go:build integration
// +build integration

package postgres_test

import (
	"testing"

	"github.com/adrianolmedo/go-restapi/domain"
	"github.com/adrianolmedo/go-restapi/postgres"
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

	input := &domain.InvoiceHeader{
		ClientID: 1,
	}

	ih := postgres.NewInvoiceHeader(db)
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
