all: test

install:
	go install -v ./...
	go test -v ./jsons
	go test -v ./utils/...
	go test -v ./gophers/runner
	go test -v .
	# gophers-lua examples/*-lua/*.lua

install-race:
	go install -v -race ./...
	go test -v -race ./jsons
	go test -v -race ./utils/...
	go test -v -race ./gophers/runner
	go test -v -race .
	# gophers-lua examples/*-lua/*.lua

test: install
	go test -v github.com/go-gophers/gophers/examples/github-go
	gophers test --debug github.com/go-gophers/gophers/examples/placehold-go
	gophers load --debug --weighted github.com/go-gophers/gophers/examples/placehold-go

test-race: install-race
	go test -v -race github.com/go-gophers/gophers/examples/github-go
	gophers test --debug --race github.com/go-gophers/gophers/examples/placehold-go
	gophers load --debug --race --weighted github.com/go-gophers/gophers/examples/placehold-go

check: install
	go tool vet -all -shadow $(shell ls -d */ | grep -v vendor/)
	golint ./... | grep -v vendor/
	- errcheck $(shell go list ./... | grep -v vendor/)

aglio:
	# npm install -g aglio
	aglio -i examples/github-go/github.apib -t flatly -o examples/github-go/github.html

prometheus:
	cd examples/prometheus && prometheus -config.file=prometheus.yml -web.listen-address=127.0.0.1:9090
