package sqlc

import (
	"testing"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/testhelper"
)

func TestCreateInvoice(t *testing.T) {
	t.Cleanup(func() {
		cleanProductsData(t)
		cleanInvoiceItemsData(t)
		cleanInvoiceHeadersData(t)
	})

	input := &domain.Invoice{
		Header: &domain.InvoiceHeader{
			ClientID: 1,
		},
		Items: domain.ItemList{
			&domain.InvoiceItem{ProductID: 1},
		},
	}

	ctx := testhelper.Ctx(t)
	db := OpenDB(ctx, t)
	defer db.Close()
	insertProductsData(ctx, t, db)

	//ih := NewInvoiceHeader(db)
	//ii := NewInvoiceItem(db)
	in := NewInvoice(db)

	if err := in.Create(ctx, input); err != nil {
		t.Fatal(err)
	}

	for _, item := range input.Items {
		if !(item.ID > 0) {
			t.Errorf("invoice for product %d not added", item.ProductID)
		}
	}
}

func cleanInvoiceItemsData(t *testing.T) {
	ctx := testhelper.Ctx(t)
	db := OpenDB(ctx, t)
	defer db.Close()

	ii := NewInvoiceItem(db)

	err := ii.DeleteAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
