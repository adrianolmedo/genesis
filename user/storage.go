package user

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/adrianolmedo/go-restapi/config"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(dbcfg config.Database) (*Repository, error) {
	if dbcfg.Engine != "postgres" {
		return nil, fmt.Errorf("database engine '%s' not implemented", dbcfg.Engine)
	}

	db, err := NewDB(dbcfg)
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r Repository) ByID(id int64) (*User, error) {
	stmt, err := r.db.Prepare("SELECT * FROM users WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	user, err := scanRowUser(stmt.QueryRow(id))
	if errors.Is(err, sql.ErrNoRows) {
		return &User{}, ErrUserNotFound
	}

	if err != nil {
		return &User{}, err
	}

	return user, nil
}

type scanner interface {
	Scan(dest ...interface{}) error
}

// scanRowUser return nulled fields of User parsed.
func scanRowUser(s scanner) (*User, error) {
	var updatedAtNull, deletedAtNull sql.NullTime
	user := &User{}

	err := s.Scan(
		&user.ID,
		//&user.UUID,
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
	user.DeletedAt = deletedAtNull.Time

	return user, nil
}
