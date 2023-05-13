//go:build integration
// +build integration

package postgres_test

import (
	"testing"

	"github.com/adrianolmedo/go-restapi/domain"
	"github.com/adrianolmedo/go-restapi/postgres"
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

	ih := postgres.NewInvoiceHeader(db)
	ii := postgres.NewInvoiceItem(db)
	invoice := postgres.NewInvoice(db, ih, ii)

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

	err := postgres.NewInvoiceHeader(db).DeleteAll()
	if err != nil {
		t.Fatal(err)
	}
}

func cleanInvoiceItemsData(t *testing.T) {
	db := openDB(t)
	defer closeDB(t, db)

	err := postgres.NewInvoiceItem(db).DeleteAll()
	if err != nil {
		t.Fatal(err)
	}
}
