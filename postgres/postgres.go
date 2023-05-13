package postgres

import (
	"database/sql"
	"fmt"

	"github.com/adrianolmedo/go-restapi/config"

	_ "github.com/lib/pq"
)

type Storage struct {
	User    User
	Product Product
	Invoice Invoice
}

func NewStorage(dbcfg config.DB) (*Storage, error) {
	if dbcfg.Engine == "" {
		return nil, fmt.Errorf("database engine not selected")
	}

	if dbcfg.Engine == "postgres" {
		db, err := NewDB(dbcfg)
		if err != nil {
			return nil, fmt.Errorf("postgres: %v", err)
		}

		return &Storage{
			User:    NewUser(db),
			Product: NewProduct(db),
			Invoice: NewInvoice(db,
				NewInvoiceHeader(db),
				NewInvoiceItem(db)),
		}, nil
	}

	return nil, fmt.Errorf("database engine '%s' not implemented", dbcfg.Engine)
}

// NewDB return a postgres database connection.
func NewDB(dbcfg config.DB) (db *sql.DB, err error) {
	// postgres://user:password@host:port/dbname?sslmode=disable
	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbcfg.User, dbcfg.Password, dbcfg.Host, dbcfg.Port, dbcfg.Name)

	db, err = sql.Open(dbcfg.Engine, conn)
	if err != nil {
		return nil, fmt.Errorf("can't open the data base %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("can't do ping %v", err)
	}

	//log.Println("Connected to postgres!")
	return db, nil
}
