package store

import (
	"github.com/adrianolmedo/go-restapi/domain"
	"github.com/adrianolmedo/go-restapi/postgres"
)

type Service struct {
	repo postgres.Product
}

func NewService(repo postgres.Product) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) Add(product *domain.Product) error {
	err := addProductService(product)
	if err != nil {
		return err
	}

	return s.repo.Create(product)
}

// addProductService business logic for adding products to the store.
// The business logic has been split into a smaller function for unit testing
// purposes, and it should do so for the other methods of the Service.
func addProductService(p *domain.Product) error {
	if !p.HasName() {
		return domain.ErrProductHasNoName
	}

	return nil
}

func (s Service) Find(id int64) (*domain.Product, error) {
	if id == 0 {
		return &domain.Product{}, domain.ErrProductNotFound
	}

	return s.repo.ByID(id)
}

func (s Service) Update(product domain.Product) error {
	if !product.HasName() {
		return domain.ErrProductHasNoName
	}

	return s.repo.Update(product)
}

func (s Service) List() ([]*domain.Product, error) {
	return s.repo.All()
}

func (s Service) Remove(id int64) error {
	if id == 0 {
		return domain.ErrProductNotFound
	}

	return s.repo.Delete(id)
}
