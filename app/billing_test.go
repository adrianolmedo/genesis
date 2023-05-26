package app

import (
	"strings"
	"testing"

	domain "github.com/adrianolmedo/go-restapi"
)

func TestGenerate(t *testing.T) {
	tt := []struct {
		name           string
		input          *domain.Invoice
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
			errExpected:    false,
			wantErrContain: "",
		},
		{
			name: "empty-item-list",
			input: &domain.Invoice{
				Header: &domain.InvoiceHeader{
					ClientID: 1,
				},
				Items: domain.ItemList{},
			},
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
			errExpected:    true,
			wantErrContain: "item list can't be empty",
		},
	}

	for _, tc := range tt {
		err := generateInvoice(tc.input)
		if (err != nil) != tc.errExpected {
			t.Fatalf("%s: unexpected error value %v", tc.name, err)
		}

		if err != nil && !strings.Contains(err.Error(), tc.wantErrContain) {
			t.Fatalf("want error string %q to contain %q", err.Error(), tc.wantErrContain)
		}
	}
}
