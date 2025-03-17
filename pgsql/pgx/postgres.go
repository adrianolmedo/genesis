package pgx

import (
	"context"
	"fmt"

	"github.com/adrianolmedo/genesis"

	"github.com/jackc/pgx/v5"
)

// newDB return a postgres database connection from dbcfg params.
func newDB(cfg genesis.Config) (*pgx.Conn, error) {
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
