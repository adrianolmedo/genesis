# postgres package

## Run integration tests

The `_test.go` files it have a build tag:

```go
go:build integration
+build integration
```

It also parses `package main` which calls `flag.Parse`, so all declared and visible flags will be parsed and available for the tests.

Example for run:

```bash
$ go test -v -tags integration -args -dbengine postgres -dbhost 127.0.0.1 -dbport 5432 -dbuser username -dbname testdb -dbpass 12345
```
