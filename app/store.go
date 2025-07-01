package app

import (
	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/pgsql"
	storage "github.com/adrianolmedo/genesis/pgsql/pgx"
)

type storeService struct {
	repoProduct  storage.Product
	repoCustomer storage.Customer
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
	err := p.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (s storeService) Find(id int64) (*domain.Product, error) {
	if id == 0 {
		return nil, domain.ErrProductNotFound
	}

	return s.repoProduct.ByID(id)
}

func (s storeService) Update(p domain.Product) error {
	err := p.Validate()
	if err != nil {
		return err
	}

	return s.repoProduct.Update(p)
}

func (s storeService) List() (domain.Products, error) {
	return s.repoProduct.All()
}

func (s storeService) AddCustomer(cx *domain.Customer) error {
	return s.repoCustomer.Create(cx)
}

func (s storeService) ListCustomers(p *pgsql.Pager) (pgsql.PagerResults, error) {
	return s.repoCustomer.All(p)
}

func (s storeService) RemoveCustomer(id int64) error {
	if id == 0 {
		return domain.ErrCustomerNotFound
	}

	return s.repoCustomer.Delete(id)
}

func (s storeService) Remove(id int64) error {
	if id == 0 {
		return domain.ErrProductNotFound
	}

	return s.repoProduct.Delete(id)
}
