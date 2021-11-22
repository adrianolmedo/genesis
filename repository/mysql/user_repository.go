package mysql

import (
	"database/sql"
	"time"

	"github.com/adrianolmedo/go-restapi-practice/user"
)

// UserRepository (before UserDAO) it's implementation of UserDAO interface of service/.
// It could be named to UserRepository or Repository.
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) Create(model *user.User) error {
	query := "INSERT INTO users (first_name, last_name, email, password, created_at) VALUES(?, ?, ?, ?, ?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	u := User{
		CreatedAt: time.Now(),
		FirstName: model.FirstName,
		LastName:  model.LastName,
		Email:     model.Email,
		Password:  model.Password,
	}

	result, err := stmt.Exec(u.FirstName, u.LastName, u.Email, u.Password, u.CreatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = id // Important!
	return nil
}

func (r UserRepository) ByID(id int64) (*user.User, error) {
	stmt, err := r.db.Prepare("SELECT * FROM users WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// As QueryRow returns a rows we can pass it directly to the mapping
	u, err := scanRowUser(stmt.QueryRow(id))
	if err != nil {
		return nil, user.ErrResourceDoesNotExist
	}

	return &user.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}, nil
}

func (r UserRepository) Update(model *user.User) error {
	stmt, err := r.db.Prepare("UPDATE users SET first_name = ?, last_name = ?, email = ?, password = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	u := User{
		UpdatedAt: time.Now(),
		ID:        model.ID,
		FirstName: model.FirstName,
		LastName:  model.LastName,
		Email:     model.Email,
		Password:  model.Password,
	}

	result, err := stmt.Exec(u.FirstName, u.LastName, u.Email, u.Password, timeToNull(u.UpdatedAt), u.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return user.ErrResourceDoesNotExist
	}
	return nil
}

func (r UserRepository) All() ([]*user.User, error) {
	stmt, err := r.db.Prepare("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		u, err := scanRowUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	resp := make([]*user.User, 0, len(users))

	assemble := func(u *User) *user.User {
		return &user.User{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Password:  u.Password,
		}
	}

	for _, u := range users {
		resp = append(resp, assemble(u))
	}

	return resp, nil
}

func (r UserRepository) Delete(id int64) error {
	stmt, err := r.db.Prepare("DELETE FROM users WHERE id = ?")
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
		return user.ErrResourceDoesNotExist
	}
	return nil
}

// User database model.
type User struct {
	CreatedAt time.Time //`json:"created_at"`
	UpdatedAt time.Time //`json:"updated_at"`
	DeletedAt time.Time //`json:"deleted_at"`
	ID        int64     //`json:"id"`
	UUID      string    //`json:"uuid"`
	FirstName string    //`json:"first_name"`
	LastName  string    //`json:"last_name"`
	Email     string    //`json:"email"`
	Password  string    //`json:"password,omitempty"`
}

// helpers...

type scanner interface {
	Scan(dest ...interface{}) error
}

// scanRowUser return nulled fields of User parsed.
func scanRowUser(s scanner) (*User, error) {
	var updatedAtNull, deletedAtNull sql.NullTime
	u := &User{}

	err := s.Scan(
		&u.ID,
		&u.UUID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&updatedAtNull,
		&deletedAtNull,
	)
	if err != nil {
		return &User{}, err
	}

	u.UpdatedAt = updatedAtNull.Time
	u.DeletedAt = updatedAtNull.Time

	return u, nil
}

// timeToNull helper to try empty time fields.
func timeToNull(t time.Time) sql.NullTime {
	null := sql.NullTime{Time: t}

	if !null.Time.IsZero() {
		null.Valid = true
	}
	return null
}
