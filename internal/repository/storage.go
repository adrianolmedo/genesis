package repository

import "github.com/adrianolmedo/go-restapi-practice/internal/domain"

// UserRepository to uncouple persistence `repository` package
// data between postgres or mysql.
type UserRepository interface {
	Create(domain.User) error
	ByID(id int64) (*domain.User, error)
	Update(domain.User) error
	All() ([]*domain.User, error)
	Delete(id int64) error
}

type LoginRepository interface {
	UserByLogin(email, password string) error
}
