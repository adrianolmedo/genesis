package service

import (
	"github.com/adrianolmedo/genesis/internal/domain"
	"github.com/adrianolmedo/genesis/internal/storage"
)

type StoreService interface {
	Add(*domain.Product) error
	Find(id int64) (*domain.Product, error)
	Update(domain.Product) error
	List() ([]*domain.Product, error)
	Remove(id int64) error
}

type storeService struct {
	repo storage.ProductRepository
}

func NewStoreService(repo storage.ProductRepository) *storeService {
	return &storeService{repo}
}

func (ss storeService) Add(product *domain.Product) error {
	if !product.HasName() {
		return domain.ErrProductHasNoName
	}

	return ss.repo.Create(product)
}

func (ss storeService) Find(id int64) (*domain.Product, error) {
	if id == 0 {
		return &domain.Product{}, domain.ErrProductNotFound
	}

	return ss.repo.ByID(id)
}

func (ss storeService) Update(product domain.Product) error {
	if !product.HasName() {
		return domain.ErrProductHasNoName
	}

	return ss.repo.Update(product)
}

func (ss storeService) List() ([]*domain.Product, error) {
	return ss.repo.All()
}

func (ss storeService) Remove(id int64) error {
	if id == 0 {
		return domain.ErrProductNotFound
	}

	return ss.repo.Delete(id)
}
