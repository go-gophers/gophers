all: test

install:
	go install -v ./...
	go test -v ./jsons
	go test -v .

install-race:
	go install -v -race ./...
	go test -v -race ./jsons
	go test -v -race .
	gophers examples/*-lua/*.lua

test: install
	go test github.com/go-gophers/gophers/examples/... -v

test-race: install-race
	go test github.com/go-gophers/gophers/examples/... -v -race

check: install
	go tool vet -all -shadow $(shell ls -d */ | grep -v vendor/)
	golint ./... | grep -v vendor/
	- errcheck $(shell go list ./... | grep -v vendor/)

aglio:
	# npm install -g aglio
	aglio -i examples/github-go/github.apib -t flatly -o examples/github-go/github.html
