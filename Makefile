all: test

install:
	go install -v ./...
	go test -v ./json
	go test -v .

test: install
	go test github.com/gophergala2016/gophers/examples/... -v

race:
	go install -v -race ./...
	env GORACE="halt_on_error=1" go test -v -race ./json
	env GORACE="halt_on_error=1" go test -v -race .
	env GORACE="halt_on_error=1" go test github.com/gophergala2016/gophers/examples/... -v -race

check: install
	go vet ./...
	golint ./...
	- errcheck
