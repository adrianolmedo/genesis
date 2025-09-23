package pgsql

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// FilterMaxLimit is the default value for the limit in a reasonable range.
const FilterMaxLimit int = 10

// Filter is a struct that encapsulates pagination details.
// It includes the limit of results per page, the current page number,
// the field to sort by, and the direction of sorting (ASC or DESC).
type Filter struct {
	limit     int
	page      int
	sort      string
	direction string
}

// NewFilter set values for a Filter and return it.
// It validates the limit and page number, and normalizes the sort direction.
// If the limit is 0 or exceeds [FilterMaxLimit], it defaults to FilterMaxLimit.
func NewFilter(limit, page int, sort, direction string) (Filter, error) {
	limit, err := validateLimit(limit)
	if err != nil {
		return Filter{}, err
	}
	page, err = validatePage(page)
	if err != nil {
		return Filter{}, err
	}
	return Filter{
		limit:     limit,
		page:      page,
		sort:      sort,
		direction: normalizeDirection(direction),
	}, nil
}

// validatePage Filter helper, ensures page number is valid.
func validatePage(p int) (int, error) {
	if p < 0 {
		return p, errors.New("positive number expected for page")
	}
	if p == 0 {
		p = 1
	}
	return p, nil
}

// validateLimit Filter helper, ensures limit is within a reasonable
// range (by default value check [FilterMaxLimit] const).
func validateLimit(n int) (int, error) {
	if n < 0 {
		return n, errors.New("positive number expected for limit")
	}
	maxLimit := FilterMaxLimit
	if n == 0 || n > maxLimit {
		n = maxLimit
	}
	return n, nil
}

// normalizeDirection Filter helper, ensures the direction is either
// ASC or DESC. By default, it returns ASC.
func normalizeDirection(dir string) string {
	dir = strings.ToUpper(dir)
	validDir := map[string]bool{"ASC": true, "DESC": true}
	if validDir[dir] {
		return dir
	}
	return "ASC"
}

// Limit restrict to subset of results.
func (f Filter) Limit() int { return f.limit }

// Page indicates the page from the client.
func (f Filter) Page() int { return f.page }

// Sort sort results by the value of a field, e.g.: ORDER BY created_at.
func (f Filter) Sort() string { return f.sort }

// Direction to display the results in DESC or ASC order based on the
// Sort value.
func (f Filter) Direction() string { return f.direction }

// OrderBy generates an SQL ORDER BY clause.
func (f Filter) OrderBy() string {
	return fmt.Sprintf(`ORDER BY %q %s`, f.sort, f.direction)
}

// LimitOffset generates an SQL LIMIT OFFSET clause.
func (f Filter) LimitOffset() string { return LimitOffset(f.limit, f.page) }

// LimitOffset returns a SQL string for LIMIT OFFSET a given limit & page.
func LimitOffset(limit, page int) string {
	if limit == 0 && page == 0 {
		return ""
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, Offset(limit, page))
}

// Offset calculates the offset for SQL queries.
func (f Filter) Offset() int { return Offset(f.limit, f.page) }

// Offset calculate offset operation from page and limit.
func Offset(limit, page int) int {
	if page == 0 {
		page = 1
	}
	return page*limit - limit
}

// Paginate calculates pagination metadata for the given totalRows.
//
// ItemFrom and ItemTo are 1-based indexes that indicate the position of the
// returned rows relative to the **sorted result set**, not the underlying
// table order. This means:
//   - When Direction = "ASC", ItemFrom is the first row number on the page
//     and ItemTo is the last row number, counting upward from the beginning.
//   - When Direction = "DESC", ItemFrom and ItemTo are still relative to the
//     descending order of the result set. For example, if totalRows = 50,
//     limit = 10, page = 1, and direction = "DESC", then ItemFrom = 1 and
//     ItemTo = 10 correspond to the first 10 rows in descending order
//     (i.e., rows 50–41 in the raw table).
//
// This ensures consistency for API clients: ItemFrom/ItemTo always describe
// the slice of the paginated result **as seen in the API response order**.
func (f Filter) Paginate(totalRows int64) FilterResult {
	if totalRows == 0 {
		return FilterResult{
			Page:       f.page,
			Limit:      f.limit,
			Sort:       f.sort,
			TotalRows:  0,
			TotalPages: 0,
			ItemFrom:   0,
			ItemTo:     0,
		}
	}
	totalPages := int(math.Ceil(float64(totalRows) / float64(f.limit)))
	var fromRow, toRow int
	if f.direction == "ASC" {
		fromRow = (f.page - 1) * f.limit
		toRow = fromRow + f.limit
		if toRow > int(totalRows) {
			toRow = int(totalRows)
		}
	} else { // DESC case
		toRow = int(totalRows) - (f.page-1)*f.limit
		fromRow = toRow - f.limit
		if fromRow < 0 {
			fromRow = 0
		}
	}
	return FilterResult{
		Page:       f.page,
		Limit:      f.limit,
		Sort:       f.sort,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		ItemFrom:   fromRow + 1, // Convert to 1-based index
		ItemTo:     toRow,
	}
}

// FilterResult contains paginated data.
type FilterResult struct {
	Limit      int    `json:"limit"`
	Page       int    `json:"page"`
	Sort       string `json:"sort"`
	TotalRows  int64  `json:"total"`
	TotalPages int    `json:"totalPages"`

	// ItemFrom is the 1-based index of the first item on this page,
	// relative to the sorted result set.
	ItemFrom int `json:"itemFrom"`

	// ItemTo is the 1-based index of the last item on this page,
	// relative to the sorted result set.
	ItemTo int `json:"itemTo"`
}

// Links generates HATEOAS pagination links.
func (f Filter) Links(path string, totalPages int) FilterLinks {
	genLink := func(page int) string {
		return fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", path, f.limit, page, f.sort)
	}
	firstPage := genLink(1)
	lastPage := genLink(totalPages)

	var previousPage, nextPage string
	if f.page > 1 {
		previousPage = genLink(f.page - 1)
	}
	if f.page < totalPages {
		nextPage = genLink(f.page + 1)
	}
	return FilterLinks{
		FirstPage:    firstPage,
		PreviousPage: previousPage,
		NextPage:     nextPage,
		LastPage:     lastPage,
	}
}

// FilterLinks follows HATEOAS principles.
type FilterLinks struct {
	FirstPage    string `json:"first"`
	PreviousPage string `json:"prev"`
	NextPage     string `json:"next"`
	LastPage     string `json:"last"`
}
