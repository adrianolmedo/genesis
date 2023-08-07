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
}

// Validate check that Name can't be empty.
func (p Product) Validate() error {
	if p.Name != "" {
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

// AddProductForm represents a subset of fields to create a Product.
// TODO: Pass DTO to http/ package.
type AddProductForm struct {
	Name         string  `json:"name"`
	Observations string  `json:"observations"`
	Price        float64 `json:"price"`
}

// UpdateProductForm represents a subset of fields to update a Product.
// TODO: Pass DTO to http/ package.
type UpdateProductForm struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Observations string  `json:"observations"`
	Price        float64 `json:"price"`
}

// ProductCardDTO subset of Product fields.
// TODO: Pass DTO to http/ package.
type ProductCardDTO struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Observations string  `json:"observations"`
	Price        float64 `json:"price"`
}
