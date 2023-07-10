package genesis

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// Direction enum for the sort field of Filter.
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

// Filter query for filtering results.
type Filter struct {
	// limit restrict to subset of results.
	limit int

	// page indicates the page from the client.
	page int

	// sort sort results by the value of a field, e.g.: ORDER BY created_at.
	sort string

	// direction to display the results in DESC or ASC order based on the
	// sort value.
	direction Direction
}

// NewFilter return Filter pointer.
func NewFilter() *Filter {
	return &Filter{}
}

// SetPage by default the firts page is 1 not 0.
func (f *Filter) SetPage(p int) error {
	if p < 0 {
		return errors.New("positive number expected for page")
	}

	if p == 0 {
		p = 1
	}

	f.page = p
	return nil
}

// Page get page.
func (f *Filter) Page() int {
	return f.page
}

// SetLimit by default 10 it's the max limit.
func (f *Filter) SetLimit(n int) error {
	if n < 0 {
		return errors.New("positive number expected for limit")
	}

	// 10 it's the max limit
	if n == 0 || n > 10 {
		n = 10
	}

	f.limit = n
	return nil
}

// Limit get limit.
func (f *Filter) Limit() int {
	return f.limit
}

// SetSort you must choose the default value.
func (f *Filter) SetSort(s string) {
	f.sort = s
}

// Sort get sort.
func (f *Filter) Sort() string {
	return f.sort
}

// SetDirection "asc" or "desc", by default will be asc.
func (f *Filter) SetDirection(direction string) {
	var d Direction
	direction = strings.ToLower(direction)

	if direction == "asc" {
		d = ASC
	}

	if direction == "desc" {
		d = DESC
	}

	f.direction = d
}

// Direction get direction.
func (f *Filter) Direction() Direction {
	return f.direction
}

// Paginate return meta data with subset of filterd results.
func (f *Filter) Paginate(rows interface{}, totalRows int) FilteredResults {
	totalPages := int(math.Ceil(float64(totalRows) / float64(f.limit)))
	//totalPages := int(math.Ceil(float64(totalRows)/float64(f.limit))) - 1
	if totalPages < 0 {
		totalPages = 0
	}

	var fromRow, toRow int

	// Set fromRow and toRow on first page.
	if f.page == 0 {
		fromRow = 1
		toRow = f.limit
	} else if f.Page() <= totalPages {
		// Calculate fromRow and toRow.
		fromRow = f.page*f.limit + 1
		toRow = (f.page + 1) * f.limit
	}

	// Or set toRow with totalRows.
	if toRow > totalRows {
		toRow = totalRows
	}

	return FilteredResults{
		Page:       f.page,
		Limit:      f.limit,
		Sort:       f.sort,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		FromRow:    fromRow,
		ToRow:      toRow,
		Rows:       rows,
	}
}

// FilteredResults that provide the filtered results and its data for
// build the pagination.
type FilteredResults struct {
	Limit      int    `json:"limit"`
	Page       int    `json:"page"`
	Sort       string `json:"sort"`
	TotalRows  int    `json:"total"`
	TotalPages int    `json:"total_pages"`
	FromRow    int    `json:"from_row"`
	ToRow      int    `json:"to_row"`

	// Rows subset of results, not all of results.
	Rows interface{} `json:"-"`
}

// GenLinksResp generate links field to JSON reponse.
func (f *Filter) GenLinksResp(path string, totalPages int) LinksResp {
	var firstPage, lastPage, previousPage, nextPage string

	// Set first page and last page for the pagination reponse.
	firstPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", path, f.limit, 1, f.sort)
	lastPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", path, f.limit, totalPages, f.sort)

	// Set previous page pagination response.
	if f.page > 1 {
		previousPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", path, f.limit, f.page-1, f.sort)
	}

	// Set next pagination response.
	if f.page < totalPages {
		nextPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", path, f.limit, f.page+1, f.sort)
	}

	// Reset previous page.
	if f.page > totalPages {
		previousPage = ""
	}

	return LinksResp{
		FirstPage:    firstPage,
		PreviousPage: previousPage,
		NextPage:     nextPage,
		LastPage:     lastPage,
	}
}

// LinksResp to complies with the HATEOAS principle to display the information
// needed to create pagination.
type LinksResp struct {
	FirstPage    string `json:"first"`
	PreviousPage string `json:"prev"`
	NextPage     string `json:"next"`
	LastPage     string `json:"last"`
}
