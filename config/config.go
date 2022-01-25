package config

import (
	"encoding/json"
	"errors"
	"strings"
)

var ErrDatabaseCantBeEmpty = errors.New("database fields configuration can't be empty")

// Config server RESTful API.
type Config struct {
	// Port for address server, if is empty by default will are 80.
	Port string `json:"port"`

	// CORS directive. Add address separated by comma.
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

	for _, v := range cfg.Database.toMap() {
		if v == "" {
			return nil, ErrDatabaseCantBeEmpty
		}
	}

	return &cfg, nil
}

// toMap convert Database struct to simple map.
func (d Database) toMap() map[string]interface{} {
	var m map[string]interface{}
	b, _ := json.Marshal(d)
	json.Unmarshal(b, &m)
	return m
}
