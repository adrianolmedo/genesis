package sqlc

import (
	"context"
	"flag"
	"testing"
	"time"

	config "github.com/adrianolmedo/genesis"

	"github.com/jackc/pgx/v5/pgxpool"
)

// $ go test -v -tags integration -args -database-url postgres://user:password@host:port/dbname?sslmode=disable
var (
	dburl = flag.String("database-url", "", "Database URL. (example \"postgres://user:password@host:port/dbname?sslmode=disable\"")
)

// TestDB test for open and close database.
func TestDB(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	db := OpenDB(ctx, t)
	db.Close()
}

// OpenDB creates a new database connection using the provided context and test.
// It returns the connection or fails the test if an error occurs.
func OpenDB(ctx context.Context, t *testing.T) *pgxpool.Pool {
	t.Helper()
	dbcfg := config.Config{
		DatabaseURL: *dburl,
	}
	db, err := NewPool(ctx, dbcfg)
	if err != nil {
		t.Fatal(err)
	}
	return db
}
