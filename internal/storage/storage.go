package storage

import (
	"database/sql"
	"fmt"

	"github.com/adrianolmedo/go-restapi-practice/config"
	"github.com/adrianolmedo/go-restapi-practice/internal/domain"
	"github.com/adrianolmedo/go-restapi-practice/internal/storage/mysql"
	"github.com/adrianolmedo/go-restapi-practice/internal/storage/postgres"
)

const (
	mySQL      config.Driver = "mysql"
	postgreSQL config.Driver = "postgres"
)

type Repositories struct {
	UserRepository  UserRepository
	LoginRepository LoginRepository
}

func NewRepositories(dbcfg config.Database) (*Repositories, error) {
	var err error
	var db *sql.DB

	switch dbcfg.Engine {

	case mySQL:
		db, err = mysql.NewStorage(dbcfg)
		if err != nil {
			return nil, fmt.Errorf("mysql: %v", err)
		}

		return &Repositories{
			UserRepository:  mysql.NewUserRepository(db),
			LoginRepository: mysql.NewLoginRepository(db),
		}, nil

	case postgreSQL:
		db, err = postgres.NewStorage(dbcfg)
		if err != nil {
			return nil, fmt.Errorf("postgres: %v", err)
		}

		return &Repositories{
			UserRepository:  postgres.NewUserRepository(db),
			LoginRepository: postgres.NewLoginRepository(db),
		}, nil

	default:
		return nil, fmt.Errorf("driver not implemented: %s", dbcfg.Engine)
	}
}

// UserRepository to uncouple persistence `repository` package
// data between postgres or mysql.
//
// If there is no data persistence, this interface should not exist,
// and instead, the model, domain or domain/value object should be imported
// from `repository` directly as field in the `service` struct.
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
