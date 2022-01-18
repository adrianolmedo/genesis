package postgres

import (
	"database/sql"
	"fmt"

	"github.com/adrianolmedo/go-restapi-practice/config"

	_ "github.com/lib/pq"
)

func New(dbcfg config.Database) (db *sql.DB, err error) {
	// postgres://user:password@server:port/dbname?sslmode=disable
	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbcfg.User, dbcfg.Password, dbcfg.Server, dbcfg.Port, dbcfg.Name)

	db, err = sql.Open(string(dbcfg.Engine), conn)
	if err != nil {
		return nil, fmt.Errorf("can't open the data base, %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("can't do ping, %v", err)
	}

	//log.Println("Connected to postgres!")
	return db, nil
}
