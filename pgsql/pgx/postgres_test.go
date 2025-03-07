//go:build integration
// +build integration

package pgx

import (
	"context"
	"flag"
	"testing"

	config "github.com/adrianolmedo/genesis"

	"github.com/jackc/pgx/v5"
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

func openDB(t *testing.T) *pgx.Conn {
	dbcfg := config.DB{
		Host:     *dbhost,
		Port:     *dbport,
		User:     *dbuser,
		Password: *dbpass,
		Name:     *dbname,
	}

	conn, err := newDB(dbcfg)
	if err != nil {
		t.Fatal(err)
	}

	return conn
}

func closeDB(t *testing.T, db *pgx.Conn) {
	if err := db.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
}
