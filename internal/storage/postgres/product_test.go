//go:build integration
// +build integration

package postgres_test

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/adrianolmedo/go-restapi/internal/domain"
	"github.com/adrianolmedo/go-restapi/internal/storage/postgres"
)

func TestCreateProduct(t *testing.T) {
	t.Cleanup(func() {
		cleanProductsData(t)
	})

	db := openDB(t)
	defer closeDB(t, db)
	p := postgres.NewProductRepository(db)

	input := &domain.Product{
		Name:         "Coca-Cola",
		Observations: "",
		Price:        3,
	}

	if err := p.Create(input); err != nil {
		t.Fatal(err)
	}

	product, err := p.ByID(input.ID)
	if err != nil {
		t.Fatal(err)
	}

	if product.CreatedAt.IsZero() {
		t.Error("expected created at")
	}

	if !product.UpdatedAt.IsZero() {
		t.Error("unexpected updated at")
	}
}

func TestProductByID(t *testing.T) {
	t.Cleanup(func() {
		cleanProductsData(t)
	})

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
		got, err := postgres.NewProductRepository(db).ByID(tc.input)
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
}

func TestUpdateProduct(t *testing.T) {
	t.Cleanup(func() {
		cleanProductsData(t)
	})

	db := openDB(t)
	defer closeDB(t, db)
	insertProductsData(t, db)

	input := domain.Product{
		ID:           1,
		Name:         "Coca-Cola",
		Observations: "",
		Price:        3,
	}

	pr := postgres.NewProductRepository(db)

	if err := pr.Update(input); err != nil {
		t.Fatal(err)
	}

	product, err := pr.ByID(input.ID)
	if err != nil {
		t.Fatal(err)
	}

	if product.Name != input.Name {
		t.Errorf("Name: want %s, got %s", input.Name, product.Name)
	}

	if product.Observations != input.Observations {
		t.Errorf("LastName: want %s, got %s", input.Observations, product.Observations)
	}

	if product.CreatedAt.IsZero() {
		t.Error("expected created at")
	}

	if product.UpdatedAt.IsZero() {
		t.Error("expected updated at")
	}
}

func insertProductsData(t *testing.T, db *sql.DB) {
	p := postgres.NewProductRepository(db)

	if err := p.Create(&domain.Product{
		Name:         "Coca-Cola",
		Observations: "",
		Price:        3,
	}); err != nil {
		t.Fatal(err)
	}

	if err := p.Create(&domain.Product{
		Name:         "Big-Cola",
		Observations: "Made in Venezuela",
		Price:        2,
	}); err != nil {
		t.Fatal(err)
	}
}

func cleanProductsData(t *testing.T) {
	db := openDB(t)
	defer closeDB(t, db)

	err := postgres.NewProductRepository(db).DeleteAll()
	if err != nil {
		t.Fatal(err)
	}
}
