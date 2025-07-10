//go:build integration
// +build integration

package sqlc

import (
	"context"
	"testing"

	domain "github.com/adrianolmedo/genesis"
)

func TestCreateTxInvoiceHeader(t *testing.T) {
	t.Cleanup(func() {
		cleanInvoiceHeadersData(t)
	})

	conn := openDB(t)
	defer closeDB(t, conn)

	tx, err := conn.Begin(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	input := &domain.InvoiceHeader{
		ClientID: 1,
	}

	ih := NewInvoiceHeader(conn)
	if err := ih.Create(tx, input); err != nil {
		tx.Rollback(context.Background())
		t.Fatal(err)
	}

	if !(input.ID > 0) {
		t.Fatal("invoice header not created")
	}

	if err := tx.Commit(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func cleanInvoiceHeadersData(t *testing.T) {
	conn := openDB(t)
	defer closeDB(t, conn)

	ih := NewInvoiceHeader(conn)

	err := ih.DeleteAll()
	if err != nil {
		t.Fatal(err)
	}
}
