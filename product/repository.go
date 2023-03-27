package product

import (
	"database/sql"
	"time"
)

type repository struct {
	db *sql.DB
}

func newRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r repository) Create(product *Product) error {
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
