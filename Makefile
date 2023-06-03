build:
	go build -o cligosrv cmd/rest/*.go

genrsa:
	openssl genrsa -out app.rsa 1024
	openssl rsa -in app.rsa -pubout > app.rsa.pub

test:
	go test ./...

debug:
	go build -gcflags "-N -l" -o cligosrv .

clean:
	if [ -f cligosrv ] ; then rm cligosrv ; fi