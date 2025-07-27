package sqlc

import (
	"context"
	"fmt"

	"github.com/adrianolmedo/genesis"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Storage represents all repositories.
type Storage struct {
	db       *pgxpool.Pool
	User     *User
	Product  *Product
	Customer *Customer
	Invoice  *Invoice
}

// NewStorage start postgres database connection, build the Storage and return it
// its pointer.
func NewStorage(ctx context.Context, cfg genesis.Config) (*Storage, error) {
	db, err := newPool(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres: %v", err)
	}

	return &Storage{
		User:     NewUser(db),
		Product:  NewProduct(db),
		Customer: NewCustomer(db),
		Invoice:  NewInvoice(db),
	}, nil
}

// Close releases the storage db connection resources.
// main function should call this method before exiting.
func (s *Storage) Close() {
	s.db.Close()
}

// newPool return a postgres database connection from cfg params.
func newPool(ctx context.Context, cfg genesis.Config) (*pgxpool.Pool, error) {
	// postgres://user:password@host:port/dbname?sslmode=disable
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	// Optional: test connection.
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't do ping %v", err)
	}

	return pool, nil
}
