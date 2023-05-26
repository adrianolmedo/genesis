build:
	go build -o rest .

genrsa:
	openssl genrsa -out app.rsa 1024
	openssl rsa -in app.rsa -pubout > app.rsa.pub

test:
	go test ./...

debug:
	go build -gcflags "-N -l" -o rest .

clean:
	if [ -f rest ] ; then rm rest ; fi