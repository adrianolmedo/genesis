package pq

import (
	"database/sql"
	"time"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/pgsql"
)

// Customer repository.
type Customer struct {
	db *sql.DB
}

// Create return one Customer or SQL error.
func (r Customer) Create(u *domain.Customer) error {
	stmt, err := r.db.Prepare(`INSERT INTO "customer" (uuid, first_name, last_name, email, password, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	u.UUID = domain.NextUUID()
	u.CreatedAt = time.Now()

	err = stmt.QueryRow(u.UUID, u.FirstName, u.LastName, u.Email, u.Password, u.CreatedAt).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

// All return filtered results by limit, offset and order for the pagination
// or return a SQL error.
func (r Customer) All(p *pgsql.Pager) (pgsql.PagerResults, error) {
	query := `SELECT * FROM "customer"`
	query += " " + p.OrderBy()
	query += " " + p.LimitOffset()

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return pgsql.PagerResults{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return pgsql.PagerResults{}, err
	}
	defer rows.Close()

	customers := make(domain.Customers, 0)
	for rows.Next() {
		c, err := scanRowCustomer(rows)
		if err != nil {
			return pgsql.PagerResults{}, err
		}
		customers = append(customers, c)
	}
	if err := rows.Err(); err != nil {
		return pgsql.PagerResults{}, err
	}

	// Get total rows to calculate total pages.
	totalRows, err := r.countAll()
	if err != nil {
		return pgsql.PagerResults{}, err
	}

	return p.Paginate(customers, totalRows), nil
}

// countAll return total of Customers in storage.
func (r Customer) countAll() (int, error) {
	stmt, err := r.db.Prepare(`SELECT COUNT (*) FROM "customer"`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var n int
	err = stmt.QueryRow().Scan(&n)
	if err != nil {
		return 0, err
	}

	return n, nil
}

// Delete delete Customer from its ID.
func (r Customer) Delete(id int) error {
	stmt, err := r.db.Prepare(`DELETE FROM "customer" WHERE id = $1`)
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
		return domain.ErrCustomerNotFound
	}
	return nil
}

// scanRowUser return nulled fields of the domain object User parsed.
func scanRowCustomer(s scanner) (*domain.Customer, error) {
	var updatedAtNull, deletedAtNull sql.NullTime
	cx := &domain.Customer{}

	err := s.Scan(
		&cx.ID,
		&cx.UUID,
		&cx.FirstName,
		&cx.LastName,
		&cx.Email,
		&cx.Password,
		&cx.CreatedAt,
		&updatedAtNull,
		&deletedAtNull,
	)
	if err != nil {
		return nil, err
	}

	cx.UpdatedAt = updatedAtNull.Time
	cx.DeletedAt = deletedAtNull.Time

	return cx, nil
}
