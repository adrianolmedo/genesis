package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) Create(product *Product) error {
	stmt, err := r.db.Prepare("INSERT INTO products (name, observations, price, created_at) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	product.CreatedAt = time.Now()

	err = stmt.QueryRow(product.Name, product.Observations, product.Price, product.CreatedAt).Scan(&product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) ByID(id int64) (*Product, error) {
	stmt, err := r.db.Prepare("SELECT * FROM products WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	product, err := scanRowProduct(stmt.QueryRow(id))
	if errors.Is(err, sql.ErrNoRows) {
		return &Product{}, ErrProductNotFound
	}

	if err != nil {
		return &Product{}, err
	}

	return product, nil
}

func (r Repository) Update(product Product) error {
	stmt, err := r.db.Prepare("UPDATE products SET name = $1, observations = $2, price = $3, updated_at = $4 WHERE id = $5")
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
		return ErrProductNotFound
	}

	return nil
}

func (r Repository) All() ([]*Product, error) {
	stmt, err := r.db.Prepare("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*Product, 0)

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

func (r Repository) Delete(id int64) error {
	stmt, err := r.db.Prepare("DELETE FROM products WHERE id = $1")
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
		return ErrProductNotFound
	}
	return nil
}

func (r Repository) DeleteAll() error {
	stmt, err := r.db.Prepare("TRUNCATE TABLE products RESTART IDENTITY CASCADE")
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

// helpers...

type scanner interface {
	Scan(dest ...interface{}) error
}

// scanRowUser return nulled fields of User parsed.
func scanRowProduct(s scanner) (*Product, error) {
	var updatedAtNull sql.NullTime
	p := &Product{}

	err := s.Scan(
		&p.ID,
		&p.Name,
		&p.Observations,
		&p.Price,
		&p.CreatedAt,
		&updatedAtNull,
	)
	if err != nil {
		return &Product{}, err
	}

	p.UpdatedAt = updatedAtNull.Time

	return p, nil
}

// timeToNull helper to try empty time fields.
func timeToNull(t time.Time) sql.NullTime {
	null := sql.NullTime{Time: t}

	if !null.Time.IsZero() {
		null.Valid = true
	}
	return null
}
