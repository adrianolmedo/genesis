//go:build integration
// +build integration

package pq

import (
	"testing"

	domain "github.com/adrianolmedo/genesis"
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

	ih := InvoiceHeader{db: db}
	if err := ih.Create(tx, input); err != nil {
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
