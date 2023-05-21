build:
	go build -o main .

test:
	go test ./...

debug:
	go build -gcflags "-N -l" -o rest .

clean:
	if [ -f main ] ; then rm main ; fi