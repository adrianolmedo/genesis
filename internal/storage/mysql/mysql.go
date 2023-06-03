package mysql

import (
	"database/sql"
	"fmt"

	"github.com/adrianolmedo/genesis/config"

	_ "github.com/go-sql-driver/mysql"
)

func New(dbcfg config.Database) (db *sql.DB, err error) {
	// En los parámetros de conexión para mysql vamos a establecer ParseTime (bool),
	// ya que sin eso no se podría hacer el maping del campo CreatedAt a time.Time:
	//
	// user:password@tcp(host:port)/dbname?tls=false&autocommit=true&parseTime=true
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=false&autocommit=true&parseTime=true",
		dbcfg.User, dbcfg.Password, dbcfg.Host, dbcfg.Port, dbcfg.Name)

	db, err = sql.Open(dbcfg.Engine, conn)
	if err != nil {
		return nil, fmt.Errorf("can't open the data base %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("can't do ping %v", err)
	}

	//log.Println("Connected to mysql!")
	return db, nil
}
