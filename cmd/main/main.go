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

	"github.com/joho/godotenv"
	"github.com/peterbourgon/ff/v3"
)

func main() {
	// Load .env file if exists.
	if err := godotenv.Load(); err != nil {
		fmt.Fprintln(os.Stderr, "No .env file found (optional)")
		os.Exit(1)
	}

	fs := flag.NewFlagSet("main", flag.ExitOnError)
	var (
		host  = fs.String("host", ":", "Internal container IP.")
		port  = fs.String("port", "80", "Internal container port.")
		dburl = fs.String("database-url", "", "Database URL. (example \"postgres://user:password@host:port/dbname?sslmode=disable\")")
	)

	// With ff we can parse flags and environment variables.
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cfg := genesis.Config{
		Host:        *host,
		Port:        *port,
		DatabaseURL: *dburl,
	}

	if err := run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		logger.Error("run", "err", err.Error())
		os.Exit(1)
	}
}

func run(cfg genesis.Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	// Load authentication credentials.
	err := jwt.LoadFiles("app.rsa", "app.rsa.pub")
	if err != nil {
		return fmt.Errorf("certificates could not be loaded: %v", err)
	}

	// Context que se cancela al recibir SIGINT/SIGTERM or Ctrl+c.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := sqlc.NewPool(ctx, cfg)
	if err != nil {
		return fmt.Errorf("error from storage: %v", err)
	}
	defer db.Close()

	s, err := sqlc.NewStorage(ctx, db, cfg)
	if err != nil {
		return fmt.Errorf("error from storage: %v", err)
	}

	// Initialize the services with the storage.
	srv := http.Router(app.NewApp(s))

	go func() {
		if err := srv.Listen(cfg.Host + cfg.Port); err != nil {
			logger.Error("HTTP server stopped with error", "err", err.Error())
		}
	}()

	// Wait stop signal.
	<-ctx.Done()
	logger.Info("Shutting down gracefully...")

	// Stop the server gracefully and close the storage (eg.: connections, workers, etc).
	if err := srv.Shutdown(); err != nil {
		return fmt.Errorf("error shutting down server: %w", err)
	}

	return nil
}
