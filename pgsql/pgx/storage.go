package pgx

import (
	"fmt"

	"github.com/adrianolmedo/genesis"
)

// Storage represents all repositories.
type Storage struct {
	User User
}

// NewStorage start postgres database connection, build the Storage and return it
// its pointer.
func NewStorage(cfg genesis.Config) (*Storage, error) {
	db, err := newDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres: %v", err)
	}

	return &Storage{
		User: User{conn: db},
	}, nil
}
