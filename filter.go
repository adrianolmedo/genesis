package genesis

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// Direction enum for the Sort field of Filter.
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
	// Limit restrict to subset of results.
	Limit int

	// MaxLimit of results per pagination.
	MaxLimit int

	// Page indicates the page from the client.
	Page int

	// Sort sort results by the value of a field, e.g.: ORDER BY created_at.
	Sort string

	// Direction to display the results in DESC or ASC order based on the
	// Sort value.
	Direction Direction
}

// NewFilter return Filter pointer with maxLimit by default value as param.
// TODO: Functional options as param.
func NewFilter(maxLimit int) *Filter {

	if maxLimit == 0 {
		maxLimit = 10
	}

	return &Filter{
		MaxLimit: maxLimit,
	}
}

func (f *Filter) SetPage(p int) error {
	if p < 0 {
		return errors.New("positive number expected for page")
	}

	if p == 0 {
		p = 1
	}

	f.Page = p
	return nil
}

func (f *Filter) SetLimit(n int) error {
	if n < 0 {
		return errors.New("positive number expected for limit")
	}

	if n == 0 || n > f.MaxLimit {
		n = f.MaxLimit
	}

	f.Limit = n
	return nil
}

func (f *Filter) SetSort(s string) {
	f.Sort = s
}

func (f *Filter) SetDirection(direction string) {
	var d Direction
	direction = strings.ToLower(direction)

	if direction == "asc" {
		d = ASC
	}

	if direction == "desc" {
		d = DESC
	}

	f.Direction = d
}

// Paginate return meta data with subset of filterd results.
func (f *Filter) Paginate(rows interface{}, totalRows int) FilteredResults {
	totalPages := int(math.Ceil(float64(totalRows) / float64(f.Limit)))
	//totalPages := int(math.Ceil(float64(totalRows)/float64(f.Limit))) - 1
	if totalPages < 0 {
		totalPages = 0
	}

	var fromRow, toRow int

	// Set fromRow and toRow on first page.
	if f.Page == 0 {
		fromRow = 1
		toRow = f.Limit
	} else if f.Page <= totalPages {
		// Calculate fromRow and toRow.
		fromRow = f.Page*f.Limit + 1
		toRow = (f.Page + 1) * f.Limit
	}

	// Or set toRow with totalRows.
	if toRow > totalRows {
		toRow = totalRows
	}

	return FilteredResults{
		Page:       f.Page,
		Limit:      f.Limit,
		Sort:       f.Sort,
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
	firstPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", path, f.Limit, 1, f.Sort)
	lastPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", path, f.Limit, totalPages, f.Sort)

	// Set previous page pagination response.
	if f.Page > 1 {
		previousPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", path, f.Limit, f.Page-1, f.Sort)
	}

	// Set next pagination response.
	if f.Page < totalPages {
		nextPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", path, f.Limit, f.Page+1, f.Sort)
	}

	// Reset previous page.
	if f.Page > totalPages {
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
