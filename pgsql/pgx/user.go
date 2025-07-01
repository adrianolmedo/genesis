package pgx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/pgsql"

	"github.com/jackc/pgx/v5"
)

// User repository.
type User struct {
	conn *pgx.Conn
}

// Create a User to the storage.
func (r User) Create(u *domain.User) error {
	u.UUID = domain.NextUUID()
	u.CreatedAt = time.Now()

	err := r.conn.QueryRow(context.Background(),
		`INSERT INTO "user" (uuid, first_name, last_name, email, password, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		u.UUID, u.FirstName, u.LastName, u.Email, u.Password, u.CreatedAt).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

// ByLogin get a User from its login data.
func (r User) ByLogin(email, password string) error {
	result, err := r.conn.Exec(context.Background(), `SELECT id FROM "user" WHERE email = $1 AND password = $2 AND deleted_at IS NULL`, email, password)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// ByID get a User from its id.
func (r User) ByID(id int64) (*domain.User, error) {
	var updatedAtNull, deletedAtNull sql.NullTime

	m := &domain.User{}

	err := r.conn.QueryRow(context.Background(), `SELECT * FROM "user" WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&m.ID, &m.UUID, &m.FirstName, &m.LastName, &m.Email, &m.Password, &m.CreatedAt, &updatedAtNull, &deletedAtNull)
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	m.UpdatedAt = pgsql.ToTimePtr(updatedAtNull)
	m.DeletedAt = pgsql.ToTimePtr(deletedAtNull)

	return m, nil
}

// Update user.
func (r User) Update(m domain.User) error {
	m.UpdatedAt = pgsql.PtrTime(time.Now())

	result, err := r.conn.Exec(context.Background(),
		`UPDATE "user" SET first_name = $1, last_name = $2, email = $3, password = $4, updated_at = $5 WHERE id = $6`,
		m.FirstName, m.LastName, m.Email, m.Password, m.UpdatedAt, m.ID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r User) All(p *pgsql.Pager) (pgsql.PagerResults, error) {
	query := `SELECT id, uuid, first_name, last_name, email, password, created_at, updated_at, deleted_at FROM "user" WHERE deleted_at IS NULL`
	query += " " + p.OrderBy()
	query += " " + p.LimitOffset()

	rows, err := r.conn.Query(context.Background(), query)
	if err != nil {
		return pgsql.PagerResults{}, err
	}
	defer rows.Close()

	users := make(domain.Users, 0)

	for rows.Next() {
		var updatedAtNull, deletedAtNull sql.NullTime
		m := &domain.User{}

		err := rows.Scan(
			&m.ID,
			&m.UUID,
			&m.FirstName,
			&m.LastName,
			&m.CreatedAt,
			&updatedAtNull,
			&deletedAtNull,
		)
		if err != nil {
			return pgsql.PagerResults{}, err
		}

		m.UpdatedAt = pgsql.ToTimePtr(updatedAtNull)
		m.DeletedAt = pgsql.ToTimePtr(deletedAtNull)

		users = append(users, m)
	}

	if err := rows.Err(); err != nil {
		return pgsql.PagerResults{}, err
	}

	// Get total rows to calculate total pages.
	var totalRows int64
	err = r.conn.QueryRow(context.Background(), `SELECT COUNT (*) FROM "user" WHERE deleted_at IS NULL`).Scan(&totalRows)
	if err != nil {
		return pgsql.PagerResults{}, err
	}

	return p.Paginate(users, totalRows), nil
}

func (r User) Delete(id int64) error {
	result, err := r.conn.Exec(context.Background(), `UPDATE "user" SET deleted_at = $1 WHERE id = $2`, time.Now(), id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r User) HardDelete(id int64) error {
	result, err := r.conn.Exec(context.Background(), `DELETE FROM "user" WHERE id = $1`)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r User) DeleteAll() error {
	_, err := r.conn.Exec(context.Background(), `TRUNCATE TABLE "user" RESTART IDENTITY`)
	if err != nil {
		return fmt.Errorf("can't truncate table: %v", err)
	}

	return nil
}
