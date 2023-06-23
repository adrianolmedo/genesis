package genesis

import (
	"errors"
	"time"
)

var ErrProductNotFound = errors.New("product not found")

// Product model.
type Product struct {
	ID           int64
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

// AddProductForm represents a subset of fields to create a Product.
type AddProductForm struct {
	Name         string  `json:"name"`
	Observations string  `json:"observations"`
	Price        float64 `json:"price"`
}

// UpdateProductForm represents a subset of fields to update a Product.
type UpdateProductForm struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Observations string  `json:"observations"`
	Price        float64 `json:"price"`
}

// ProductCardDTO subset of Product fields.
type ProductCardDTO struct {
	ID           int64   `json:"id,omitempty"`
	Name         string  `json:"name"`
	Observations string  `json:"observations"`
	Price        float64 `json:"price"`
}

// ProductList collection of Product presentation.
type ProductList []ProductCardDTO
