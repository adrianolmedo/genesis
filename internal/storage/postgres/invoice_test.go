package postgres_test

import (
	"testing"

	"github.com/adrianolmedo/go-restapi/internal/domain"
	"github.com/adrianolmedo/go-restapi/internal/storage/postgres"
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

	ih := postgres.NewInvoiceHeaderRepository(db)
	ii := postgres.NewInvoiceItemRepository(db)
	invoice := postgres.NewInvoiceRepository(db, ih, ii)

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

func cleanInvoiceItemsData(t *testing.T) {
	db := openDB(t)
	defer closeDB(t, db)

	err := postgres.NewInvoiceItemRepository(db).DeleteAll()
	if err != nil {
		t.Fatal(err)
	}
}
