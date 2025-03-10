package pgsql

import "fmt"

// LimitOffset returns a SQL string for LIMIT OFFSET a given limit & page.
func LimitOffset(limit, page int) string {
	if limit == 0 && page == 0 {
		return ""
	}

	if page == 0 {
		page = 1
	}

	offset := page*limit - limit

	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}
