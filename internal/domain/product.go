package domain

import (
	"errors"
	"time"
)

var ErrProductHasNoName = errors.New("the product has no name")
var ErrProductNotFound = errors.New("product not found")

type Product struct {
	ID           int64
	Name         string
	Observations string
	Price        float64

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p Product) HasName() bool {
	return p.Name != ""
}

type AddProductForm struct {
	Name         string  `json:"name"`
	Observations string  `json:"observations"`
	Price        float64 `json:"price"`
}

type UpdateProductForm struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Observations string  `json:"observations"`
	Price        float64 `json:"price"`
}

type ProductCardDTO struct {
	ID           int64   `json:"id,omitempty"`
	Name         string  `json:"name"`
	Observations string  `json:"observations"`
	Price        float64 `json:"price"`
}

type ProductList []ProductCardDTO
