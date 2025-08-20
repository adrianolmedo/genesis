//go:build integration
// +build integration

package pgx

import (
	"context"
	"flag"
	"testing"
	"time"

	config "github.com/adrianolmedo/genesis"

	"github.com/jackc/pgx/v5"
)

// $ go test -v -tags integration -args -database-url postgres://user:password@host:port/dbname?sslmode=disable
var (
	dburl = flag.String("database-url", "", "Database URL. (example \"postgres://user:password@host:port/dbname?sslmode=disable\"")
)

// TestDB test for open & close database.
func TestDB(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	conn := openDB(ctx, t)
	closeDB(ctx, t, conn)
}

// openDB creates a new database connection using the provided context and test.
// It returns the connection or fails the test if an error occurs.
func openDB(ctx context.Context, t *testing.T) *pgx.Conn {
	t.Helper()

	dbcfg := config.Config{
		DatabaseURL: *dburl,
	}

	conn, err := newDB(ctx, dbcfg)
	if err != nil {
		t.Fatal(err)
	}

	return conn
}

// closeDB closes the database connection and fails the test if an error occurs.
// It is a helper function to ensure proper cleanup after tests.
// It should be deferred after opening a connection.
func closeDB(ctx context.Context, t *testing.T, db *pgx.Conn) {
	t.Helper()

	if err := db.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
}

// testCtx creates a context with a timeout for testing purposes.
func testCtx(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)
	return ctx
}
