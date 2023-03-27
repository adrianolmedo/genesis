package product

import "database/sql"

type Service struct {
	repo *repository
}

func NewService(db *sql.DB) *Service {
	return &Service{
		repo: newRepository(db),
	}
}

func (s Service) Add(product *Product) error {
	if !product.HasName() {
		return ErrProductHasNoName
	}

	return s.repo.Create(product)
}
