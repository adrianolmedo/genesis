# Binary file name.
BINARY = genesis

# Generate the RSA files needed to create the user credentials.
genrsa:
	openssl genrsa -out app.rsa 1024
	openssl rsa -in app.rsa -pubout > app.rsa.pub

# Execute Go test command.
test:
	sqlc generate
	go test ./...

# Generate docs package for Swagger files with swag.
# More info visit: https://github.com/swaggo/swag#getting-started
swagger:
	swag fmt -d http/
	swag init -g http/router.go

# Execute the Go build command to generate a debuggable binary with some IDE.
debug:
	sqlc generate
	go build -gcflags "-N -l" -o $(BINARY) .

# Execute the Go build command to compile and generate the binary in the root.
build: swagger
	sqlc generate
	go build -o $(BINARY) cmd/rest/*.go

# Execute rm command to delete binary file.
clean:
	if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi
