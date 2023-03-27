package user

import "database/sql"

type Service struct {
	repo *repository
}

func NewService(db *sql.DB) *Service {
	return &Service{
		repo: newRepository(db),
	}
}

// Find a User by its ID.
func (s Service) Find(id int64) (*User, error) {
	if id == 0 {
		return &User{}, ErrUserNotFound
	}

	return s.repo.ByID(id)
}
