build:
	go build cmd/rest/rest.go

genrsa:
	openssl genrsa -out app.rsa 1024
	openssl rsa -in app.rsa -pubout > app.rsa.pub

test:
	go test ./...

debug:
	go build -gcflags "-N -l" cmd/rest/rest.go

clean:
	if [ -f rest ] ; then rm rest ; fi
