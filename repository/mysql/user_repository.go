package mysql

import (
	"database/sql"
	"time"

	"go-restapi-practice/user"
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

func (u UserRepository) Create(req *user.Request) error {
	query := "INSERT INTO users (first_name, last_name, email, password, created_at) VALUES(?, ?, ?, ?, ?)"
	stmt, err := u.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// TO-DO: Cambiar el nombre de la variable m.
	m := User{
		CreatedAt: time.Now(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	}

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

func (u UserRepository) ByID(id int64) (*user.Response, error) {
	stmt, err := u.db.Prepare("SELECT * FROM users WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// As QueryRow returns a rows we can pass it directly to the mapping
	userDB, err := scanRowUser(stmt.QueryRow(id))
	if err != nil {
		return nil, user.ErrResourceDoesNotExist
	}

	return &user.Response{
		ID:        userDB.ID,
		FirstName: userDB.FirstName,
		LastName:  userDB.LastName,
		Email:     userDB.Email,
	}, nil
}

func (u UserRepository) Update(req *user.Request) error {
	stmt, err := u.db.Prepare("UPDATE users SET first_name = ?, last_name = ?, email = ?, password = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	m := User{
		UpdatedAt: time.Now(),
		ID:        req.ID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	}

	result, err := stmt.Exec(m.FirstName, m.LastName, m.Email, m.Password, timeToNull(m.UpdatedAt), m.ID)
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

func (u UserRepository) All() ([]*user.Response, error) {
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

	users := make([]*User, 0)

	for rows.Next() {
		user, err := scanRowUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	resp := make([]*user.Response, 0)

	for i := 0; i < len(users); i++ {
		resp[i].ID = users[i].ID
	}

	return resp, nil
}

func (u UserRepository) Delete(id int64) error {
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
		return user.ErrResourceDoesNotExist
	}
	return nil
}

// User databse model.
type User struct {
	CreatedAt time.Time //`json:"created_at"`
	UpdatedAt time.Time //`json:"updated_at"`
	DeletedAt time.Time //`json:"deleted_at"`
	ID        int64     //`json:"id"`
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
	user := &User{}

	err := s.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&updatedAtNull,
		&deletedAtNull,
	)
	if err != nil {
		return &User{}, err
	}

	user.UpdatedAt = updatedAtNull.Time
	user.DeletedAt = updatedAtNull.Time

	return user, nil
}

// timeToNull helper to try empty time fields.
func timeToNull(t time.Time) sql.NullTime {
	null := sql.NullTime{Time: t}

	if !null.Time.IsZero() {
		null.Valid = true
	}
	return null
}
