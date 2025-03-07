//go:build integration
// +build integration

package pq

import (
	"database/sql"
	"flag"
	"testing"

	config "github.com/adrianolmedo/genesis"
)

// $ go test -v -tags integration -args -dbengine postgres -dbhost 127.0.0.1 -dbport 5432 -dbuser username -dbname foodb -dbpass 12345
var (
	dbhost = flag.String("dbhost", "", "Database host.")
	dbport = flag.String("dbport", "", "Database port.")
	dbuser = flag.String("dbuser", "", "Database user.")
	dbpass = flag.String("dbpass", "", "Database password.")
	dbname = flag.String("dbname", "", "Database name.")
)

// TestDB test for open & close database.
func TestDB(t *testing.T) {
	db := openDB(t)
	closeDB(t, db)
}

func openDB(t *testing.T) *sql.DB {
	dbcfg := config.DB{
		Host:     *dbhost,
		Port:     *dbport,
		User:     *dbuser,
		Password: *dbpass,
		Name:     *dbname,
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
