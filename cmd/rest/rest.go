package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/app"
	"github.com/adrianolmedo/genesis/http"
	"github.com/adrianolmedo/genesis/http/jwt"
	"github.com/adrianolmedo/genesis/logger"
	"github.com/adrianolmedo/genesis/pgsql/sqlc"

	"github.com/peterbourgon/ff/v3"
)

func main() {
	fs := flag.NewFlagSet("rest", flag.ExitOnError)
	var (
		host   = fs.String("host", ":", "Internal container IP.")
		port   = fs.String("port", "80", "Internal container port.")
		cors   = fs.String("cors", "", "CORS directive, write address separated by comma.")
		dbhost = fs.String("dbhost", "", "Database host.")
		dbport = fs.String("dbport", "", "Database port.")
		dbuser = fs.String("dbuser", "", "Database user.")
		dbpass = fs.String("dbpass", "", "Database password.")
		dbname = fs.String("dbname", "", "Database name.")
	)

	// Pass env vars to flags.
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	cfg := genesis.Config{
		Host:       *host,
		Port:       *port,
		CORS:       *cors,
		DBHost:     *dbhost,
		DBPort:     *dbport,
		DBUser:     *dbuser,
		DBPassword: *dbpass,
		DBName:     *dbname,
	}

	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		logger.Error("run", "err", err.Error())
		os.Exit(1)
	}
}

func run(cfg genesis.Config) error {
	// Load authentication credentials.
	err := jwt.LoadFiles("app.rsa", "app.rsa.pub")
	if err != nil {
		return fmt.Errorf("certificates could not be loaded: %v", err)
	}

	// Context que se cancela al recibir SIGINT/SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	s, err := sqlc.NewStorage(ctx, cfg)
	if err != nil {
		return fmt.Errorf("error from storage: %v", err)
	}

	return http.Router(app.NewServices(s)).Listen(cfg.Host + cfg.Port)
}
