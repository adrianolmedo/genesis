#!/bin/sh

for d in $(go list ./...); do
	echo "testing package $d"
	go test -v $d
done