package sqlc

import (
	"context"
	"testing"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/testhelper"
)

func TestCreateTxInvoiceHeader(t *testing.T) {
	t.Cleanup(func() {
		cleanInvoiceHeadersData(t)
	})

	ctx := testhelper.Ctx(t)
	db := OpenDB(ctx, t)
	defer db.Close()

	tx, err := db.Begin(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	input := &domain.InvoiceHeader{
		ClientID: 1,
	}

	ih := NewInvoiceHeader(db)
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
	ctx := testhelper.Ctx(t)
	db := OpenDB(ctx, t)
	defer db.Close()

	ih := NewInvoiceHeader(db)

	err := ih.DeleteAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
