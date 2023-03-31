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
// The business logic has been split into a smaller function for unit testing
// purposes, and it should do so for the other methods of the Service.
func addProductService(p *Product) error {
	if !p.HasName() {
		return ErrProductHasNoName
	}

	return nil
}

func (s Service) Find(id int64) (*Product, error) {
	if id == 0 {
		return &Product{}, ErrProductNotFound
	}

	return s.repo.ByID(id)
}

func (s Service) Update(product Product) error {
	if !product.HasName() {
		return ErrProductHasNoName
	}

	return s.repo.Update(product)
}

func (s Service) List() ([]*Product, error) {
	return s.repo.All()
}

func (s Service) Remove(id int64) error {
	if id == 0 {
		return ErrProductNotFound
	}

	return s.repo.Delete(id)
}
