package pgsql

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// PagerMaxLimit is the default value for the limit in a reasonable range.
const PagerMaxLimit uint = 10

// Pager for filtering paginated results.
type Pager struct {
	limit     int
	page      int
	sort      string
	direction string
}

// NewPager constructor, ensures valid defaults when creating a Pager.
func NewPager(limit, page int, sort, direction string) (*Pager, error) {
	limit, err := validateLimit(limit)
	if err != nil {
		return nil, err
	}

	page, err = validatePage(page)
	if err != nil {
		return nil, err
	}

	return &Pager{
		limit:     limit,
		page:      page,
		sort:      sort,
		direction: normalizeDirection(direction),
	}, nil
}

// Limit restrict to subset of results.
func (p *Pager) Limit() int {
	return p.limit
}

// Page indicates the page from the client.
func (p *Pager) Page() int {
	return p.page
}

// Sort sort results by the value of a field, e.g.: ORDER BY created_at.
func (p *Pager) Sort() string {
	return p.sort
}

// Direction to display the results in DESC or ASC order based on the
// Sort value.
func (p *Pager) Direction() string {
	return p.direction
}

// validatePage Pager helper, ensures page number is valid.
func validatePage(p int) (int, error) {
	if p < 0 {
		return p, errors.New("positive number expected for page")
	}

	if p == 0 {
		p = 1
	}

	return p, nil
}

// validateLimit Pager helper, ensures limit is within a reasonable
// range (by default value check PagerMaxLimit const).
func validateLimit(n int) (int, error) {
	if n < 0 {
		return n, errors.New("positive number expected for limit")
	}

	// 10 it's the max limit
	if n == 0 || n > int(PagerMaxLimit) {
		n = int(PagerMaxLimit)
	}

	return n, nil
}

// normalizeDirection Pager helper, ensures the direction is either
// ASC or DESC.
func normalizeDirection(dir string) string {
	dir = strings.ToUpper(dir)
	validDirections := map[string]bool{"ASC": true, "DESC": true}

	if validDirections[dir] {
		return dir
	}
	return "ASC"
}

// OrderBy generates an SQL ORDER BY clause.
func (p *Pager) OrderBy() string {
	return fmt.Sprintf(`ORDER BY %q %s`, p.sort, p.direction)
}

// LimitOffset generates an SQL LIMIT OFFSET clause.
func (p *Pager) LimitOffset() string {
	return LimitOffset(p.limit, p.page)
}

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

// Paginate calculates pagination details.
func (p *Pager) Paginate(rows any, totalRows int) PagerResults {
	if totalRows == 0 {
		return PagerResults{
			Page:       p.page,
			Limit:      p.limit,
			Sort:       p.sort,
			TotalRows:  0,
			TotalPages: 0,
			FromRow:    0,
			ToRow:      0,
			Rows:       rows,
		}
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(p.limit)))

	var fromRow, toRow int

	if p.direction == "ASC" {
		fromRow = (p.page - 1) * p.limit
		toRow = fromRow + p.limit
		if toRow > totalRows {
			toRow = totalRows
		}
	} else { // DESC case
		toRow = totalRows - (p.page-1)*p.limit
		fromRow = toRow - p.limit
		if fromRow < 0 {
			fromRow = 0
		}
	}

	return PagerResults{
		Page:       p.page,
		Limit:      p.limit,
		Sort:       p.sort,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		FromRow:    fromRow + 1, // Convert to 1-based index
		ToRow:      toRow,
		Rows:       rows,
	}
}

// PagerResults contains paginated data.
type PagerResults struct {
	Limit      int    `json:"limit"`
	Page       int    `json:"page"`
	Sort       string `json:"sort"`
	TotalRows  int    `json:"total"`
	TotalPages int    `json:"totalPages"`
	FromRow    int    `json:"fromRow"`
	ToRow      int    `json:"toRow"`

	// Rows subset of results, not all of results.
	Rows any `json:"-"`
}

// Links generates HATEOAS pagination links.
func (p *Pager) Links(path string, totalPages int) PagerLinks {
	genLink := func(page int) string {
		return fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", path, p.limit, page, p.sort)
	}

	firstPage := genLink(1)
	lastPage := genLink(totalPages)

	var previousPage, nextPage string
	if p.page > 1 {
		previousPage = genLink(p.page - 1)
	}
	if p.page < totalPages {
		nextPage = genLink(p.page + 1)
	}

	return PagerLinks{
		FirstPage:    firstPage,
		PreviousPage: previousPage,
		NextPage:     nextPage,
		LastPage:     lastPage,
	}
}

// PagerLinks follows HATEOAS principles.
type PagerLinks struct {
	FirstPage    string `json:"first"`
	PreviousPage string `json:"prev"`
	NextPage     string `json:"next"`
	LastPage     string `json:"last"`
}
