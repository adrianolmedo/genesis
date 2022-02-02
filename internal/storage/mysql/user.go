package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/adrianolmedo/go-restapi/internal/domain"
)

// UserRepository (before UserDAO) it's implementation of UserDAO interface of service/.
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) Create(user *domain.User) error {
	query := "INSERT INTO users (uuid, first_name, last_name, email, password, created_at) VALUES(?, ?, ?, ?, ?, ?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	user.UUID = domain.NextUserID()
	user.CreatedAt = time.Now()

	result, err := stmt.Exec(user.UUID, user.FirstName, user.LastName, user.Email, user.Password, user.CreatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id // Important!
	return nil
}

func (r UserRepository) ByID(id int64) (*domain.User, error) {
	stmt, err := r.db.Prepare("SELECT * FROM users WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// As QueryRow returns a rows we can pass it directly to the mapping
	user, err := scanRowUser(stmt.QueryRow(id))
	if errors.Is(err, sql.ErrNoRows) {
		return &domain.User{}, domain.ErrUserNotFound
	}

	if err != nil {
		return &domain.User{}, err
	}

	return user, nil
}

func (r UserRepository) Update(user domain.User) error {
	stmt, err := r.db.Prepare("UPDATE users SET first_name = ?, last_name = ?, email = ?, password = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	user.UpdatedAt = time.Now()

	result, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Password, timeToNull(user.UpdatedAt), user.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r UserRepository) All() ([]*domain.User, error) {
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

	users := make([]*domain.User, 0)
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

	return users, nil
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
		return domain.ErrUserNotFound
	}
	return nil
}

func (r UserRepository) DeleteAll() error {
	stmt, err := r.db.Prepare("TRUNCATE TABLE users")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("can't truncate table: %v", err)
	}

	return nil
}

// helpers...

type scanner interface {
	Scan(dest ...interface{}) error
}

// scanRowUser return nulled fields of User parsed.
func scanRowUser(s scanner) (*domain.User, error) {
	var updatedAtNull, deletedAtNull sql.NullTime
	user := &domain.User{}

	err := s.Scan(
		&user.ID,
		&user.UUID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&updatedAtNull,
		&deletedAtNull,
	)
	if err != nil {
		return &domain.User{}, err
	}

	user.UpdatedAt = updatedAtNull.Time
	user.DeletedAt = deletedAtNull.Time

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
