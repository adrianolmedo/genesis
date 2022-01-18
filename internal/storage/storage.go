package storage

import (
	"database/sql"
	"fmt"

	"github.com/adrianolmedo/go-restapi-practice/config"
	"github.com/adrianolmedo/go-restapi-practice/internal/domain"
	"github.com/adrianolmedo/go-restapi-practice/internal/storage/mysql"
	"github.com/adrianolmedo/go-restapi-practice/internal/storage/postgres"
)

type RepoProvider interface {
	ProvideRepo() (*Repository, error)
}

type Storage struct {
	dbcfg config.Database
}

func New(dbcfg config.Database) *Storage {
	return &Storage{dbcfg}
}

func (s Storage) ProvideRepo() (*Repository, error) {
	var err error
	var db *sql.DB

	switch s.dbcfg.Engine {

	case "mysql":
		db, err = mysql.New(s.dbcfg)
		if err != nil {
			return nil, fmt.Errorf("mysql: %v", err)
		}

		return &Repository{
			User:  mysql.NewUserRepository(db),
			Login: mysql.NewLoginRepository(db),
		}, nil

	case "postgres":
		db, err = postgres.New(s.dbcfg)
		if err != nil {
			return nil, fmt.Errorf("postgres: %v", err)
		}

		return &Repository{
			User:  postgres.NewUserRepository(db),
			Login: postgres.NewLoginRepository(db),
		}, nil

	default:
		return nil, fmt.Errorf("driver not implemented: %s", s.dbcfg.Engine)
	}
}

type Repository struct {
	User  UserRepository
	Login LoginRepository
}

// UserRepository to uncouple persistence `repository` package
// data between postgres or mysql.
type UserRepository interface {
	Create(*domain.User) error
	ByID(id int64) (*domain.User, error)
	Update(domain.User) error
	All() ([]*domain.User, error)
	Delete(id int64) error
}

type LoginRepository interface {
	UserByLogin(email, password string) error
}
