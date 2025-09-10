package store

import (
	"strings"
	"testing"
)

func TestAddProduct(t *testing.T) {
	tt := []struct {
		name           string
		input          *Product
		errExpected    bool
		wantErrContain string
	}{
		{
			name: "successful-test", // test name
			input: &Product{
				Name:         "Coca-Cola",
				Observations: "",
				Price:        3,
			},
			errExpected:    false,
			wantErrContain: "",
		},
		{
			name: "empty-field-test",
			input: &Product{
				Name:         "",
				Observations: "Made in Venezuela",
				Price:        2,
			},
			errExpected:    true,
			wantErrContain: "the product has no name",
		},
	}
	for _, tc := range tt {
		err := addProduct(tc.input)
		if (err != nil) != tc.errExpected {
			t.Fatalf("%s: unexpected error value, %v", tc.name, err)
		}
		if err != nil && !strings.Contains(err.Error(), tc.wantErrContain) {
			t.Fatalf("want error string %q to contain %q", err.Error(), tc.wantErrContain)
		}
	}
}
