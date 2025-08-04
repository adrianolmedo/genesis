package pq

import (
	"database/sql"
	"fmt"

	"github.com/adrianolmedo/genesis"

	// blank import to init postgres library.
	_ "github.com/lib/pq"
)

// newDB return a postgres database connection from dbcfg params.
func newDB(cfg genesis.Config) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", cfg.DBURL)
	if err != nil {
		return nil, fmt.Errorf("can't open the data base %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("can't do ping %v", err)
	}

	//log.Println("Connected to postgres!")
	return db, nil
}
