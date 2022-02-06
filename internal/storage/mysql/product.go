package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/adrianolmedo/go-restapi/internal/domain"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r ProductRepository) Create(product *domain.Product) error {
	query := "INSERT INTO products (name, observations, price, created_at) VALUES(?, ?, ?, ?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	product.CreatedAt = time.Now()

	result, err := stmt.Exec(product.Name, product.Observations, product.Price, product.CreatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	product.ID = id // Important!
	return nil
}

func (r ProductRepository) ByID(id int64) (*domain.Product, error) {
	stmt, err := r.db.Prepare("SELECT * FROM products WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// As QueryRow returns a rows we can pass it directly to the mapping
	product, err := scanRowProduct(stmt.QueryRow(id))
	if errors.Is(err, sql.ErrNoRows) {
		return &domain.Product{}, domain.ErrProductNotFound
	}

	if err != nil {
		return &domain.Product{}, err
	}

	return product, nil
}

func (r ProductRepository) Update(product domain.Product) error {
	stmt, err := r.db.Prepare("UPDATE products SET name = ?, observations = ?, price = ?, updated_at = ? WHERE id = ?")
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

func (r ProductRepository) All() ([]*domain.Product, error) {
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

	products := make([]*domain.Product, 0)
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

func (r ProductRepository) Delete(id int64) error {
	stmt, err := r.db.Prepare("DELETE FROM products WHERE id = ?")
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

func (r ProductRepository) Reset() error {
	stmt, err := r.db.Prepare("ALTER TABLE products AUTO_INCREMENT = 1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("can't alter table: %v", err)
	}

	return nil
}

// scanRowProduct return nulled fields of Product parsed.
func scanRowProduct(s scanner) (*domain.Product, error) {
	var updatedAtNull sql.NullTime
	p := &domain.Product{}

	err := s.Scan(
		&p.ID,
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
