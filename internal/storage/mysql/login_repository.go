package mysql

import (
	"database/sql"
)

type LoginRepository struct {
	db *sql.DB
}

func NewLoginRepository(db *sql.DB) *LoginRepository {
	return &LoginRepository{
		db: db,
	}
}

func (r LoginRepository) UserByLogin(email, password string) error {
	stmt, err := r.db.Prepare("SELECT email, password FROM users WHERE email = ? AND password = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(email, password).Err()
	if err != nil {
		return err
	}

	return nil
}
