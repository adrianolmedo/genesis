package storage

import (
	"database/sql"
	"fmt"

	"github.com/adrianolmedo/go-restapi/config"
)

func New(dbcfg config.DB) (*sql.DB, error) {
	if dbcfg.Engine == "" {
		return nil, fmt.Errorf("database engine not selected")
	}

	if dbcfg.Engine != "postgres" {
		return nil, fmt.Errorf("database engine '%s' not implemented", dbcfg.Engine)
	}

	db, err := newPSQL(dbcfg)
	if err != nil {
		return nil, err
	}

	return db, nil
}
