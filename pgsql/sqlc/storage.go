package sqlc

import (
	"context"
	"fmt"

	"github.com/adrianolmedo/genesis"
	"github.com/jackc/pgx/v5"
)

// Storage represents all repositories.
type Storage struct {
	User *User
	//Product  Product
	//Customer Customer
	//Invoice  Invoice
}

// NewStorage start postgres database connection, build the Storage and return it
// its pointer.
func NewStorage(cfg genesis.Config) (*Storage, error) {
	conn, err := newConn(cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres: %v", err)
	}

	return &Storage{
		User: NewUser(conn),
		//Product:  Product{conn: db},
		//Customer: Customer{conn: db},
		/*Invoice: Invoice{
			conn:   db,
			header: InvoiceHeader{conn: db},
			items:  InvoiceItem{conn: db},
		}*,*/
	}, nil
}

// newConn return a postgres database connection from cfg params.
func newConn(cfg genesis.Config) (*pgx.Conn, error) {
	// postgres://user:password@host:port/dbname?sslmode=disable
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can't do ping %v", err)
	}

	//defer conn.Close(context.Background())
	return conn, nil
}
