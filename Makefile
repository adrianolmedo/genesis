install:
	go install cmd/rest/rest.go

build:
	go build cmd/rest/rest.go

test:
	go test ./...

debug:
	go build -gcflags "-N -l" cmd/rest/rest.go

clean:
	if [ -f rest ] ; then rm rest ; fi
