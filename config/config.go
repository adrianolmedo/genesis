package config

import (
	"strings"
)

// Config server RESTful API.
type Config struct {
	// Port for address server, if is empty by default will are 80.
	Port string `json:"port"`

	// CORS directive. Add address separated by comma. Example, "127.0.0.1,172.17.0.1".
	CORS string `json:"cors"`

	Database `json:"database"`
}

// Database config.
type Database struct {
	// Engine eg.: "mysql" or "postgres".
	Engine string `json:"engine"`

	// Host when is running the database Engine.
	Host string `json:"host"`

	// Port of database Engine server.
	Port string `json:"port"`

	// User of database, eg.: "root".
	User string `json:"user"`

	// Password of User database
	Password string `json:"password"`

	// Name of SQL database.
	Name string `json:"name"`
}

func New(port, cors, dbengine, dbhost, dbport, dbuser, dbpass, dbname string) (*Config, error) {
	cfg := Config{
		Port: port,
		CORS: strings.Join(strings.Fields(cors), ""), // remove whitespaces
		Database: Database{
			Engine:   dbengine,
			Host:     dbhost,
			Port:     dbport,
			User:     dbuser,
			Password: dbpass,
			Name:     dbname,
		},
	}

	return &cfg, nil
}
