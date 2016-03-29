all: test

install:
	go install -v ./...
	go test -v ./jsons
	go test -v .

test: install
	gophers examples/*-lua/*.lua
	go test github.com/go-gophers/gophers/examples/... -v

race:
	go install -v -race ./...
	env GORACE="halt_on_error=1" go test -v -race ./jsons
	env GORACE="halt_on_error=1" go test -v -race .
	env GORACE="halt_on_error=1" go test github.com/go-gophers/gophers/examples/... -v -race

check: install
	go tool vet -all -shadow $(shell ls -d */ | grep -v vendor/)
	golint ./... | grep -v vendor/
	- errcheck $(shell go list ./... | grep -v vendor/)

aglio:
	# npm install -g aglio
	aglio -i examples/github-go/github.apib -t flatly -o examples/github-go/github.html
