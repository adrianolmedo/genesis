package app

import (
	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/postgres"
)

type storeService struct {
	repoProduct  postgres.Product
	repoCustomer postgres.Customer
}

func (s storeService) Add(product *domain.Product) error {
	err := addProduct(product)
	if err != nil {
		return err
	}

	return s.repoProduct.Create(product)
}

// addProduct application logic for adding products to the store.
// The application logic has been split into a smaller function for unit testing
// purposes, and it should do so for the other methods of the Service.
func addProduct(p *domain.Product) error {
	if err := p.Validate(); err != nil {
		return err
	}

	return nil
}

func (s storeService) Find(id int64) (*domain.Product, error) {
	if id == 0 {
		return &domain.Product{}, domain.ErrProductNotFound
	}

	return s.repoProduct.ByID(id)
}

func (s storeService) Update(p domain.Product) error {
	if err := p.Validate(); err != nil {
		return err
	}

	return s.repoProduct.Update(p)
}

func (s storeService) List() ([]*domain.Product, error) {
	return s.repoProduct.All()
}

func (s storeService) AddCustomer(c *domain.Customer) error {
	return s.repoCustomer.Create(c)
}

func (s storeService) ListCustomers(filter domain.Filter) (domain.Customers, error) {
	return s.repoCustomer.All(filter)
}

func (s storeService) Remove(id int64) error {
	if id == 0 {
		return domain.ErrProductNotFound
	}

	return s.repoProduct.Delete(id)
}
