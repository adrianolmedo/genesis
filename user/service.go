package user

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

// Find a User by its ID.
func (s Service) Find(id int64) (*User, error) {
	if id == 0 {
		return &User{}, ErrUserNotFound
	}

	return s.repo.ByID(id)
}
