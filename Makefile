all: test

install:
	go install -v ./...
	go test -v ./jsons
	go test -v .

test: install
	gophers examples/*-lua/*.lua
	go test github.com/gophergala2016/gophers/examples/... -v

race:
	go install -v -race ./...
	env GORACE="halt_on_error=1" go test -v -race ./jsons
	env GORACE="halt_on_error=1" go test -v -race .
	env GORACE="halt_on_error=1" go test github.com/gophergala2016/gophers/examples/... -v -race

check: install
	go vet ./...
	golint ./...
	- errcheck

aglio:
	# npm install -g aglio
	aglio -i examples/github-go/github.apib -t flatly -o examples/github-go/github.html
