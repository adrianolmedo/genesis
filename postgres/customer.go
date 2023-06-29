package postgres

import (
	"database/sql"
	"fmt"
	"math"
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

func (r Customer) countAll(f domain.Filter) (int, error) {
	stmt, err := r.db.Prepare("SELECT COUNT (*) FROM customers")
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

// All build with limit, offset and order the filter results and then return it
// for the pagination or a SQL error.
func (r Customer) All(f domain.Filter) (domain.FilterResults, error) {
	// Get data with limit, offset & order.
	query := "SELECT * FROM customers"
	query += " " + orderBy(f)
	//query += " " + limitOffset(f)
	offset := f.Page * f.Limit
	query += " " + fmt.Sprintf("LIMIT %d OFFSET %d", f.Limit, offset)

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return domain.FilterResults{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return domain.FilterResults{}, err
	}
	defer rows.Close()

	customers := make(domain.Customers, 0)
	for rows.Next() {
		c, err := scanRowCustomer(rows)
		if err != nil {
			return domain.FilterResults{}, err
		}
		customers = append(customers, c)
	}
	if err := rows.Err(); err != nil {
		return domain.FilterResults{}, err
	}

	// Get total rows to calculate total pages.
	totalRows, err := r.countAll(f)
	if err != nil {
		return domain.FilterResults{}, err
	}

	totalPages := int(math.Ceil(float64(totalRows)/float64(f.Limit))) - 1

	var fromRow, toRow int

	// Set fromRow and toRow on first page.
	if f.Page == 0 {
		fromRow = 1
		toRow = f.Limit
	} else {
		if f.Page <= totalPages {
			// Calculate fromRow and toRow.
			fromRow = f.Page*f.Limit + 1
			toRow = (f.Page + 1) * f.Limit
		}
	}

	// Or set toRow with totalRows.
	if toRow > totalRows {
		toRow = totalRows
	}

	return domain.FilterResults{
		TotalRows:  totalRows,
		TotalPages: totalPages,
		FromRow:    fromRow,
		ToRow:      toRow,
		Rows:       customers,
	}, nil
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
