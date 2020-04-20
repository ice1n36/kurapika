.PHONY: gen lint test

VERSION := `git vertag get`
COMMIT  := `git rev-parse HEAD`

gen:
	go generate ./...

lint: gen
	golangci-lint run

test: lint
	go test -v --race ./...
