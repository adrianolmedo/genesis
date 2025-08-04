package genesis

import "fmt"

// Config server RESTful API.
type Config struct {
	// Host where is running the app.
	Host string

	// Port for address server, if it is empty by default it will be 80.
	Port string

	// DBURL is the database URL connection string.
	// Example: "postgres://user:password@host:port/dbname?sslmode=disable".
	DBURL string
}

// Validate checks if the configuration is valid.
func (c Config) Validate() error {
	if c.Host == "" || c.Port == "" {
		return fmt.Errorf("host and port must be specified")
	}

	if c.DBURL == "" {
		return fmt.Errorf("database URL is required")
	}

	return nil
}
