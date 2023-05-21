package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/adrianolmedo/go-restapi/config"
	"github.com/adrianolmedo/go-restapi/postgres"
	"github.com/adrianolmedo/go-restapi/rest"
	"github.com/adrianolmedo/go-restapi/rest/jwt"

	"github.com/peterbourgon/ff/v3"
)

func main() {
	fs := flag.NewFlagSet("rest", flag.ExitOnError)
	var (
		port     = fs.String("port", "80", "Internal container port.")
		cors     = fs.String("cors", "", "CORS directive, write address separated by comma.")
		dbhost   = fs.String("dbhost", "", "Database host.")
		dbengine = fs.String("dbengine", "", "Database engine, choose mysql or postgres.")
		dbport   = fs.String("dbport", "", "Database port.")
		dbuser   = fs.String("dbuser", "", "Database user.")
		dbpass   = fs.String("dbpass", "", "Database password.")
		dbname   = fs.String("dbname", "", "Database name.")
	)

	// Pass env vars to flags.
	ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix())

	err := run(&config.Config{
		Port: *port,
		CORS: *cors,
		DB: config.DB{
			Engine:   *dbengine,
			Host:     *dbhost,
			Port:     *dbport,
			User:     *dbuser,
			Password: *dbpass,
			Name:     *dbname,
		},
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(cfg *config.Config) error {
	// Load authentication credentials.
	err := jwt.LoadFiles("app.sra", "app.sra.pub")
	if err != nil {
		return fmt.Errorf("certificates could not be loaded: %v", err)
	}

	strg, err := postgres.NewStorage(cfg.DB)
	if err != nil {
		return fmt.Errorf("error from storage: %v", err)
	}

	return rest.Routes(strg).Listen(":" + cfg.Port)
}
