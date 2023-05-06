//go:build integration
// +build integration

package storage_test

import (
	"testing"

	"github.com/adrianolmedo/go-restapi/domain"
	"github.com/adrianolmedo/go-restapi/storage"
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

	ih := storage.NewInvoiceHeaderRepository(db)
	ii := storage.NewInvoiceItemRepository(db)
	invoice := storage.NewInvoiceRepository(db, ih, ii)

	input := &domain.Invoice{
		Header: &domain.InvoiceHeader{
			ClientID: 1,
		},
		Items: domain.ItemList{
			&domain.InvoiceItem{ProductID: 1},
		},
	}

	if err := invoice.Create(input); err != nil {
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

	err := storage.NewInvoiceHeaderRepository(db).DeleteAll()
	if err != nil {
		t.Fatal(err)
	}
}

func cleanInvoiceItemsData(t *testing.T) {
	db := openDB(t)
	defer closeDB(t, db)

	err := storage.NewInvoiceItemRepository(db).DeleteAll()
	if err != nil {
		t.Fatal(err)
	}
}
