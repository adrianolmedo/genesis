package postgres

import (
	"database/sql"
	"fmt"
	"time"

	domain "github.com/adrianolmedo/genesis"
)

// orderBy if f.Sort is empty, the default field is "created_at".
func orderBy(f domain.Filter) string {
	if f.Sort == "" {
		f.Sort = "created_at"
	}

	// TODO: Direction enum.
	//if f.Direction == "" {
	//	f.Direction = "ASC"
	//}

	return fmt.Sprintf("ORDER BY %s %s", f.Sort, f.Direction)
}

// limitOffset the max limit by default is 10.
func limitOffset(p domain.Filter) string {
	if p.Limit == 0 && p.Page == 0 {
		return ""
	}

	maxLimit := 10

	if p.Limit == 0 || p.Limit > maxLimit {
		p.Limit = maxLimit
	}

	if p.Page == 0 {
		p.Page = 1
	}

	offset := p.Page*p.Limit - p.Limit

	return fmt.Sprintf("LIMIT %d OFFSET %d", p.Limit, offset)
}

// timeToNull helper to try empty time fields.
func timeToNull(t time.Time) sql.NullTime {
	null := sql.NullTime{Time: t}

	if !null.Time.IsZero() {
		null.Valid = true
	}
	return null
}
