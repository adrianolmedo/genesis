package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/adrianolmedo/go-restapi/config"
	"github.com/adrianolmedo/go-restapi/internal/app"

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

	cfg, err := config.New(*port, *cors, *dbengine, *dbhost, *dbport, *dbuser, *dbpass, *dbname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if err := app.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
