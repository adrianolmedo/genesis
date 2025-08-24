package pgx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/pgsql"
	"github.com/adrianolmedo/genesis/user"

	"github.com/jackc/pgx/v5"
)

// UserRepo repository.
type UserRepo struct {
	conn *pgx.Conn
}

// Create a User to the storage.
func (r UserRepo) Create(ctx context.Context, u *user.User) error {
	u.UUID = genesis.NextUUID()
	u.CreatedAt = time.Now()

	err := r.conn.QueryRow(ctx,
		`INSERT INTO "user" (uuid, first_name, last_name, email, password, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		u.UUID, u.FirstName, u.LastName, u.Email, u.Password, u.CreatedAt).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

// ByLogin get a User from its login data.
func (r UserRepo) ByLogin(ctx context.Context, email, password string) error {
	result, err := r.conn.Exec(ctx, `SELECT id FROM "user" WHERE email = $1 AND password = $2 AND deleted_at IS NULL`, email, password)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return user.ErrNotFound
	}

	return nil
}

// ByID get a User from its id.
func (r UserRepo) ByID(ctx context.Context, id int64) (*user.User, error) {
	var updatedAtNull, deletedAtNull sql.NullTime

	m := &user.User{}

	err := r.conn.QueryRow(ctx, `SELECT * FROM "user" WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&m.ID, &m.UUID, &m.FirstName, &m.LastName, &m.Email, &m.Password, &m.CreatedAt, &updatedAtNull, &deletedAtNull)
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return nil, user.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	m.UpdatedAt = pgsql.NullTimeToPtr(updatedAtNull)
	m.DeletedAt = pgsql.NullTimeToPtr(deletedAtNull)

	return m, nil
}

// Update user.
func (r UserRepo) Update(ctx context.Context, m user.User) error {
	m.UpdatedAt = pgsql.TimeToPtr(time.Now())

	result, err := r.conn.Exec(ctx,
		`UPDATE "user" SET first_name = $1, last_name = $2, email = $3, password = $4, updated_at = $5 WHERE id = $6`,
		m.FirstName, m.LastName, m.Email, m.Password, m.UpdatedAt, m.ID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return user.ErrNotFound
	}

	return nil
}

func (r UserRepo) All(ctx context.Context, p *pgsql.Pager) (pgsql.PagerResults, error) {
	query := `SELECT id, uuid, first_name, last_name, email, password, created_at, updated_at, deleted_at FROM "user" WHERE deleted_at IS NULL`
	query += " " + p.OrderBy()
	query += " " + p.LimitOffset()

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return pgsql.PagerResults{}, err
	}
	defer rows.Close()

	users := make(user.Users, 0)

	for rows.Next() {
		var updatedAtNull, deletedAtNull sql.NullTime
		m := &user.User{}

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

		m.UpdatedAt = pgsql.NullTimeToPtr(updatedAtNull)
		m.DeletedAt = pgsql.NullTimeToPtr(deletedAtNull)

		users = append(users, m)
	}

	if err := rows.Err(); err != nil {
		return pgsql.PagerResults{}, err
	}

	// Get total rows to calculate total pages.
	var totalRows int64
	err = r.conn.QueryRow(ctx, `SELECT COUNT (*) FROM "user" WHERE deleted_at IS NULL`).Scan(&totalRows)
	if err != nil {
		return pgsql.PagerResults{}, err
	}

	return p.Paginate(users, totalRows), nil
}

func (r UserRepo) Delete(ctx context.Context, id int64) error {
	result, err := r.conn.Exec(ctx, `UPDATE "user" SET deleted_at = $1 WHERE id = $2`, time.Now(), id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return user.ErrNotFound
	}

	return nil
}

func (r UserRepo) HardDelete(ctx context.Context, id int64) error {
	result, err := r.conn.Exec(ctx, `DELETE FROM "user" WHERE id = $1`)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return user.ErrNotFound
	}

	return nil
}

func (r UserRepo) DeleteAll(ctx context.Context) error {
	_, err := r.conn.Exec(ctx, `TRUNCATE TABLE "user" RESTART IDENTITY`)
	if err != nil {
		return fmt.Errorf("can't truncate table: %v", err)
	}

	return nil
}
