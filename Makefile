all: test

install:
	go install -v .
	go test -v ./json
	go test -v .

test: install
	go test github.com/gophergala2016/gophers/examples/... -v

check: install
	go vet ./...
	golint ./...
	- errcheck
