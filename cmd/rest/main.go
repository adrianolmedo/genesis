package main

import (
	"database/sql"
	"flag"
	"log"
	"strings"

	"github.com/adrianolmedo/go-restapi-practice/auth"
	"github.com/adrianolmedo/go-restapi-practice/config"
	"github.com/adrianolmedo/go-restapi-practice/repository/mysql"
	"github.com/adrianolmedo/go-restapi-practice/repository/postgres"
	"github.com/adrianolmedo/go-restapi-practice/server/rest"
	"github.com/adrianolmedo/go-restapi-practice/user"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	MySQL      config.Driver = "mysql"
	PostgreSQL config.Driver = "postgres"
)

func main() {
	// Set up configuration.
	c := flag.String("c", "config.json", "Load configuration from `FILE`")
	flag.Parse()

	var err error
	var db *sql.DB
	var conf *config.Server
	var userRepository user.Repository // interface

	// Load authentication credentials.
	err = auth.LoadFiles("auth/app.sra", "auth/app.sra.pub")
	if err != nil {
		log.Fatalf("certificates could not be loaded: %v", err)
	}

	conf, err = config.NewServer(*c)
	if err != nil {
		log.Fatalf("error setting server: %v\n", err)
	}

	// - Set up storage from configuration loaded.
	switch conf.Engine {

	case MySQL:
		db, err = mysql.NewRepository(conf.Database)
		if err != nil {
			log.Fatalf("error from mysql storage: %v\n", err)
		}
		userRepository = mysql.NewUserRepository(db)

	case PostgreSQL:
		db, err = postgres.NewRepository(conf.Database)
		if err != nil {
			log.Fatalf("error from postgres storage: %v\n", err)
		}
		userRepository = postgres.NewUserRepository(db)

	default:
		log.Fatalf("driver not implemented %s:", conf.Engine)
	}

	// Echo framework.
	e := echo.New()

	// - Load middlewares.
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// - CORS restricted with GET, PUT, POST or DELETE method.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(conf.CORS, ","),
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// - Call handlers.
	rest.Handlers(e, userRepository)
	rest.HandlersAuthRequired(e, userRepository)

	// - Up server.
	err = e.Start(conf.Address())
	if err != nil {
		log.Printf("error server: %v\n", err)
	}
}
