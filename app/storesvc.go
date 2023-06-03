package app

import (
	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/postgres"
)

type storeService struct {
	repo postgres.Product
}

func (s storeService) Add(product *domain.Product) error {
	err := addProduct(product)
	if err != nil {
		return err
	}

	return s.repo.Create(product)
}

// addProduct application logic for adding products to the store.
// The application logic has been split into a smaller function for unit testing
// purposes, and it should do so for the other methods of the Service.
func addProduct(p *domain.Product) error {
	if !p.HasName() {
		return domain.ErrProductHasNoName
	}

	return nil
}

func (s storeService) Find(id int64) (*domain.Product, error) {
	if id == 0 {
		return &domain.Product{}, domain.ErrProductNotFound
	}

	return s.repo.ByID(id)
}

func (s storeService) Update(product domain.Product) error {
	if !product.HasName() {
		return domain.ErrProductHasNoName
	}

	return s.repo.Update(product)
}

func (s storeService) List() ([]*domain.Product, error) {
	return s.repo.All()
}

func (s storeService) Remove(id int64) error {
	if id == 0 {
		return domain.ErrProductNotFound
	}

	return s.repo.Delete(id)
}
