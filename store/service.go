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
	err := addProductService(product)
	if err != nil {
		return err
	}

	return s.repo.Create(product)
}

// addProductService business logic for adding products to the store.
func addProductService(p *Product) error {
	if !p.HasName() {
		return ErrProductHasNoName
	}

	return nil
}
