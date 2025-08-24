package pgx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/pgsql"
	"github.com/adrianolmedo/genesis/store"
	"github.com/adrianolmedo/genesis/user"

	"github.com/jackc/pgx/v5"
)

// Product repository.
type Product struct {
	conn *pgx.Conn
}

// Create add one product to the storage.
func (r Product) Create(ctx context.Context, m *store.Product) error {
	m.UUID = genesis.NextUUID()
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
func (r Product) ByID(ctx context.Context, id int64) (*store.Product, error) {
	var updatedAtNull, deletedAtNull sql.NullTime

	m := &store.Product{}

	err := r.conn.QueryRow(ctx, `SELECT * FROM "product" WHERE id = $1 AND deleted_at IS NULL`, id).Scan(
		&m.ID, &m.UUID, &m.Name, &m.Observations, &m.Price, &m.CreatedAt, &updatedAtNull, &deletedAtNull)
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return nil, store.ErrProductNotFound
	}

	if err != nil {
		return nil, err
	}

	m.UpdatedAt = pgsql.NullTimeToPtr(updatedAtNull)
	m.DeletedAt = pgsql.NullTimeToPtr(deletedAtNull)

	return m, nil
}

// Update product.
func (r Product) Update(ctx context.Context, m store.Product) error {
	m.UpdatedAt = pgsql.TimeToPtr(time.Now())

	result, err := r.conn.Exec(ctx,
		`UPDATE "product" SET name = $1, observations = $2, price = $3, updated_at = $4 WHERE id = $5`,
		m.Name, m.Observations, m.Price, m.UpdatedAt, m.ID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return store.ErrProductNotFound
	}

	return nil
}

// All get a collection of all products.
func (r Product) All(ctx context.Context) (store.Products, error) {
	rows, err := r.conn.Query(ctx, `SELECT * FROM "product" WHERE deleted_at IS NULL`)
	if err != nil {
		return store.Products{}, err
	}
	defer rows.Close()

	products := make(store.Products, 0)

	for rows.Next() {
		var updatedAtNull, deletedAtNull sql.NullTime
		m := &store.Product{}

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
			return store.Products{}, err
		}

		m.UpdatedAt = pgsql.NullTimeToPtr(updatedAtNull)
		m.DeletedAt = pgsql.NullTimeToPtr(deletedAtNull)

		products = append(products, m)
	}

	if err := rows.Err(); err != nil {
		return store.Products{}, err
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
		return store.ErrProductNotFound
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
		return user.ErrNotFound
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
