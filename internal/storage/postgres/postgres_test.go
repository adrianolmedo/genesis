package postgres_test

import (
	"database/sql"
	"testing"

	"github.com/adrianolmedo/go-restapi-practice/config"
	"github.com/adrianolmedo/go-restapi-practice/internal/storage/postgres"
)

// dbcfg credentials for Postgres testing database connection.
var dbcfg = config.Database{
	Engine:   "postgres",
	Server:   "127.0.0.1",
	Port:     "5432",
	User:     "postgres",
	Password: "1234567@",
	Name:     "go_testing_restapi",
}

// TestDB can open & close.
func TestDB(t *testing.T) {
	db := openDB(t)
	closeDB(t, db)
}

func openDB(t *testing.T) *sql.DB {
	db, err := postgres.New(dbcfg)
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
