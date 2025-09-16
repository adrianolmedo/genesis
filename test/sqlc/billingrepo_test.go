package sqlc

import (
	"context"
	"testing"

	"github.com/adrianolmedo/genesis/billing"
	"github.com/adrianolmedo/genesis/test"
)

func TestCreateInvoice(t *testing.T) {
	t.Cleanup(func() {
		cleanProductsData(t)
		cleanInvoiceItemsData(t)
		cleanInvoiceHeadersData(t)
	})
	input := &billing.Invoice{
		Header: &billing.InvoiceHeader{
			ClientID: 1,
		},
		Items: billing.ItemList{
			billing.InvoiceItem{ProductID: 1},
		},
	}
	ctx := test.Ctx(t)
	db := openDB(ctx, t)
	defer db.Close()
	insertProductsData(ctx, t, db)
	//ih := NewInvoiceHeader(db)
	//ii := NewInvoiceItem(db)
	in := billing.NewRepo(db)
	if err := in.CreateInvoice(ctx, input); err != nil {
		t.Fatal(err)
	}
	for _, item := range input.Items {
		if !(item.ID > 0) {
			t.Errorf("invoice for product %d not added", item.ProductID)
		}
	}
}

func TestCreateTxInvoiceHeader(t *testing.T) {
	t.Cleanup(func() {
		cleanInvoiceHeadersData(t)
	})
	ctx := test.Ctx(t)
	db := openDB(ctx, t)
	defer db.Close()
	tx, err := db.Begin(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	input := &billing.InvoiceHeader{
		ClientID: 1,
	}
	r := billing.NewRepo(db)
	if err := r.CreateHeader(ctx, tx, input); err != nil {
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
	ctx := test.Ctx(t)
	db := openDB(ctx, t)
	defer db.Close()
	ih := billing.NewRepo(db)
	err := ih.DeleteAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func cleanInvoiceItemsData(t *testing.T) {
	ctx := test.Ctx(t)
	db := openDB(ctx, t)
	defer db.Close()
	ii := billing.NewRepo(db)
	err := ii.DeleteAllItems(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
