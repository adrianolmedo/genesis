package pgx

import (
	"context"
	"fmt"

	"github.com/adrianolmedo/genesis"

	"github.com/jackc/pgx/v5"
)

// newDB return a postgres database connection from dbcfg params.
func newDB(ctx context.Context, cfg genesis.Config) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, cfg.DBURL)
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
