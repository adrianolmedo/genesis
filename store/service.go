package store

import (
	"context"

	"github.com/adrianolmedo/genesis/pgsql"
)

// Service provides methods to manage products and customers in the store.
type Service struct {
	productRepo  *ProductRepo
	customerRepo *CustomerRepo
}

// NewService creates a new store service with the provided repositories.
func NewService(productRepo *ProductRepo, customerRepo *CustomerRepo) *Service {
	return &Service{
		productRepo:  productRepo,
		customerRepo: customerRepo,
	}
}

func (s Service) Add(ctx context.Context, p *Product) error {
	err := addProduct(p)
	if err != nil {
		return err
	}

	return s.productRepo.Create(ctx, p)
}

// addProduct application logic for adding products to the store.
// The application logic has been split into a smaller function for unit testing
// purposes, and it should do so for the other methods of the Service.
func addProduct(p *Product) error {
	err := p.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (s Service) Find(ctx context.Context, id int64) (*Product, error) {
	if id == 0 {
		return nil, ErrProductNotFound
	}

	return s.productRepo.ByID(ctx, id)
}

func (s Service) Update(ctx context.Context, p Product) error {
	err := p.Validate()
	if err != nil {
		return err
	}

	return s.productRepo.Update(ctx, p)
}

func (s Service) List(ctx context.Context) (Products, error) {
	return s.productRepo.All(ctx)
}

func (s Service) AddCustomer(ctx context.Context, cx *Customer) error {
	return s.customerRepo.Create(ctx, cx)
}

func (s Service) ListCustomers(ctx context.Context, p *pgsql.Pager) (pgsql.PagerResult, error) {
	return s.customerRepo.List(ctx, p)
}

func (s Service) RemoveCustomer(ctx context.Context, id int64) error {
	if id == 0 {
		return ErrCustomerNotFound
	}

	return s.customerRepo.Delete(ctx, id)
}

func (s Service) Remove(ctx context.Context, id int64) error {
	if id == 0 {
		return ErrProductNotFound
	}

	return s.productRepo.Delete(ctx, id)
}
