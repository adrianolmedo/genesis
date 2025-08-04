include .env.development

# Binary file name.
BINARY = genesis

help:
	@echo "Commands availables:"
	@echo "  make genrsa         - Generate RSA keys for JWT auth"
	@echo "  make goose-up       - Execute migrations up"
	@echo "  make goose-down     - Revert last migration"
	@echo "  make goose-status   - Show status of migrations"
	@echo "  make goose-to       - Execute migrations up to a specific version"
	@echo "  make reset-db       - Reset the database to a clean state"
	@echo "  make test           - Execute tests"
	@echo "  make swagger        - Generate Swagger documentation"
	@echo "  make debug          - Compiling for debugging with IDE"
	@echo "  make build          - Compiling binary"
	@echo "  make clean          - Remove the binary file generated"

# Generate the RSA files needed to create the user credentials.
genrsa:
	openssl genrsa -out app.rsa 1024
	openssl rsa -in app.rsa -pubout > app.rsa.pub

goose-up:
	goose -dir migrations $(DBENGINE) "$(DBURL)" up

goose-down:
	goose -dir migrations $(DBENGINE) "$(DBURL)" down

goose-status:
	goose -dir migrations $(DBENGINE) "$(DBURL)" status

goose-to:
	goose -dir migrations $(DBENGINE) "$(DBURL)" up-to $(version)

reset-db:
	bash scripts/reset-db.sh

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

# PHONY is for commands that are not files.
# This prevents make from looking for a file with the same name as the target.
# It ensures that the command is always executed when called.
.PHONY: help genrsa goose-up goose-down goose-status goose-to test swagger debug build clean