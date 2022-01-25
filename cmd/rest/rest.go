package main

import (
	"flag"
	"log"

	"github.com/adrianolmedo/go-restapi-practice/config"
	"github.com/adrianolmedo/go-restapi-practice/internal/app"
)

func main() {
	// Pass env vars to flags
	port := flag.String("port", "80", "Internal container port.")
	cors := flag.String("cors", "", "CORS directive, write address separted by comma.")
	dbhost := flag.String("dbhost", "", "Database host.")
	dbengine := flag.String("dbengine", "", "Database engine, choose mysql or postgres.")
	dbport := flag.String("dbport", "", "Database port.")
	dbuser := flag.String("dbuser", "", "Database user.")
	dbpass := flag.String("dbpass", "", "Database password.")
	dbname := flag.String("dbname", "", "Database name.")
	flag.Parse()

	cfg, err := config.New(*port, *cors, *dbhost, *dbengine, *dbport, *dbuser, *dbpass, *dbname)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
