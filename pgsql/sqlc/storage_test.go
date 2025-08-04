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

// $ go test -v -tags integration -args -dburl postgres://user:password@host:port/dbname?sslmode=disable
var (
	dburl = flag.String("dburl", "", "Database URL. (example \"postgres://user:password@host:port/dbname?sslmode=disable\"")
)

// TestDB test for open and close database.
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
		DBURL: *dburl,
	}

	db, err := newPool(ctx, dbcfg)
	if err != nil {
		t.Fatal(err)
	}

	return db
}
