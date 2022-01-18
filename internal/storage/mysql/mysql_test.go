package mysql_test

import (
	"database/sql"
	"testing"

	"github.com/adrianolmedo/go-restapi-practice/config"
	"github.com/adrianolmedo/go-restapi-practice/internal/storage/mysql"
)

// dbcfg credentials for Postgres testing database connection.
var dbcfg = config.Database{
	Engine:   "mysql",
	Server:   "127.0.0.1",
	Port:     "3306",
	User:     "pmadrian",
	Password: "3eD5gfiqjSYO@x%k",
	Name:     "go_restapi_practice",
}

func TestDB(t *testing.T) {
	db := openDB(t)
	closeDB(t, db)
}

func openDB(t *testing.T) *sql.DB {
	db, err := mysql.New(dbcfg)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func closeDB(t *testing.T, db *sql.DB) {
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}
