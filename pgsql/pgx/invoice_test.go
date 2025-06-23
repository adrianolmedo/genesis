//go:build integration
// +build integration

package pgx

import (
	"testing"

	domain "github.com/adrianolmedo/genesis"
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

	conn := openDB(t)
	defer closeDB(t, conn)
	insertProductsData(t, conn)

	ih := InvoiceHeader{conn: conn}
	ii := InvoiceItem{conn: conn}
	in := Invoice{conn: conn, header: ih, items: ii}

	if err := in.Create(input); err != nil {
		t.Fatal(err)
	}

	for _, item := range input.Items {
		if !(item.ID > 0) {
			t.Errorf("invoice for product %d not added", item.ProductID)
		}
	}
}

func cleanInvoiceItemsData(t *testing.T) {
	conn := openDB(t)
	defer closeDB(t, conn)

	ii := InvoiceItem{conn: conn}

	err := ii.DeleteAll()
	if err != nil {
		t.Fatal(err)
	}
}
