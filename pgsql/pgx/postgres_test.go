package pgx

import (
	"context"
	"flag"
	"testing"
	"time"

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
		DBHost:     *dbhost,
		DBPort:     *dbport,
		DBUser:     *dbuser,
		DBPassword: *dbpass,
		DBName:     *dbname,
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
