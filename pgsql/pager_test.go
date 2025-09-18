package pgsql_test

import (
	"testing"

	"github.com/adrianolmedo/genesis/pgsql"
)

func TestPager(t *testing.T) {
	tt := []struct {
		name        string // test name
		limit       int
		page        int
		sort        string
		direction   string
		errExpected bool
	}{
		{
			name:        "page-zero",
			limit:       2,
			page:        0,
			sort:        "created_at",
			direction:   "ASC",
			errExpected: false,
		},
		{
			name:        "negative-limit",
			limit:       -2,
			page:        0,
			sort:        "created_at",
			direction:   "ASC",
			errExpected: true,
		},
		{
			name:        "negative-page",
			limit:       2,
			page:        -1,
			sort:        "created_at",
			direction:   "ASC",
			errExpected: true,
		},
	}
	for _, tc := range tt {
		_, err := pgsql.NewPager(tc.limit, tc.page, tc.sort, tc.direction)
		errReceived := err != nil
		if errReceived != tc.errExpected {
			t.Fatalf("%s: NewPager(%d, %d, %q, %q): unexpected error status: %v",
				tc.name, tc.limit, tc.page, tc.sort, tc.direction, err)
		}
	}
}

func TestLimitOffset(t *testing.T) {
	tt := []struct {
		name  string // test name
		page  int
		limit int
		want  string
	}{
		{
			name: "empty-pagination",
			want: "",
		},
		{
			name:  "page-one",
			page:  0,
			limit: 5,
			want:  "LIMIT 5 OFFSET 0",
		},
		{
			name:  "page-two",
			page:  2,
			limit: 10,
			want:  "LIMIT 10 OFFSET 10",
		},
	}
	for _, tc := range tt {
		got := pgsql.LimitOffset(tc.limit, tc.page)
		if tc.want != got {
			t.Errorf("%s: got %q, want %q", tc.name, got, tc.want)
		}
	}
}
