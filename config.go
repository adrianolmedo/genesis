package genesis

import "fmt"

// Config holds the application configuration.
type Config struct {
	// Host where is running the app.
	Host string

	// Port for address server, if it is empty by default it will be 80.
	Port string

	// DatabaseURL is the database connection string.
	// Example: "postgres://user:password@host:port/dbname?sslmode=disable".
	DatabaseURL string
}

// Validate checks if the configuration is valid.
func (c Config) Validate() error {
	if c.Host == "" || c.Port == "" {
		return fmt.Errorf("host and port must be specified")
	}

	if c.DatabaseURL == "" {
		return fmt.Errorf("database URL is required")
	}

	return nil
}
