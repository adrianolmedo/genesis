package billing

import (
	"strings"
	"testing"
)

func TestGenerateInvoice(t *testing.T) {
	tt := []struct {
		name           string
		input          *Invoice
		errExpected    bool
		wantErrContain string
	}{
		{
			name: "successful",
			input: &Invoice{
				Header: &InvoiceHeader{
					ClientID: 1,
				},
				Items: ItemList{
					InvoiceItem{ProductID: 1},
				},
			},
			errExpected:    false,
			wantErrContain: "",
		},
		{
			name: "empty-item-list",
			input: &Invoice{
				Header: &InvoiceHeader{
					ClientID: 1,
				},
				Items: ItemList{},
			},
			errExpected:    true,
			wantErrContain: "item list can't be empty",
		},
		{
			name: "nil-item-list",
			input: &Invoice{
				Header: &InvoiceHeader{
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
