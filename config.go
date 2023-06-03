package genesis

// Config server RESTful API.
type Config struct {
	// Port for address server, if it is empty by default it will be 80.
	Port string

	// CORS directive. Add address separated by comma. Example, "127.0.0.1,172.17.0.1".
	CORS string

	// Database.
	DB
}

// DB Database config.
type DB struct {
	// Engine eg.: "mysql" or "postgres".
	Engine string

	// Host when is running the database Engine.
	Host string

	// Port of database Engine server.
	Port string

	// User of database, eg.: "root".
	User string

	// Password of User database.
	Password string

	// Name of SQL database.
	Name string
}
