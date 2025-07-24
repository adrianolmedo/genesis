package sqlc

import (
	"context"
	"fmt"

	"github.com/adrianolmedo/genesis"

	"github.com/jackc/pgx/v5"
)

// Storage represents all repositories.
type Storage struct {
	User     *User
	Product  *Product
	Customer *Customer
	Invoice  *Invoice
}

// NewStorage start postgres database connection, build the Storage and return it
// its pointer.
func NewStorage(ctx context.Context, cfg genesis.Config) (*Storage, error) {
	conn, err := newConn(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres: %v", err)
	}

	return &Storage{
		User:     NewUser(conn),
		Product:  NewProduct(conn),
		Customer: NewCustomer(conn),
		Invoice:  NewInvoice(conn),
	}, nil
}

// newConn return a postgres database connection from cfg params.
func newConn(ctx context.Context, cfg genesis.Config) (*pgx.Conn, error) {
	// postgres://user:password@host:port/dbname?sslmode=disable
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't do ping %v", err)
	}

	//defer conn.Close(ctx)
	return conn, nil
}
