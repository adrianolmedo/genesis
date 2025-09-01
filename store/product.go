package store

import (
	"errors"

	"github.com/adrianolmedo/genesis"
)

var ErrProductNotFound = errors.New("product not found")

// Product domain model.
type Product struct {
	ID           int64
	UUID         string
	Name         string
	Observations string
	Price        int64

	genesis.AuditFields
}

// Validate check integrity of fields.
func (p Product) Validate() error {
	if p.Name == "" {
		return errors.New("the product has no name")
	}

	return nil
}

// Products collection of Product.
type Products []Product

// IsEmpty return true if is empty.
func (ps Products) IsEmpty() bool {
	return len(ps) == 0
}
