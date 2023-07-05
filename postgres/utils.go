package postgres

import (
	"database/sql"
	"fmt"
	"time"
)

// limitOffset returns a SQL string for a given limit & offset.
func limitOffset(limit, page int) string {
	if limit == 0 && page == 0 {
		return ""
	}

	offset := page*limit - limit
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}

// timeToNull helper to try empty time fields.
func timeToNull(t time.Time) sql.NullTime {
	null := sql.NullTime{Time: t}

	if !null.Time.IsZero() {
		null.Valid = true
	}
	return null
}
