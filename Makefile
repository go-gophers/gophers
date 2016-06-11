all: test

install:
	go install -v ./...
	go test -v ./jsons
	go test -v .
	gophers-lua examples/*-lua/*.lua

install-race:
	go install -v -race ./...
	go test -v -race ./jsons
	go test -v -race .
	gophers-lua examples/*-lua/*.lua

test: install
	go test -v github.com/go-gophers/gophers/examples/...

test-race: install-race
	go test -v -race github.com/go-gophers/gophers/examples/...

bench: install
	go test -bench=. -benchtime=10s github.com/go-gophers/gophers/examples/...

bench-race: install-race
	go test -bench=. -benchtime=10s -race github.com/go-gophers/gophers/examples/...

check: install
	go tool vet -all -shadow $(shell ls -d */ | grep -v vendor/)
	golint ./... | grep -v vendor/
	- errcheck $(shell go list ./... | grep -v vendor/)

aglio:
	# npm install -g aglio
	aglio -i examples/github-go/github.apib -t flatly -o examples/github-go/github.html
