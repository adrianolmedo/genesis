package service_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/adrianolmedo/go-restapi/internal/domain"
	"github.com/adrianolmedo/go-restapi/internal/mock"
	"github.com/adrianolmedo/go-restapi/internal/service"
	"github.com/adrianolmedo/go-restapi/internal/storage"
)

func TestAdd(t *testing.T) {
	tt := []struct {
		name           string
		input          *domain.Product
		mock           storage.ProductRepository
		errExpected    bool
		wantErrContain string
	}{
		{
			name: "successful",
			input: &domain.Product{
				Name:         "Coca-Cola",
				Observations: "",
				Price:        3,
			},
			mock:           mock.ProductRepositoryOk{},
			errExpected:    false,
			wantErrContain: "",
		},
		{
			name: "empty-field",
			input: &domain.Product{
				Name:         "",
				Observations: "Made in Venezuela",
				Price:        2,
			},
			mock:           mock.ProductRepositoryOk{},
			errExpected:    true,
			wantErrContain: "the product has no name",
		},
		{
			name: "error-from-repository",
			input: &domain.Product{
				Name:         "Coca-Cola",
				Observations: "",
				Price:        3,
			},
			mock:           mock.ProductRepositoryError{},
			errExpected:    true,
			wantErrContain: "mock error",
		},
	}

	for _, tc := range tt {
		err := service.NewStoreService(tc.mock).Add(tc.input)
		if (err != nil) != tc.errExpected {
			t.Fatalf("%s: unexpected error value %v", tc.name, err)
		}

		if err != nil && !strings.Contains(err.Error(), tc.wantErrContain) {
			t.Fatalf("want error string %q to contain %q", err.Error(), tc.wantErrContain)
		}
	}
}

func TestProductList(t *testing.T) {
	want := domain.ProductList{
		{
			ID:           1,
			Name:         "Coca-Cola",
			Observations: "",
			Price:        3,
		},
		{
			ID:           2,
			Name:         "Big-Cola",
			Observations: "Made in Venezuela",
			Price:        2,
		},
	}

	errExpected := false
	products, err := service.NewStoreService(mock.ProductRepositoryOk{}).List()
	if (err != nil) != errExpected {
		t.Fatalf("unexpected error value %v", err)
	}

	got := make(domain.ProductList, 0, len(products))

	assemble := func(p *domain.Product) domain.ProductCardDTO {
		return domain.ProductCardDTO{
			ID:           p.ID,
			Name:         p.Name,
			Observations: p.Observations,
			Price:        p.Price,
		}
	}

	for _, v := range products {
		got = append(got, assemble(v))
	}

	// Only for testing purposes, you may want to use reflect.DeepEqual.
	// It compares two elements of any type recursively.
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}
