package store

import "testing"

func TestProduct(t *testing.T) {
	tt := []struct {
		name        string
		model       Product
		errExpected bool
	}{
		{
			name:        "empty-model-test",
			model:       Product{},
			errExpected: true,
		},
		{
			name:        "empty-fields-test",
			model:       Product{Name: "", Observations: "", Price: 0.0},
			errExpected: true,
		},
		{
			name: "filled-fields-test",
			model: Product{
				Name:         "Protein",
				Observations: "Lorem ipsum",
				Price:        03333},
			errExpected: false,
		},
	}
	for _, tc := range tt {
		err := tc.model.Validate()
		errReceived := err != nil
		if errReceived != tc.errExpected {
			t.Fatalf("%s: unexpected error value, %v", tc.name, err)
		}
	}
}
