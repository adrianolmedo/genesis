//go:build integration
// +build integration

package mysql_test

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/adrianolmedo/genesis/internal/domain"
	"github.com/adrianolmedo/genesis/internal/storage/mysql"
)

func TestCreateProduct(t *testing.T) {
	db := openDB(t)
	defer closeDB(t, db)
	pr := mysql.NewProductRepository(db)

	input := &domain.Product{
		Name:         "Coca-Cola",
		Observations: "",
		Price:        3,
	}

	if err := pr.Create(input); err != nil {
		t.Fatal(err)
	}

	product, err := pr.ByID(input.ID)
	if err != nil {
		t.Fatal(err)
	}

	if product.CreatedAt.IsZero() {
		t.Error("expected created at")
	}

	if !product.UpdatedAt.IsZero() {
		t.Error("unexpected updated at")
	}

	cleanProductsData(t, db, input.ID)
}

func TestProductByID(t *testing.T) {
	db := openDB(t)
	defer closeDB(t, db)
	insertProductsData(t, db)

	tt := []struct {
		name           string
		input          int64
		wantName       string
		errExpected    bool
		wantErrContain string
	}{
		{
			name:           "ok",
			input:          1,
			wantName:       "Coca-Cola",
			errExpected:    false,
			wantErrContain: "",
		},
		{
			name:           "product-not-found",
			input:          3,
			wantName:       "Big-Cola",
			errExpected:    true,
			wantErrContain: "product not found",
		},
	}

	for _, tc := range tt {
		got, err := mysql.NewProductRepository(db).ByID(tc.input)
		if (err != nil) != tc.errExpected {
			t.Fatalf("%s: ByID(%d): unexpected error status: %v", tc.name, tc.input, err)
		}

		if err != nil && !strings.Contains(err.Error(), tc.wantErrContain) {
			t.Fatalf("want error string %q to contain %q", err.Error(), tc.wantErrContain)
		}

		if !tc.errExpected && tc.wantName != got.Name {
			t.Fatalf("%s: ByID(%d): want %s, got %s", tc.name, tc.input, tc.wantName, got.Name)
		}
	}

	cleanProductsData(t, db, 2)
	cleanProductsData(t, db, 1)
}

func insertProductsData(t *testing.T, db *sql.DB) {
	pr := mysql.NewProductRepository(db)

	if err := pr.Create(&domain.Product{
		Name:         "Coca-Cola",
		Observations: "",
		Price:        3,
	}); err != nil {
		t.Fatal(err)
	}

	if err := pr.Create(&domain.Product{
		Name:         "Big-Cola",
		Observations: "Made in Venezuela",
		Price:        2,
	}); err != nil {
		t.Fatal(err)
	}
}

func cleanProductsData(t *testing.T, db *sql.DB, productID int64) {
	pr := mysql.NewProductRepository(db)

	if err := pr.Delete(productID); err != nil {
		t.Fatal(err)
	}

	if err := pr.Reset(); err != nil {
		t.Fatal(err)
	}
}
