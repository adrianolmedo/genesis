//go:build integration
// +build integration

package mysql_test

import (
	"database/sql"
	"flag"
	"testing"

	"github.com/adrianolmedo/go-restapi/config"
	"github.com/adrianolmedo/go-restapi/internal/storage/mysql"
)

// $ go test -v -tags integration -args -dbengine mysql -dbhost 127.0.0.1 -dbport 3306 -dbuser username -dbname foodb -dbpass 12345
var (
	dbhost   = flag.String("dbhost", "", "Database host.")
	dbengine = flag.String("dbengine", "", "Database engine, choose mysql or postgres.")
	dbport   = flag.String("dbport", "", "Database port.")
	dbuser   = flag.String("dbuser", "", "Database user.")
	dbpass   = flag.String("dbpass", "", "Database password.")
	dbname   = flag.String("dbname", "", "Database name.")
)

func TestDB(t *testing.T) {
	db := openDB(t)
	closeDB(t, db)
}

func openDB(t *testing.T) *sql.DB {
	dbcfg := config.Database{
		Engine:   *dbengine,
		Host:     *dbhost,
		Port:     *dbport,
		User:     *dbuser,
		Password: *dbpass,
		Name:     *dbname,
	}

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
