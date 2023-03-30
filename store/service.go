package store

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) Add(product *Product) error {
	if !product.HasName() {
		return ErrProductHasNoName
	}

	return s.repo.Create(product)
}
