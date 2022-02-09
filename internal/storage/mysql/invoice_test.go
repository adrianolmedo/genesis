//go:build integration
// +build integration

package mysql_test

import (
	"database/sql"
	"testing"

	"github.com/adrianolmedo/go-restapi/internal/domain"
	"github.com/adrianolmedo/go-restapi/internal/storage/mysql"
)

func TestCreateInvoice(t *testing.T) {
	db := openDB(t)
	defer closeDB(t, db)
	insertProductsData(t, db)

	ih := mysql.NewInvoiceHeaderRepository(db)
	ii := mysql.NewInvoiceItemRepository(db)
	invoice := mysql.NewInvoiceRepository(db, ih, ii)

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

	cleanInvoiceItemsData(t, db, input.Items)
	cleanInvoiceHeadersData(t, db, input.Header.ID)
	cleanProductsData(t, db, 2)
	cleanProductsData(t, db, 1)
}

func cleanInvoiceItemsData(t *testing.T, db *sql.DB, items domain.ItemList) {
	ii := mysql.NewInvoiceItemRepository(db)

	for _, item := range items {
		if err := ii.Delete(item.ID); err != nil {
			t.Errorf("delete item id %d: %s", item.ID, err)
		}
	}

	if err := ii.Reset(); err != nil {
		t.Fatal(err)
	}
}
