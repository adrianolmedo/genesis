package genesis

// Direction enum.
type Direction int

const (
	// ASC filter results ascending.
	ASC Direction = iota

	// DESC filter results descending.
	DESC
)

func (d Direction) String() string {
	return [2]string{"ASC", "DESC"}[d]
}

// Filter request model for filtering results.
type Filter struct {
	// Limit restrict to subset of results.
	Limit int

	// Page indicates the page from the client.
	Page int

	// Sort sort results by the value of a field, e.g.: ORDER BY created_at.
	Sort string

	// Direction to display the results in DESC or ASC order based on the
	// Sort value.
	Direction Direction

	// MaxLimit of results per pagination.
	MaxLimit int
}

// FilterResults model that provide the filtered results and its data for
// build the pagination.
type FilterResults struct {
	// Rows subset of results, not all of results.
	Rows       interface{}
	TotalRows  int
	TotalPages int
	FromRow    int
	ToRow      int
}

// PaginationLinks it's a DTO to complies with the HATEOAS principle to display
// the information needed to create pagination.
type PaginationLinks struct {
	Limit        int    `json:"limit"`
	Page         int    `json:"page"`
	Sort         string `json:"sort"`
	TotalRows    int    `json:"total"`
	TotalPages   int    `json:"total_pages"`
	FirstPage    string `json:"first_page"`
	PreviousPage string `json:"previous_page"`
	NextPage     string `json:"next_page"`
	LastPage     string `json:"last_page"`
	FromRow      int    `json:"from_row"`
	ToRow        int    `json:"to_row"`
}
