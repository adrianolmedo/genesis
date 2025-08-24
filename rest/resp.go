package rest

import "github.com/adrianolmedo/genesis/pgsql"

// resp represents a standard successful API response.
type resp struct {
	Status string `json:"status"`
	respDetails
}

// respDetails holds the details of the response, including code, data, message, and optional details.
type respDetails struct {
	Code    string `json:"code,omitempty"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// errorResp represents a standard error API response.
type errorResp struct {
	Status string      `json:"status"`
	Error  respDetails `json:"error"`
}

// pagerResp represents a paginated API response using existing DTOs.
type pagerResp struct {
	Meta  pgsql.PagerResults `json:"meta"`  // Uses your existing struct for metadata
	Data  any                `json:"data"`  // Holds the actual paginated data
	Links pgsql.PagerLinks   `json:"links"` // Uses your existing struct for pagination links
}
