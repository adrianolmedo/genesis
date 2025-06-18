package pgx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	domain "github.com/adrianolmedo/genesis"

	"github.com/jackc/pgx/v5"
)

// Product repository.
type Product struct {
	conn *pgx.Conn
}

// Create create one product.
func (r Product) Create(p *domain.Product) error {
	p.UUID = domain.NextUUID()
	p.CreatedAt = time.Now()

	err := r.conn.QueryRow(context.Background(),
		`INSERT INTO "product"
		(uuid, name, observations, price, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		p.UUID, p.Name, p.Observations, p.Price, p.CreatedAt).Scan(&p.ID)
	if err != nil {
		return err
	}

	return nil
}

// ByID get one product by its id.
func (r Product) ByID(id int) (*domain.Product, error) {
	m := &domain.Product{}

	err := r.conn.QueryRow(context.Background(), `SELECT * FROM "product" WHERE id = $1 AND deleted_at IS NULL`, id).Scan(
		m.ID, m.UUID, m.Name, m.Observations, m.Price, m.CreatedAt, m.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrProductNotFound
	}

	return m, nil
}

// Update product.
func (r Product) Update(m domain.Product) error {
	m.UpdatedAt = time.Now()

	result, err := r.conn.Exec(context.Background(),
		`UPDATE "product" SET name = $1, observations = $2, price = $3, updated_at = $4 WHERE id = $5`,
		m.Name, m.Observations, m.Price, m.UpdatedAt, m.ID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// All get a collection of all prodycts.
func (r Product) All() (domain.Products, error) {
	rows, err := r.conn.Query(context.Background(), `SELECT * FROM "product" WHERE deleted_at IS NULL`)
	if err != nil {
		return domain.Products{}, err
	}

	products := make(domain.Products, 0)

	for rows.Next() {
		var updatedAtNull, createdAtNull sql.NullTime
		m := &domain.Product{}

		err := rows.Scan(
			&m.ID,
			&m.UUID,
			&m.Name,
			&m.Observations,
			&m.CreatedAt,
			&updatedAtNull,
			&createdAtNull,
		)
		if err != nil {
			return domain.Products{}, err
		}

		m.UpdatedAt = updatedAtNull.Time
		m.CreatedAt = createdAtNull.Time

		products = append(products, m)
	}

	if err := rows.Err(); err != nil {
		return domain.Products{}, err
	}

	return products, nil
}

// Delete product by its id.
func (r Product) Delete(id int) error {
	result, err := r.conn.Exec(context.Background(), `UPDATE "product" SET deleted_at = $1 WHERE id = $2`, time.Now(), id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

// HardDelete delete product from the storage.
func (r Product) HardDelete(id uint) error {
	result, err := r.conn.Exec(context.Background(), `DELETE FROM "product" WHERE id = $1`)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// DeleteAll delete all products.
func (r Product) DeleteAll() error {
	_, err := r.conn.Exec(context.Background(), `TRUNCATE TABLE "product" RESTART IDENTITY`)
	if err != nil {
		return fmt.Errorf("can't truncate table: %v", err)
	}

	return nil
}
