package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/compose"
	"github.com/adrianolmedo/genesis/logger"
	"github.com/adrianolmedo/genesis/pgsql/sqlc"
	"github.com/adrianolmedo/genesis/rest"
	"github.com/adrianolmedo/genesis/rest/jwt"

	"github.com/joho/godotenv"
	"github.com/peterbourgon/ff/v3"
)

func main() {
	// Load .env file if exists.
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found (optional)")
	}

	// With ff we can parse flags and environment variables.
	fs := flag.NewFlagSet("main", flag.ExitOnError)
	var (
		host  = fs.String("host", ":", "Internal container IP.")
		port  = fs.String("port", "80", "Internal container port.")
		dburl = fs.String("database-url", "", "Database URL. (example \"postgres://user:password@host:port/dbname?sslmode=disable\")")
	)
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: ", err)
		os.Exit(1)
	}
	cfg := genesis.Config{
		Host:        *host,
		Port:        *port,
		DatabaseURL: *dburl,
	}
	if err := run(context.Background(), cfg); err != nil {
		fmt.Fprintln(os.Stderr, "error: ", err)
		logger.Error("run", "err", err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context, cfg genesis.Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	// Load authentication credentials.
	err := jwt.LoadFiles("app.rsa", "app.rsa.pub")
	if err != nil {
		return fmt.Errorf("certificates could not be loaded: %v", err)
	}

	// Context that is canceled upon receiving SIGINT/SIGTERM or Ctrl+c as stop signal.
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize the storage.
	db, err := sqlc.NewPool(ctx, cfg)
	if err != nil {
		return fmt.Errorf("error from storage: %v", err)
	}
	defer db.Close()
	s, err := sqlc.NewStorage(ctx, db, cfg)
	if err != nil {
		return fmt.Errorf("error from storage: %v", err)
	}

	// Initialize the server with its dependencies.
	srv := rest.Router(compose.NewServices(s))
	go func() {
		if err := srv.Listen(cfg.Host + cfg.Port); err != nil {
			logger.Error("HTTP server stopped with error", "err", err.Error())
		}
	}()

	// Wait stop signal.
	<-ctx.Done()
	fmt.Println("\nShutting down gracefully...")

	// Stop the server gracefully and close the storage (eg.: connections, workers, etc).
	if err := srv.Shutdown(); err != nil {
		return fmt.Errorf("error shutting down server: %w", err)
	}
	return nil
}
