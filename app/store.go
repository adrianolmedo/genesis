package app

import (
	"context"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/pgsql"
	storage "github.com/adrianolmedo/genesis/pgsql/sqlc"
)

type storeService struct {
	repoProduct  *storage.Product
	repoCustomer *storage.Customer
}

func (s storeService) Add(ctx context.Context, p *domain.Product) error {
	err := addProduct(p)
	if err != nil {
		return err
	}

	return s.repoProduct.Create(ctx, p)
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

func (s storeService) Find(ctx context.Context, id int64) (*domain.Product, error) {
	if id == 0 {
		return nil, domain.ErrProductNotFound
	}

	return s.repoProduct.ByID(ctx, id)
}

func (s storeService) Update(ctx context.Context, p domain.Product) error {
	err := p.Validate()
	if err != nil {
		return err
	}

	return s.repoProduct.Update(ctx, p)
}

func (s storeService) List(ctx context.Context) (domain.Products, error) {
	return s.repoProduct.All(ctx)
}

func (s storeService) AddCustomer(ctx context.Context, cx *domain.Customer) error {
	return s.repoCustomer.Create(ctx, cx)
}

func (s storeService) ListCustomers(ctx context.Context, p *pgsql.Pager) (pgsql.PagerResults, error) {
	return s.repoCustomer.List(ctx, p)
}

func (s storeService) RemoveCustomer(ctx context.Context, id int64) error {
	if id == 0 {
		return domain.ErrCustomerNotFound
	}

	return s.repoCustomer.Delete(ctx, id)
}

func (s storeService) Remove(ctx context.Context, id int64) error {
	if id == 0 {
		return domain.ErrProductNotFound
	}

	return s.repoProduct.Delete(ctx, id)
}
