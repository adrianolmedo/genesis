package postgres

import (
	"database/sql"
	"time"

	domain "github.com/adrianolmedo/genesis"
)

// Customer repository.
type Customer struct {
	db *sql.DB
}

// Create return one Customer or SQL error.
func (r Customer) Create(u *domain.Customer) error {
	stmt, err := r.db.Prepare("INSERT INTO customers (uuid, first_name, last_name, email, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	u.UUID = domain.NextUUID()
	u.CreatedAt = time.Now()

	err = stmt.QueryRow(u.UUID, u.FirstName, u.LastName, u.Email, u.CreatedAt).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

// All return a Customer collection or SQL error.
func (r Customer) All(filter domain.Filter) (domain.Customers, error) {
	query := "SELECT * FROM customers"
	query += " " + orderBy(filter)
	query += " " + limitOffset(filter)

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := make(domain.Customers, 0)

	for rows.Next() {
		c, err := scanRowCustomer(rows)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

// scanRowUser return nulled fields of the domain object User parsed.
func scanRowCustomer(s scanner) (domain.Customer, error) {
	var updatedAtNull, deletedAtNull sql.NullTime
	c := domain.Customer{}

	err := s.Scan(
		&c.ID,
		&c.UUID,
		&c.FirstName,
		&c.LastName,
		&c.Email,
		&c.CreatedAt,
		&updatedAtNull,
		&deletedAtNull,
	)
	if err != nil {
		return domain.Customer{}, err
	}

	c.UpdatedAt = updatedAtNull.Time
	c.DeletedAt = deletedAtNull.Time

	return c, nil
}
