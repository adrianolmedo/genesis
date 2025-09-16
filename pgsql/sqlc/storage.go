package sqlc

import (
	"context"
	"fmt"
	"time"

	"github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/billing"
	"github.com/adrianolmedo/genesis/store"
	"github.com/adrianolmedo/genesis/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Storage represents all repositories.
type Storage struct {
	db       *pgxpool.Pool
	User     *user.Repo
	Product  *store.ProductRepo
	Customer *store.CustomerRepo
	Invoice  *billing.Repo
}

// NewStorage creates a new Storage instance with all repositories.
func NewStorage(ctx context.Context, db *pgxpool.Pool, cfg genesis.Config) (*Storage, error) {
	return &Storage{
		db:       db,
		User:     user.NewRepo(db),
		Product:  store.NewProductRepo(db),
		Customer: store.NewCustomerRepo(db),
		Invoice:  billing.NewRepo(db),
	}, nil
}

// NewPool return a postgres database connection from cfg params.
func NewPool(ctx context.Context, cfg genesis.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't do ping %v", err)
	}
	return pool, nil
}
