package mysql

import (
	"database/sql"
	"fmt"
	"log"

	"go-restapi-practice/config"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func NewStorage(database config.Database) (*sql.DB, error) {
	var err error

	// En los parámetros de conexión para mysql vamos a establecer ParseTime (bool),
	// ya que sin eso no se podría hacer el maping del campo CreatedAt a time.Time:
	// user:password@tcp(server:port)/dbname?tls=false&autocommit=true&parseTime=true
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=false&autocommit=true&parseTime=true",
		database.User, database.Password, database.Server, database.Port, database.Name)

	db, err = sql.Open(string(database.Engine), conn)
	if err != nil {
		return nil, fmt.Errorf("can't open the data base, %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("can't do ping, %v", err)
	}

	log.Println("connected to mysql!")
	return db, nil
}
