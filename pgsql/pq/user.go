package pq

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/pgsql"
)

// User repository.
type User struct {
	db *sql.DB
}

// Create a User to the storage.
func (r User) Create(u *domain.User) error {
	stmt, err := r.db.Prepare(`INSERT INTO "user" 
	(uuid, first_name, last_name, email, password, created_at)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	u.UUID = domain.NextUUID()
	u.CreatedAt = time.Now()

	err = stmt.QueryRow(u.UUID, u.FirstName, u.LastName, u.Email, u.Password, u.CreatedAt).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

// ByLogin get a User from its login data.
func (r User) ByLogin(email, password string) error {
	stmt, err := r.db.Prepare(`SELECT id FROM "user" WHERE email = $1 AND password = $2 AND deleted_at IS NULL`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(email, password)
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

// ByID get a User from its id.
func (r User) ByID(id uint) (*domain.User, error) {
	stmt, err := r.db.Prepare(`SELECT * FROM "user" WHERE id = $1 AND deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	u, err := scanRowUser(stmt.QueryRow(id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return u, nil
}

// Update update User.
func (r User) Update(u domain.User) error {
	stmt, err := r.db.Prepare(`UPDATE "user" SET first_name = $1, last_name = $2, email = $3, password = $4, updated_at = $5 WHERE id = $6`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	u.UpdatedAt = pgsql.PtrTime(time.Now())

	result, err := stmt.Exec(u.FirstName, u.LastName, u.Email, u.Password, u.UpdatedAt, u.ID)
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

// All get User collection.
func (r User) All() (domain.Users, error) {
	stmt, err := r.db.Prepare(`SELECT * FROM "user" WHERE deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make(domain.Users, 0)

	for rows.Next() {
		var updatedAtNull, deletedAtNull sql.NullTime
		u := domain.User{}

		err := rows.Scan(
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
			return nil, err
		}

		u.UpdatedAt = pgsql.ToTimePtr(updatedAtNull)
		u.DeletedAt = pgsql.ToTimePtr(deletedAtNull)

		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Delete mark user as deleted.
func (r User) Delete(id uint) error {
	stmt, err := r.db.Prepare(`UPDATE "user" SET deleted_at = $1 WHERE id = $2`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(time.Now(), id)
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

// HardDelete delete user from the storage.
func (r User) HardDelete(id uint) error {
	stmt, err := r.db.Prepare(`DELETE FROM "user" WHERE id = $1`)
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

// DeleteAll delete all users.
func (r User) DeleteAll() error {
	stmt, err := r.db.Prepare(`TRUNCATE TABLE "user" RESTART IDENTITY`)
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
	Scan(dest ...any) error
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

	m.UpdatedAt = pgsql.ToTimePtr(updatedAtNull)
	m.DeletedAt = pgsql.ToTimePtr(deletedAtNull)

	return m, nil
}
