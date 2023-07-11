package postgres

import (
	"testing"
)

func TestLimitOffset(t *testing.T) {
	tt := []struct {
		name  string
		page  int
		limit int
		want  string
	}{
		{
			name: "empty-pagination",
			want: "",
		},
		{
			name:  "page-1",
			page:  0,
			limit: 5,
			want:  "LIMIT 5 OFFSET 0",
		},
		{
			name:  "page-2",
			page:  2,
			limit: 10,
			want:  "LIMIT 10 OFFSET 10",
		},
	}

	for _, tc := range tt {
		got := limitOffset(tc.limit, tc.page)
		if tc.want != got {
			t.Errorf("%s: want %q, got %q", tc.name, tc.want, got)
		}
	}
}
