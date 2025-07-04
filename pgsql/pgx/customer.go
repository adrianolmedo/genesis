package pgx

import (
	"context"
	"database/sql"
	"time"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/pgsql"

	"github.com/jackc/pgx/v5"
)

// Customer repository.
type Customer struct {
	conn *pgx.Conn
}

// Create return one Customer or SQL error.
func (r Customer) Create(m *domain.Customer) error {
	m.UUID = domain.NextUUID()
	m.CreatedAt = time.Now()

	err := r.conn.QueryRow(context.Background(),
		`INSERT INTO "customer" (uuid, first_name, last_name, email, password, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		m.UUID, m.FirstName, m.LastName, m.Email).Scan(&m.ID)
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

	rows, err := r.conn.Query(context.Background(), query)
	if err != nil {
		return pgsql.PagerResults{}, err
	}
	defer rows.Close()

	customers := make(domain.Customers, 0)

	for rows.Next() {
		var updatedAtNull, deletedAtNull sql.NullTime
		m := &domain.Customer{}

		err := rows.Scan(
			&m.ID,
			&m.UUID,
			&m.FirstName,
			&m.LastName,
			&m.Email,
			&m.Password,
			&m.CreatedAt,
			&updatedAtNull,
			&deletedAtNull,
		)
		if err != nil {
			return pgsql.PagerResults{}, err
		}

		m.UpdatedAt = pgsql.PtrFromNullTime(updatedAtNull)
		m.DeletedAt = pgsql.PtrFromNullTime(deletedAtNull)

		customers = append(customers, m)
	}

	if err := rows.Err(); err != nil {
		return pgsql.PagerResults{}, err
	}

	// Get total rows to calculate total pages.
	var totalRows int64
	err = r.conn.QueryRow(context.Background(), `SELECT COUNT (*) FROM "customer" WHERE deleted_at IS NULL`).Scan(&totalRows)
	if err != nil {
		return pgsql.PagerResults{}, err
	}

	return p.Paginate(customers, totalRows), nil
}

// Delete soft delete Customer from its ID.
func (r Customer) Delete(id int64) error {
	result, err := r.conn.Exec(context.Background(), `UPDATE "customer" SET deleted_at = $1 WHERE id = $2`, time.Now(), id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrCustomerNotFound
	}

	return nil
}
