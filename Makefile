PACKAGES := $(shell go list ./... | grep -v vendor/ | grep -v examples/)
EXAMPLES_TESTING := $(shell go list ./... | grep examples/testing/)

export GODEBUG := netdns=go

all: check

init:
	go get -u github.com/kisielk/errcheck github.com/golang/lint/golint

install:
	go install -v $(PACKAGES)
	cd examples/gophers/placehold && gophers test --debug
	cd examples/gophers/placehold && gophers load --debug
	go test -v $(PACKAGES)
	cd examples/gophers/placehold && go build -o /dev/null -v main-test.go
	cd examples/gophers/placehold && go build -o /dev/null -v main-load.go

install-race:
	go install -v -race $(PACKAGES)
	go test -v -race $(PACKAGES)

test-testing: install
	go test -v $(EXAMPLES_TESTING)

test-testing-race: install-race
	go test -v -race $(EXAMPLES_TESTING)

test-gophers: install
	cd examples/gophers/placehold && go run -v main-test.go -disable-ipv6

load-gophers: install
	cd examples/gophers/placehold && go run -v main-load.go -disable-ipv6

# test-gophers-race: install-race
# 	gophers test --debug $(EXAMPLES_GOPHERS)

check: install
	go vet $(PACKAGES)
	- errcheck $(PACKAGES)
	for package in $(PACKAGES) ; do \
		golint $$package ; \
	done

aglio:
	# npm install -g aglio
	aglio -i examples/testing/github/github.apib -t flatly -o examples/testing/github/github.html

prometheus:
	cd examples/prometheus && prometheus -config.file=prometheus.yml -web.listen-address=127.0.0.1:9090
