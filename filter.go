package genesis

// TODO: Direction enum.

// Filter for filtering results.
type Filter struct {
	// Field are the filter fields for the query.
	Fields string

	// Limit restrict to subset of results.
	Limit int

	// Page indicates the pagination of the results.
	Page int

	// Sort results by the value of a field, eg: ORDER BY created_at.
	Sort string

	// Direction to display the results in ASC or DESC order based on the
	// Sort value.
	Direction string
}
