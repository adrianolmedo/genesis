package store

import (
	"strings"
	"testing"

	"github.com/adrianolmedo/go-restapi/domain"
)

func TestAddProductService(t *testing.T) {
	tt := []struct {
		name           string
		input          *domain.Product
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
			errExpected:    true,
			wantErrContain: "the product has no name",
		},
	}

	for _, tc := range tt {
		err := addProductService(tc.input)
		if (err != nil) != tc.errExpected {
			t.Fatalf("%s: unexpected error value %v", tc.name, err)
		}

		if err != nil && !strings.Contains(err.Error(), tc.wantErrContain) {
			t.Fatalf("want error string %q to contain %q", err.Error(), tc.wantErrContain)
		}
	}
}
