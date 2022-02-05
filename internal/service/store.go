package service

import (
	"github.com/adrianolmedo/go-restapi/internal/domain"
	"github.com/adrianolmedo/go-restapi/internal/storage"
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

func NewStoreService(repo storage.ProductRepository) StoreService {
	return &storeService{repo}
}

func (ss storeService) Add(product *domain.Product) error {
	if !product.HasName() {
		return domain.ErrProductHasNoName
	}

	return ss.repo.Create(product)
}

func (ss storeService) Find(id int64) (*domain.Product, error) {
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
	return ss.repo.Delete(id)
}
