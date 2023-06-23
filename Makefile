# Binary file name.
BINARY = goclisrv

# Run the Go build command to compile and generate the binary in the root.
build:
	go build -o $(BINARY) cmd/rest/*.go

# Generate the RSA files needed to create the user credentials.
genrsa:
	openssl genrsa -out app.rsa 1024
	openssl rsa -in app.rsa -pubout > app.rsa.pub

# Run Go test command.
test:
	go test ./...

# Run the Go build command to generate a debuggable binary with some IDE.
debug:
	go build -gcflags "-N -l" -o $(BINARY) .

# Run rm command to delete binary file.
clean:
	if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi