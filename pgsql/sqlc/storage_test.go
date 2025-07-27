//go:build integration
// +build integration

package sqlc

import (
	"context"
	"flag"
	"testing"
	"time"

	config "github.com/adrianolmedo/genesis"

	"github.com/jackc/pgx/v5/pgxpool"
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	db := openDB(ctx, t)
	db.Close()
}

// openDB creates a new database connection using the provided context and test.
// It returns the connection or fails the test if an error occurs.
func openDB(ctx context.Context, t *testing.T) *pgxpool.Pool {
	t.Helper()

	dbcfg := config.Config{
		DBHost:     *dbhost,
		DBPort:     *dbport,
		DBUser:     *dbuser,
		DBPassword: *dbpass,
		DBName:     *dbname,
	}

	db, err := newPool(ctx, dbcfg)
	if err != nil {
		t.Fatal(err)
	}

	return db
}
