package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/adrianolmedo/go-restapi/domain"
)

type User struct {
	db *sql.DB
}

func (r User) Create(u *domain.User) error {
	stmt, err := r.db.Prepare("INSERT INTO users (uuid, first_name, last_name, email, password, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	u.UUID = domain.NextUserID()
	u.CreatedAt = time.Now()

	err = stmt.QueryRow(u.UUID, u.FirstName, u.LastName, u.Email, u.Password, u.CreatedAt).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r User) ByID(id int64) (*domain.User, error) {
	stmt, err := r.db.Prepare("SELECT * FROM users WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	u, err := scanRowUser(stmt.QueryRow(id))
	if errors.Is(err, sql.ErrNoRows) {
		return &domain.User{}, domain.ErrUserNotFound
	}

	if err != nil {
		return &domain.User{}, err
	}

	return u, nil
}

func (r User) Update(u domain.User) error {
	stmt, err := r.db.Prepare("UPDATE users SET first_name = $1, last_name = $2, email = $3, password = $4, updated_at = $5 WHERE id = $6")
	if err != nil {
		return err
	}
	defer stmt.Close()

	u.UpdatedAt = time.Now()

	result, err := stmt.Exec(u.FirstName, u.LastName, u.Email, u.Password, timeToNull(u.UpdatedAt), u.ID)
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

func (r User) All() ([]*domain.User, error) {
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

func (r User) Delete(id int64) error {
	stmt, err := r.db.Prepare("DELETE FROM users WHERE id = $1")
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

func (r User) DeleteAll() error {
	stmt, err := r.db.Prepare("TRUNCATE TABLE users RESTART IDENTITY")
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

type scanner interface {
	Scan(dest ...interface{}) error
}

// scanRowUser return nulled fields of the domain object User parsed.
func scanRowUser(s scanner) (*domain.User, error) {
	var updatedAtNull, deletedAtNull sql.NullTime
	m := &domain.User{}

	err := s.Scan(
		&m.ID,
		&m.UUID,
		&m.FirstName,
		&m.LastName,
		&m.Email,
		&m.Password,
		&m.CreatedAt,
		&updatedAtNull,
		&deletedAtNull,
	)
	if err != nil {
		return &domain.User{}, err
	}

	m.UpdatedAt = updatedAtNull.Time
	m.DeletedAt = deletedAtNull.Time

	return m, nil
}

// timeToNull helper to try empty time fields.
func timeToNull(t time.Time) sql.NullTime {
	null := sql.NullTime{Time: t}

	if !null.Time.IsZero() {
		null.Valid = true
	}
	return null
}
