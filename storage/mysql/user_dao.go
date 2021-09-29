package mysql

import (
	"database/sql"

	"go-restapi-practice/model"
)

// UserDAO it's implementation of UserDAO interface of service/.
// It could be named to UserRepository or Repository.
type UserDAO struct {
	db *sql.DB
}

func NewUserDAO(db *sql.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (u UserDAO) Create(m *model.User) error {
	query := "INSERT INTO users (first_name, last_name, email, password, created_at) VALUES(?, ?, ?, ?, ?)"
	stmt, err := u.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(m.FirstName, m.LastName, m.Email, m.Password, m.CreatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	m.ID = id // Important!
	return nil
}

func (UserDAO) ByID(id int64) (*model.User, error) {
	return nil, nil
}

func (UserDAO) Update(m *model.User) error {
	return nil
}

func (s UserDAO) All() (model.Users, error) {
	return model.Users{}, nil
}

func (UserDAO) Delete(id int64) error {
	return nil
}
