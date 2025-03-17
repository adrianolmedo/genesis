package genesis

// Config server RESTful API.
type Config struct {
	// Host where is running the app.
	Host string

	// Port for address server, if it is empty by default it will be 80.
	Port string

	// CORS directive. Add address separated by comma. Example, "127.0.0.1,172.17.0.1".
	CORS string

	// Host where is running the database Engine.
	DBHost string

	// Port of database Engine server.
	DBPort string

	// User of database, eg.: "root".
	DBUser string

	// Password of User database.
	DBPassword string

	// Name of SQL database.
	DBName string
}
