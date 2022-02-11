package storage

import (
	"database/sql"
	"fmt"

	"github.com/adrianolmedo/go-restapi/config"
	"github.com/adrianolmedo/go-restapi/internal/domain"
	"github.com/adrianolmedo/go-restapi/internal/storage/mysql"
	"github.com/adrianolmedo/go-restapi/internal/storage/postgres"
)

type Storage interface {
	ProvideRepository() (*Repository, error)
}

type storage struct {
	dbcfg config.Database
}

func New(dbcfg config.Database) *storage {
	return &storage{dbcfg}
}

func (s storage) ProvideRepository() (*Repository, error) {
	var err error
	var db *sql.DB

	switch s.dbcfg.Engine {

	case "mysql":
		db, err = mysql.New(s.dbcfg)
		if err != nil {
			return nil, fmt.Errorf("mysql: %v", err)
		}

		return &Repository{
			User:    mysql.NewUserRepository(db),
			Login:   mysql.NewLoginRepository(db),
			Product: mysql.NewProductRepository(db),
			Invoice: mysql.NewInvoiceRepository(db,
				mysql.NewInvoiceHeaderRepository(db),
				mysql.NewInvoiceItemRepository(db)),
		}, nil

	case "postgres":
		db, err = postgres.New(s.dbcfg)
		if err != nil {
			return nil, fmt.Errorf("postgres: %v", err)
		}

		return &Repository{
			User:    postgres.NewUserRepository(db),
			Login:   postgres.NewLoginRepository(db),
			Product: postgres.NewProductRepository(db),
			Invoice: postgres.NewInvoiceRepository(db,
				postgres.NewInvoiceHeaderRepository(db),
				postgres.NewInvoiceItemRepository(db)),
		}, nil

	default:
		return nil, fmt.Errorf("database engine '%s' not implemented", s.dbcfg.Engine)
	}
}

type Repository struct {
	User    UserRepository
	Login   LoginRepository
	Product ProductRepository
	Invoice InvoiceRepository
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

type ProductRepository interface {
	Create(*domain.Product) error
	ByID(id int64) (*domain.Product, error)
	Update(domain.Product) error
	All() ([]*domain.Product, error)
	Delete(id int64) error
}

type InvoiceHeaderRepository interface {
	CreateTx(*sql.Tx, *domain.InvoiceHeader) error
}

type InvoiceItemRepository interface {
	CreateTx(*sql.Tx, int64, domain.ItemList) error
}

type InvoiceRepository interface {
	Create(m *domain.Invoice) error
}
