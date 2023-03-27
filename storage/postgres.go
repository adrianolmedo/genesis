package storage

import (
	"database/sql"
	"fmt"

	"github.com/adrianolmedo/go-restapi/config"

	_ "github.com/lib/pq"
)

// newPSQL return a postgres database connection.
func newPSQL(dbcfg config.DB) (db *sql.DB, err error) {
	// postgres://user:password@host:port/dbname?sslmode=disable
	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbcfg.User, dbcfg.Password, dbcfg.Host, dbcfg.Port, dbcfg.Name)

	db, err = sql.Open(dbcfg.Engine, conn)
	if err != nil {
		return nil, fmt.Errorf("can't open the data base %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("can't do ping %v", err)
	}

	//log.Println("Connected to postgres!")
	return db, nil
}
