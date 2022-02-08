package mock

import (
	"errors"
	"time"

	"github.com/adrianolmedo/go-restapi/internal/domain"
)

type ProductRepositoryOk struct{}

func (ProductRepositoryOk) Create(*domain.Product) error {
	return nil
}

func (ProductRepositoryOk) ByID(id int64) (*domain.Product, error) {
	if id == 1 {
		return &domain.Product{
			ID:           1,
			Name:         "Coca-Cola",
			Observations: "",
			Price:        3,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Time{},
		}, nil
	}

	if id == 2 {
		return &domain.Product{
			ID:           2,
			Name:         "Big Cola",
			Observations: "Made in Venezuela",
			Price:        2,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Time{},
		}, nil
	}

	return &domain.Product{}, domain.ErrProductNotFound
}

func (ProductRepositoryOk) Update(domain.Product) error {
	return nil
}

func (ProductRepositoryOk) All() ([]*domain.Product, error) {
	products := []*domain.Product{
		{
			ID:           1,
			Name:         "Coca-Cola",
			Observations: "",
			Price:        3,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Time{},
		},
		{
			ID:           2,
			Name:         "Big-Cola",
			Observations: "Made in Venezuela",
			Price:        2,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Time{},
		},
	}
	return products, nil
}

func (ProductRepositoryOk) Delete(int64) error {
	return nil
}

// ---

type ProductRepositoryError struct{}

func (ProductRepositoryError) Create(*domain.Product) error {
	return errors.New("mock error")
}

func (ProductRepositoryError) ByID(int64) (*domain.Product, error) {
	return &domain.Product{}, errors.New("mock error")
}

func (ProductRepositoryError) Update(domain.Product) error {
	return errors.New("mock error")
}

func (ProductRepositoryError) All() ([]*domain.Product, error) {
	products := make([]*domain.Product, 0)
	return products, errors.New("mock error")
}

func (ProductRepositoryError) Delete(int64) error {
	return errors.New("mock error")
}
