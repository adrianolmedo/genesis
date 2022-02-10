package mysql

import (
	"database/sql"
	"errors"

	"github.com/adrianolmedo/go-restapi/internal/domain"
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
	stmt, err := r.db.Prepare("SELECT id FROM users WHERE email = ? AND password = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var id int64

	err = stmt.QueryRow(email, password).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrUserNotFound
	}

	if err != nil {
		return err
	}

	return nil
}
