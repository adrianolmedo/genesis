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

	ctx := testCtx(t)
	conn := openDB(ctx, t)
	defer closeDB(ctx, t, conn)

	tx, err := conn.Begin(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	input := &domain.InvoiceHeader{
		ClientID: 1,
	}

	ih := NewInvoiceHeader(conn)
	if err := ih.Create(ctx, tx, input); err != nil {
		tx.Rollback(ctx)
		t.Fatal(err)
	}

	if !(input.ID > 0) {
		t.Fatal("invoice header not created")
	}

	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
}

func cleanInvoiceHeadersData(t *testing.T) {
	ctx := testCtx(t)
	conn := openDB(ctx, t)
	defer closeDB(ctx, t, conn)

	ih := NewInvoiceHeader(conn)

	err := ih.DeleteAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
