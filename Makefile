# Binary file name.
BINARY = genesis

# Generate the RSA files needed to create the user credentials.
genrsa:
	openssl genrsa -out app.rsa 1024
	openssl rsa -in app.rsa -pubout > app.rsa.pub

# Run Go test command.
test:
	go test ./...

# Generate docs package for Swagger files with swag.
# More info visit: https://github.com/swaggo/swag#getting-started
swagger:
	swag fmt -d http/
	swag init -g http/router.go

# Run the Go build command to generate a debuggable binary with some IDE.
debug:
	go build -gcflags "-N -l" -o $(BINARY) .

# Run the Go build command to compile and generate the binary in the root.
build: swagger
	go build -o $(BINARY) cmd/rest/*.go

# Run rm command to delete binary file.
clean:
	if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi
