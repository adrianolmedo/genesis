package pgx

import (
	"fmt"

	config "github.com/adrianolmedo/genesis"
)

// Storage represents all repositories.
type Storage struct {
	User User
}

// NewStorage start postgres database connection, build the Storage and return it
// its pointer.
func NewStorage(dbcfg config.DB) (*Storage, error) {
	db, err := newDB(dbcfg)
	if err != nil {
		return nil, fmt.Errorf("postgres: %v", err)
	}

	return &Storage{
		User: User{conn: db},
	}, nil
}
