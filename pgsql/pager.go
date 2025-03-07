package pgsql

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// Direction enum for the sort field of Pager.
type Direction int

const (
	// ASC sort results ascending.
	ASC Direction = iota

	// DESC sort results descending.
	DESC
)

func (d Direction) String() string {
	return [2]string{"ASC", "DESC"}[d]
}

// Pager query for filtering paginated results.
type Pager struct {
	limit     int
	page      int
	sort      string
	direction Direction
}

// NewPager constructor.
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
		direction: getDirection(direction),
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
func (p *Pager) Direction() Direction {
	return p.direction
}

// validatePage Pager helper.
func validatePage(p int) (int, error) {
	if p < 0 {
		return p, errors.New("positive number expected for page")
	}

	if p == 0 {
		p = 1
	}

	return p, nil
}

// validateLimit Pager helper, by default 10 it's the max limit.
func validateLimit(n int) (int, error) {
	if n < 0 {
		return n, errors.New("positive number expected for limit")
	}

	// 10 it's the max limit
	if n == 0 || n > 10 {
		n = 10
	}

	return n, nil
}

// getDirection Pager helper.
func getDirection(dir string) Direction {
	var d Direction
	dir = strings.ToLower(dir)

	if dir == "asc" {
		d = ASC
	}

	if dir == "desc" {
		d = DESC
	}

	return d
}

// Paginate return meta data with subset of filterd results.
func (p *Pager) Paginate(rows any, totalRows int) PagerResults {
	totalPages := int(math.Ceil(float64(totalRows) / float64(p.limit)))
	totalPages = int(math.Max(float64(totalPages), 0))

	var fromRow, toRow int

	// Set fromRow and toRow on first page.
	if p.page == 0 {
		fromRow = 1
		toRow = p.limit
	} else if p.page <= totalPages {
		// Calculate fromRow and toRow.
		fromRow = p.page*p.limit + 1
		toRow = (p.page + 1) * p.limit
	}

	// Or set toRow with totalRows.
	if toRow > totalRows {
		toRow = totalRows
	}

	return PagerResults{
		Page:       p.page,
		Limit:      p.limit,
		Sort:       p.sort,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		FromRow:    fromRow,
		ToRow:      toRow,
		Rows:       rows,
	}
}

// PagerResults that provide the filtered results and its data for
// build the pagination.
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

// GenLinks generate links field to JSON reponse.
func (p *Pager) GenLinks(path string, totalPages int) PagerLinks {
	var firstPage, lastPage, previousPage, nextPage string

	// Set first page and last page for the pagination reponse.
	firstPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", path, p.limit, 1, p.sort)
	lastPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", path, p.limit, totalPages, p.sort)

	// Set previous page pagination response.
	if p.page > 1 {
		previousPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", path, p.limit, p.page-1, p.sort)
	}

	// Set next pagination response.
	if p.page < totalPages {
		nextPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", path, p.limit, p.page+1, p.sort)
	}

	// Reset previous page.
	if p.page > totalPages {
		previousPage = ""
	}

	return PagerLinks{
		FirstPage:    firstPage,
		PreviousPage: previousPage,
		NextPage:     nextPage,
		LastPage:     lastPage,
	}
}

// PagerLinks to complies with the HATEOAS principle to display the information
// needed to create pagination.
type PagerLinks struct {
	FirstPage    string `json:"first"`
	PreviousPage string `json:"prev"`
	NextPage     string `json:"next"`
	LastPage     string `json:"last"`
}
