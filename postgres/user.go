package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	domain "github.com/adrianolmedo/genesis"
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
// TODO: Only select data that have deleted_at empty.
func (r User) ByLogin(email, password string) error {
	stmt, err := r.db.Prepare(`SELECT id FROM "user" WHERE email = $1 AND password = $2`)
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
// TODO: Only select data that have deleted_at empty.
func (r User) ByID(id uint) (*domain.User, error) {
	stmt, err := r.db.Prepare(`SELECT * FROM "user" WHERE id = $1 AND deleted_at IS NULL`)
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

// Update update User.
func (r User) Update(u domain.User) error {
	stmt, err := r.db.Prepare(`UPDATE "user" SET first_name = $1, last_name = $2, email = $3, password = $4, updated_at = $5 WHERE id = $6`)
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
		var updatedNull, deletedNull sql.NullTime
		u := domain.User{}

		err := rows.Scan(
			&u.ID,
			&u.UUID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Password,
			&u.CreatedAt,
			&updatedNull,
			&deletedNull,
		)
		if err != nil {
			return nil, err
		}

		u.UpdatedAt = updatedNull.Time
		u.DeletedAt = deletedNull.Time

		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Delete delete user from its is.
// TODO: Convert to soft delete.
func (r User) Delete(id int64) error {
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
// TODO: Move this query to tests.
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
	Scan(dest ...interface{}) error
}

// scanRowUser return nulled fields of the domain object User parsed.
// TODO: Check how to do this without using scanner interface.
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
