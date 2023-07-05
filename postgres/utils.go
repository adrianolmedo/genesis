package postgres

import (
	"database/sql"
	"fmt"
	"time"

	domain "github.com/adrianolmedo/genesis"
)

// orderBy if f.Sort is empty, the default field is "created_at".
func orderBy(f *domain.Filter) string {
	//if f.Sort == "" {
	//	f.Sort = "created_at"
	//}

	return fmt.Sprintf("ORDER BY %s %s", f.Sort, f.Direction)
}

// groupBy if f.Sort is empty, the default field is "created_at".
func groupBy(f *domain.Filter) string {
	//if f.Sort == "" {
	//	f.Sort = "created_at"
	//}

	return fmt.Sprintf("GROUP BY %s", f.Sort)
}

// limitOffset returns a SQL string for a given limit & offset. If the MaxLimit
// is 0 by default it will set to 10.
func limitOffset(f *domain.Filter) string {
	if f.Limit == 0 && f.Page == 0 {
		return ""
	}

	//if f.MaxLimit == 0 {
	//	f.MaxLimit = 10
	//}

	//if f.Limit == 0 || f.Limit > f.MaxLimit {
	//	f.Limit = f.MaxLimit
	//}

	//if f.Page == 0 {
	//	f.Page = 1
	//}

	offset := f.Page*f.Limit - f.Limit
	//offset := f.Page * f.Limit
	return fmt.Sprintf("LIMIT %d OFFSET %d", f.Limit, offset)
}

// timeToNull helper to try empty time fields.
func timeToNull(t time.Time) sql.NullTime {
	null := sql.NullTime{Time: t}

	if !null.Time.IsZero() {
		null.Valid = true
	}
	return null
}
