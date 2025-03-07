package pq

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	domain "github.com/adrianolmedo/genesis"
)

// Product repository.
type Product struct {
	db *sql.DB
}

// Create create one product.
func (p Product) Create(product *domain.Product) error {
	stmt, err := p.db.Prepare(`INSERT INTO "product" (uuid, name, observations, price, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	product.UUID = domain.NextUUID()
	product.CreatedAt = time.Now()

	err = stmt.QueryRow(product.UUID, product.Name, product.Observations, product.Price, product.CreatedAt).Scan(&product.ID)
	if err != nil {
		return err
	}

	return nil
}

// ByID get one product by its id.
// TODO: Only select data that have deleted_at empty.
func (p Product) ByID(id int) (*domain.Product, error) {
	stmt, err := p.db.Prepare(`SELECT * FROM "product" WHERE id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	product, err := scanRowProduct(stmt.QueryRow(id))
	if errors.Is(err, sql.ErrNoRows) {
		return &domain.Product{}, domain.ErrProductNotFound
	}

	if err != nil {
		return &domain.Product{}, err
	}

	return product, nil
}

// Update product.
func (p Product) Update(product domain.Product) error {
	stmt, err := p.db.Prepare(`UPDATE "product" SET name = $1, observations = $2, price = $3, updated_at = $4 WHERE id = $5`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	product.UpdatedAt = time.Now()

	result, err := stmt.Exec(product.Name, product.Observations, product.Price, timeToNull(product.UpdatedAt), product.ID)
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
// TODO: Only select data that have deleted_at empty.
func (p Product) All() (domain.Products, error) {
	stmt, err := p.db.Prepare(`SELECT * FROM "product"`)
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
// TODO: Convert to soft delete.
func (p Product) Delete(id int) error {
	stmt, err := p.db.Prepare(`DELETE FROM "product" WHERE id = $1`)
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
// TODO: Move to tests.
func (p Product) DeleteAll() error {
	stmt, err := p.db.Prepare(`TRUNCATE TABLE "product" RESTART IDENTITY CASCADE`)
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
// TODO: Check how to do this without using scanner interface.
func scanRowProduct(s scanner) (*domain.Product, error) {
	var updatedAtNull sql.NullTime
	p := &domain.Product{}

	err := s.Scan(
		&p.ID,
		&p.UUID,
		&p.Name,
		&p.Observations,
		&p.Price,
		&p.CreatedAt,
		&updatedAtNull,
	)
	if err != nil {
		return &domain.Product{}, err
	}

	p.UpdatedAt = updatedAtNull.Time

	return p, nil
}
