package service

import (
	"database/sql"
	"fmt"

	"github.com/adrianolmedo/go-restapi-practice/config"
	"github.com/adrianolmedo/go-restapi-practice/internal/repository"
	"github.com/adrianolmedo/go-restapi-practice/internal/repository/mysql"
	"github.com/adrianolmedo/go-restapi-practice/internal/repository/postgres"
)

const (
	mySQL      = "mysql"
	postgreSQL = "postgres"
)

type Service struct {
	UserService  UserService
	LoginService LoginService
}

func New(dbcfg config.Database) (*Service, error) {
	repos, err := newStorage(dbcfg)
	if err != nil {
		return nil, fmt.Errorf("error from repository: %v", err)
	}

	return &Service{
		UserService:  NewUserService(repos.userRepository),
		LoginService: NewLoginService(repos.loginRepository),
	}, nil
}

// storage is a collection of all repositories interfaces.
type storage struct {
	userRepository  repository.UserRepository
	loginRepository repository.LoginRepository
}

func newStorage(dbcfg config.Database) (*storage, error) {
	var err error
	var db *sql.DB

	switch dbcfg.Engine {

	case mySQL:
		db, err = mysql.New(dbcfg)
		if err != nil {
			return nil, fmt.Errorf("mysql: %v", err)
		}

		return &storage{
			userRepository:  mysql.NewUserRepository(db),
			loginRepository: mysql.NewLoginRepository(db),
		}, nil

	case postgreSQL:
		db, err = postgres.New(dbcfg)
		if err != nil {
			return nil, fmt.Errorf("postgres: %v", err)
		}

		return &storage{
			userRepository:  postgres.NewUserRepository(db),
			loginRepository: postgres.NewLoginRepository(db),
		}, nil

	default:
		return nil, fmt.Errorf("driver not implemented: %s", dbcfg.Engine)
	}
}
