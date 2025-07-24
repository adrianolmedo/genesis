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

// Product repository.
type Product struct {
	conn *pgx.Conn
}

// Create add one product to the storage.
func (r Product) Create(ctx context.Context, m *domain.Product) error {
	m.UUID = domain.NextUUID()
	m.CreatedAt = time.Now()

	err := r.conn.QueryRow(ctx,
		`INSERT INTO "product" (uuid, name, observations, price, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		m.UUID, m.Name, m.Observations, m.Price, m.CreatedAt).Scan(&m.ID)
	if err != nil {
		return err
	}

	return nil
}

// ByID get one product by its id.
func (r Product) ByID(ctx context.Context, id int64) (*domain.Product, error) {
	var updatedAtNull, deletedAtNull sql.NullTime

	m := &domain.Product{}

	err := r.conn.QueryRow(ctx, `SELECT * FROM "product" WHERE id = $1 AND deleted_at IS NULL`, id).Scan(
		&m.ID, &m.UUID, &m.Name, &m.Observations, &m.Price, &m.CreatedAt, &updatedAtNull, &deletedAtNull)
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrProductNotFound
	}

	if err != nil {
		return nil, err
	}

	m.UpdatedAt = pgsql.PtrFromNullTime(updatedAtNull)
	m.DeletedAt = pgsql.PtrFromNullTime(deletedAtNull)

	return m, nil
}

// Update product.
func (r Product) Update(ctx context.Context, m domain.Product) error {
	m.UpdatedAt = pgsql.TimePtr(time.Now())

	result, err := r.conn.Exec(ctx,
		`UPDATE "product" SET name = $1, observations = $2, price = $3, updated_at = $4 WHERE id = $5`,
		m.Name, m.Observations, m.Price, m.UpdatedAt, m.ID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

// All get a collection of all products.
func (r Product) All(ctx context.Context) (domain.Products, error) {
	rows, err := r.conn.Query(ctx, `SELECT * FROM "product" WHERE deleted_at IS NULL`)
	if err != nil {
		return domain.Products{}, err
	}
	defer rows.Close()

	products := make(domain.Products, 0)

	for rows.Next() {
		var updatedAtNull, deletedAtNull sql.NullTime
		m := &domain.Product{}

		err := rows.Scan(
			&m.ID,
			&m.UUID,
			&m.Name,
			&m.Observations,
			&m.CreatedAt,
			&updatedAtNull,
			&deletedAtNull,
		)
		if err != nil {
			return domain.Products{}, err
		}

		m.UpdatedAt = pgsql.PtrFromNullTime(updatedAtNull)
		m.DeletedAt = pgsql.PtrFromNullTime(deletedAtNull)

		products = append(products, m)
	}

	if err := rows.Err(); err != nil {
		return domain.Products{}, err
	}

	return products, nil
}

// Delete product by its id.
func (r Product) Delete(ctx context.Context, id int64) error {
	result, err := r.conn.Exec(ctx, `UPDATE "product" SET deleted_at = $1 WHERE id = $2`, time.Now(), id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

// HardDelete delete product from the storage.
func (r Product) HardDelete(ctx context.Context, id int64) error {
	result, err := r.conn.Exec(ctx, `DELETE FROM "product" WHERE id = $1`)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// DeleteAll delete all products.
func (r Product) DeleteAll(ctx context.Context) error {
	_, err := r.conn.Exec(ctx, `TRUNCATE TABLE "product" RESTART IDENTITY CASCADE`)
	if err != nil {
		return fmt.Errorf("can't truncate table: %v", err)
	}

	return nil
}
