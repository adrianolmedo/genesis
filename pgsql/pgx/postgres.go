package pgx

import (
	"context"
	"fmt"

	config "github.com/adrianolmedo/genesis"

	"github.com/jackc/pgx/v5"
)

// newDB return a postgres database connection from dbcfg params.
func newDB(dbcfg config.DB) (*pgx.Conn, error) {
	// postgres://user:password@host:port/dbname?sslmode=disable
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbcfg.User, dbcfg.Password, dbcfg.Host, dbcfg.Port, dbcfg.Name)

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
