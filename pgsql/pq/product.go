package pq

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/pgsql"
)

// Product repository.
type Product struct {
	db *sql.DB
}

// Create create one product.
func (r Product) Create(m *domain.Product) error {
	stmt, err := r.db.Prepare(`INSERT INTO "product" (uuid, name, observations, price, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	m.UUID = domain.NextUUID()
	m.CreatedAt = time.Now()

	err = stmt.QueryRow(m.UUID, m.Name, m.Observations, m.Price, m.CreatedAt).Scan(&m.ID)
	if err != nil {
		return err
	}

	return nil
}

// ByID get one product by its id.
func (r Product) ByID(id int64) (*domain.Product, error) {
	stmt, err := r.db.Prepare(`SELECT * FROM "product" WHERE id = $1 AND deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	product, err := scanRowProduct(stmt.QueryRow(id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrProductNotFound
	}

	if err != nil {
		return nil, err
	}

	return product, nil
}

// Update product.
func (r Product) Update(p domain.Product) error {
	stmt, err := r.db.Prepare(`UPDATE "product" SET name = $1, observations = $2, price = $3, updated_at = $4 WHERE id = $5`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	p.UpdatedAt = pgsql.TimePtr(time.Now())

	result, err := stmt.Exec(p.Name, p.Observations, p.Price, p.UpdatedAt, p.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

// All get a collection of all prodycts.
func (r Product) All() (domain.Products, error) {
	stmt, err := r.db.Prepare(`SELECT * FROM "product" WHERE deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make(domain.Products, 0)

	for rows.Next() {
		p, err := scanRowProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// Delete product by its id.
func (r Product) Delete(id int) error {
	stmt, err := r.db.Prepare(`UPDATE "product" SET deleted_at = $1 WHERE id = $2`)
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
		return domain.ErrProductNotFound
	}
	return nil
}

// HardDelete delete product from the storage.
func (r Product) HardDelete(id uint) error {
	stmt, err := r.db.Prepare(`DELETE FROM "product" WHERE id = $1`)
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
		return domain.ErrProductNotFound
	}
	return nil
}

// DeleteAll delete all products.
func (r Product) DeleteAll() error {
	stmt, err := r.db.Prepare(`TRUNCATE TABLE "product" RESTART IDENTITY CASCADE`)
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

// scanRowProduct return null fields of the domain object Product parsed.
func scanRowProduct(s scanner) (*domain.Product, error) {
	var updatedAtNull, deletedAtNull sql.NullTime
	p := &domain.Product{}

	err := s.Scan(
		&p.ID,
		&p.UUID,
		&p.Name,
		&p.Observations,
		&p.Price,
		&p.CreatedAt,
		&updatedAtNull,
		&deletedAtNull,
	)
	if err != nil {
		return &domain.Product{}, err
	}

	p.UpdatedAt = pgsql.PtrFromNullTime(updatedAtNull)
	p.DeletedAt = pgsql.PtrFromNullTime(deletedAtNull)

	return p, nil
}
