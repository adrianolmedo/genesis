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

func (u UserDAO) ByID(id int64) (*model.User, error) {
	stmt, err := u.db.Prepare("SELECT * FROM users WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// As QueryRow returns a rows we can pass it directly to the mapping
	user, err := model.ScanRowUser(stmt.QueryRow(id))
	if err != nil {
		return nil, model.ErrResourceDoesNotExist
	}
	return user, nil
}

func (u UserDAO) Update(m *model.User) error {
	stmt, err := u.db.Prepare("UPDATE users SET first_name = ?, last_name = ?, email = ?, password = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(m.FirstName, m.LastName, m.Email, m.Password, model.TimeToNull(m.UpdatedAt), m.ID)
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
	stmt, err := u.db.Prepare("DELETE FROM users WHERE id = ?")
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
