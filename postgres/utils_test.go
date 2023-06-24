package postgres

import (
	"testing"

	domain "github.com/adrianolmedo/genesis"
)

func TestLimitOffset(t *testing.T) {
	tt := []struct {
		name  string
		input domain.Filter
		want  string
	}{
		{
			name:  "empty-pagination",
			input: domain.Filter{},
			want:  "",
		},
		{
			name: "page-1",
			input: domain.Filter{
				Page:     0,
				Limit:    5,
				MaxLimit: 0,
			},
			want: "LIMIT 5 OFFSET 0",
		},
		{
			name: "page-2",
			input: domain.Filter{
				Page:     2,
				Limit:    10,
				MaxLimit: 10,
			},
			want: "LIMIT 10 OFFSET 10",
		},
	}

	for _, tc := range tt {
		if got := limitOffset(tc.input); tc.want != got {
			t.Errorf("%s: want %q, got %q", tc.name, tc.want, got)
		}
	}
}
