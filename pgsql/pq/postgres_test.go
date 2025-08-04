//go:build integration
// +build integration

package pq

import (
	"database/sql"
	"flag"
	"testing"

	config "github.com/adrianolmedo/genesis"
)

// $ go test -v -tags integration -args -dburl postgres://user:password@host:port/dbname?sslmode=disable
var (
	dburl = flag.String("dburl", "", "Database URL. (example \"postgres://user:password@host:port/dbname?sslmode=disable\"")
)

// TestDB test for open & close database.
func TestDB(t *testing.T) {
	db := openDB(t)
	closeDB(t, db)
}

func openDB(t *testing.T) *sql.DB {
	dbcfg := config.Config{
		DBURL: *dburl,
	}

	db, err := newDB(dbcfg)
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
