package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/adrianolmedo/go-restapi-practice/config"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

func NewRepository(database config.Database) (*sql.DB, error) {
	var err error

	// postgres://user:password@server:port/dbname?sslmode=disable
	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		database.User, database.Password, database.Server, database.Port, database.Name)

	db, err = sql.Open(string(database.Engine), conn)
	if err != nil {
		return nil, fmt.Errorf("can't open the data base, %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("can't do ping, %v", err)
	}

	log.Println("connected to postgres!")
	return db, nil
}
