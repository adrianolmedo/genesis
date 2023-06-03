//go:build integration
// +build integration

package postgres

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

	db := openDB(t)
	defer closeDB(t, db)
	insertProductsData(t, db)

	ih := InvoiceHeader{db: db}
	ii := InvoiceItem{db: db}
	in := Invoice{db: db, header: ih, items: ii}

	input := &domain.Invoice{
		Header: &domain.InvoiceHeader{
			ClientID: 1,
		},
		Items: domain.ItemList{
			&domain.InvoiceItem{ProductID: 1},
		},
	}

	if err := in.Create(input); err != nil {
		t.Fatal(err)
	}

	for _, item := range input.Items {
		if !(item.ID > 0) {
			t.Errorf("invoice for product %d not added", item.ProductID)
		}
	}
}

func cleanInvoiceHeadersData(t *testing.T) {
	db := openDB(t)
	defer closeDB(t, db)

	ih := InvoiceHeader{db: db}

	err := ih.DeleteAll()
	if err != nil {
		t.Fatal(err)
	}
}

func cleanInvoiceItemsData(t *testing.T) {
	db := openDB(t)
	defer closeDB(t, db)

	ii := InvoiceItem{db: db}

	err := ii.DeleteAll()
	if err != nil {
		t.Fatal(err)
	}
}
