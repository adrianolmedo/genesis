package rest

import "github.com/adrianolmedo/genesis/pgsql"

// resp represents a standard successful API response.
type resp struct {
	Status string `json:"status"`
	detailsResp
}

// errorResp represents a standard error API response.
type errorResp struct {
	Status string      `json:"status"`
	Error  detailsResp `json:"error"`
}

// detailsResp holds the details of the response, including code, data and message.
type detailsResp struct {
	Code    string `json:"code,omitempty"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
}

// pagerResp represents a paginated API response using existing DTOs.
type pagerResp struct {
	Meta  pgsql.PagerResult `json:"meta"`  // Uses your existing struct for metadata
	Data  any               `json:"data"`  // Holds the actual paginated data
	Links pgsql.PagerLinks  `json:"links"` // Uses your existing struct for pagination links
}
