package gorestapi

import (
	"github.com/adrianolmedo/go-restapi/domain"
	"github.com/adrianolmedo/go-restapi/postgres"
)

type StoreService struct {
	repo postgres.Product
}

func (s StoreService) Add(product *domain.Product) error {
	err := addProduct(product)
	if err != nil {
		return err
	}

	return s.repo.Create(product)
}

// addProduct business logic for adding products to the store.
// The business logic has been split into a smaller function for unit testing
// purposes, and it should do so for the other methods of the Service.
func addProduct(p *domain.Product) error {
	if !p.HasName() {
		return domain.ErrProductHasNoName
	}

	return nil
}

func (s StoreService) Find(id int64) (*domain.Product, error) {
	if id == 0 {
		return &domain.Product{}, domain.ErrProductNotFound
	}

	return s.repo.ByID(id)
}

func (s StoreService) Update(product domain.Product) error {
	if !product.HasName() {
		return domain.ErrProductHasNoName
	}

	return s.repo.Update(product)
}

func (s StoreService) List() ([]*domain.Product, error) {
	return s.repo.All()
}

func (s StoreService) Remove(id int64) error {
	if id == 0 {
		return domain.ErrProductNotFound
	}

	return s.repo.Delete(id)
}
