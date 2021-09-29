package main

import (
	"database/sql"
	"flag"
	"log"
	"strings"

	"go-rest-practice/auth"
	"go-rest-practice/config"
	"go-rest-practice/server/rest"
	"go-rest-practice/service"
	"go-rest-practice/storage/mysql"
	"go-rest-practice/storage/postgres"

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
	var srv *config.Server
	var user service.UserDAO // interface

	// Load authentication credentials.
	err = auth.LoadFiles("auth/app.sra", "auth/app.sra.pub")
	if err != nil {
		log.Fatalf("certificates could not be loaded: %v", err)
	}

	srv, err = config.NewServer(*c)
	if err != nil {
		log.Fatalf("error setting server: %v\n", err)
	}

	// - Set up storage from configuration loaded.
	switch srv.Engine {
	case MySQL:
		db, err = mysql.NewStorage(srv.Database)
		if err != nil {
			log.Fatalf("error from mysql storage: %v\n", err)
		}
		user = mysql.NewUserDAO(db)
	case PostgreSQL:
		db, err = postgres.NewStorage(srv.Database)
		if err != nil {
			log.Fatalf("error from postgres storage: %v\n", err)
		}
		user = postgres.NewUserDAO(db)
	default:
		log.Fatalf("driver not implemented %s:", srv.Engine)
	}

	// Echo framework.
	e := echo.New()

	// - Load middlewares.
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// - CORS restricted with GET, PUT, POST or DELETE method.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(srv.CORS, ","),
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// - Call handlers.
	rest.Handlers(e, user)
	rest.HandlersAuthRequired(e, user)

	// - Up server.
	err = e.Start(srv.Address())
	if err != nil {
		log.Printf("error server: %v\n", err)
	}
}
