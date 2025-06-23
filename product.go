package genesis

import (
	"errors"
	"time"
)

var ErrProductNotFound = errors.New("product not found")

// Product model.
type Product struct {
	ID           int
	UUID         string
	Name         string
	Observations string
	Price        float64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// Validate check integrity of fields.
func (p Product) Validate() error {
	if p.Name == "" {
		return errors.New("the product has no name")
	}

	return nil
}

// Products collection of Product.
type Products []*Product

// IsEmpty return true if is empty.
func (ps Products) IsEmpty() bool {
	return len(ps) == 0
}
