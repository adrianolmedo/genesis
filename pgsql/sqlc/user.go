package sqlc

import (
	"context"
	"time"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/pgsql/sqlc/dbgen"
)

// User repository.
type User struct {
	dbg *dbgen.Queries
}

func (r User) Create(m *domain.User) error {
	m.UUID = domain.NextUUID()
	m.CreatedAt = time.Now()

	id, err := r.dbg.UserCreate(context.Background(), dbgen.UserCreateParams{
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
	})

	if err != nil {
		return err
	}

	m.ID = uint(id)

	return nil
}
