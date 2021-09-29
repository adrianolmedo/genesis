package postgres

import (
	"database/sql"
	"time"

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
	stmt, err := u.db.Prepare("INSERT INTO users (first_name, last_name, email, password, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(m.FirstName, m.LastName, m.Email, m.Password, m.CreatedAt).Scan(&m.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u UserDAO) ByID(id int64) (*model.User, error) {
	stmt, err := u.db.Prepare("SELECT * FROM users WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	user, err := model.ScanRowUser(stmt.QueryRow(id))
	if err != nil {
		return nil, model.ErrResourceDoesNotExist
	}
	return user, nil
}

func (u UserDAO) Update(m *model.User) error {
	stmt, err := u.db.Prepare("UPDATE users SET first_name = $1, last_name = $2, email = $3, password = $4, updated_at = $5 WHERE id = $6")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(m.FirstName, m.LastName, m.Email, m.Password, timeToNull(m.UpdatedAt), m.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return model.ErrResourceDoesNotExist
	}
	return nil
}

func (u UserDAO) All() (model.Users, error) {
	stmt, err := u.db.Prepare("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make(model.Users, 0)
	for rows.Next() {
		user, err := model.ScanRowUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u UserDAO) Delete(id int64) error {
	stmt, err := u.db.Prepare("DELETE FROM users WHERE id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return model.ErrResourceDoesNotExist
	}
	return nil
}

// timeToNull helper to try empty time fields.
func timeToNull(t time.Time) sql.NullTime {
	null := sql.NullTime{Time: t}
	if !null.Time.IsZero() {
		null.Valid = true
	}
	return null
}
