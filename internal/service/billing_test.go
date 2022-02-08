package service_test

import (
	"strings"
	"testing"

	"github.com/adrianolmedo/go-restapi/internal/domain"
	"github.com/adrianolmedo/go-restapi/internal/mock"
	"github.com/adrianolmedo/go-restapi/internal/service"
	"github.com/adrianolmedo/go-restapi/internal/storage"
)

func TestGenerate(t *testing.T) {
	tt := []struct {
		name           string
		input          *domain.Invoice
		mock           storage.InvoiceRepository
		errExpected    bool
		wantErrContain string
	}{
		{
			name: "successful",
			input: &domain.Invoice{
				Header: &domain.InvoiceHeader{
					ClientID: 1,
				},
				Items: domain.ItemList{
					&domain.InvoiceItem{ProductID: 1},
				},
			},
			mock:           mock.InvoiceRepositoryOk{},
			errExpected:    false,
			wantErrContain: "",
		},
		{
			name: "product-with-id-zero",
			input: &domain.Invoice{
				Header: &domain.InvoiceHeader{
					ClientID: 1,
				},
				Items: domain.ItemList{
					&domain.InvoiceItem{ProductID: 0},
				},
			},
			mock:           mock.InvoiceRepositoryOk{},
			errExpected:    true,
			wantErrContain: "product not found with id 0",
		},
		{
			name: "not-products",
			input: &domain.Invoice{
				Header: &domain.InvoiceHeader{
					ClientID: 1,
				},
				Items: domain.ItemList{
					&domain.InvoiceItem{},
				},
			},
			mock:           mock.InvoiceRepositoryOk{},
			errExpected:    true,
			wantErrContain: "product not found with id 0",
		},
		{
			name: "empty-item-list",
			input: &domain.Invoice{
				Header: &domain.InvoiceHeader{
					ClientID: 1,
				},
				Items: domain.ItemList{},
			},
			mock:           mock.InvoiceRepositoryOk{},
			errExpected:    true,
			wantErrContain: "item list can't be empty",
		},
		{
			name: "nil-item-list",
			input: &domain.Invoice{
				Header: &domain.InvoiceHeader{
					ClientID: 1,
				},
			},
			mock:           mock.InvoiceRepositoryOk{},
			errExpected:    true,
			wantErrContain: "item list can't be empty",
		},
		{
			name: "error-from-repository",
			input: &domain.Invoice{
				Header: &domain.InvoiceHeader{
					ClientID: 1,
				},
				Items: domain.ItemList{
					&domain.InvoiceItem{ProductID: 1},
				},
			},
			mock:           mock.InvoiceRepositoryError{},
			errExpected:    true,
			wantErrContain: "mock error",
		},
	}

	for _, tc := range tt {
		err := service.NewBillingService(tc.mock).Generate(tc.input)
		if (err != nil) != tc.errExpected {
			t.Fatalf("%s: unexpected error value %v", tc.name, err)
		}

		if err != nil && !strings.Contains(err.Error(), tc.wantErrContain) {
			t.Fatalf("want error string %q to contain %q", err.Error(), tc.wantErrContain)
		}
	}
}
