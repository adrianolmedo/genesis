//go:build integration
// +build integration

package sqlc

import (
	"strings"
	"testing"

	domain "github.com/adrianolmedo/genesis"

	"github.com/jackc/pgx/v5"
)

func TestCreateProduct(t *testing.T) {
	t.Cleanup(func() {
		cleanProductsData(t)
	})

	conn := openDB(t)
	defer closeDB(t, conn)

	p := NewProduct(conn)

	input := &domain.Product{
		Name:         "Coca-Cola",
		Observations: "",
		Price:        3,
	}

	if err := p.Create(input); err != nil {
		t.Fatal(err)
	}

	got, err := p.ByID(1)
	if err != nil {
		t.Fatal(err)
	}

	if got.CreatedAt.IsZero() {
		t.Error("expected created at")
	}

	if got.UpdatedAt != nil {
		t.Error("unexpected updated at")
	}

	if got.DeletedAt != nil {
		t.Error("unexpected deleted at")
	}
}

func TestProductByID(t *testing.T) {
	t.Cleanup(func() {
		cleanProductsData(t)
	})

	conn := openDB(t)
	defer closeDB(t, conn)
	insertProductsData(t, conn)

	tt := []struct {
		name           string
		input          int64
		wantName       string
		errExpected    bool
		wantErrContain string
	}{
		{
			name:           "ok-test", // test name
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

	p := NewProduct(conn)

	for _, tc := range tt {
		got, err := p.ByID(tc.input)
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

	conn := openDB(t)
	defer closeDB(t, conn)
	insertProductsData(t, conn)

	input := domain.Product{
		ID:           1,
		Name:         "Coca-Cola",
		Observations: "",
		Price:        3,
	}

	p := NewProduct(conn)

	if err := p.Update(input); err != nil {
		t.Fatal(err)
	}

	got, err := p.ByID(input.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.Name != input.Name {
		t.Errorf("Name: want %s, got %s", input.Name, got.Name)
	}

	if got.Observations != input.Observations {
		t.Errorf("Observations: want %s, got %s", input.Observations, got.Observations)
	}

	if got.CreatedAt.IsZero() {
		t.Error("expected created at")
	}

	if got.UpdatedAt.IsZero() {
		t.Error("expected updated at")
	}

	if got.DeletedAt != nil {
		t.Error("unexpected deleted at")
	}
}

// insertProductsData add default `product` data.
func insertProductsData(t *testing.T, conn *pgx.Conn) {
	p := NewProduct(conn)

	// Add first product
	if err := p.Create(&domain.Product{
		Name:         "Coca-Cola",
		Observations: "",
		Price:        3,
	}); err != nil {
		t.Fatal(err)
	}

	// Add second product
	if err := p.Create(&domain.Product{
		Name:         "Big-Cola",
		Observations: "Made in Venezuela",
		Price:        2,
	}); err != nil {
		t.Fatal(err)
	}
}

// cleanProductsData delete all rows of `product` table.
func cleanProductsData(t *testing.T) {
	conn := openDB(t)
	defer closeDB(t, conn)

	p := NewProduct(conn)

	err := p.DeleteAll()
	if err != nil {
		t.Fatal(err)
	}
}
