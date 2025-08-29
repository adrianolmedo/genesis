package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/adrianolmedo/genesis/config"
	"github.com/adrianolmedo/genesis/internal/app"

	"github.com/peterbourgon/ff/v3"
)

func main() {
	// Pass env vars to flags.
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
	ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix())

	err := app.Run(&config.Config{
		Port: *port,
		CORS: strings.Join(strings.Fields(*cors), ""), // remove whitespaces
		Database: config.Database{
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
