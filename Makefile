PACKAGES := $(shell go list ./... | grep -v vendor/ | grep -v examples/)
EXAMPLES := $(shell go list ./... | grep examples/)

export CGO_ENABLED := 0
export GODEBUG := netdns=go

all: test

init:
	go get -u github.com/kisielk/errcheck github.com/golang/lint/golint

install:
	go install -v $(PACKAGES)
	go test -v $(PACKAGES)

install-race:
	go install -v -race $(PACKAGES)
	go test -v -race $(PACKAGES)

# test: install
# 	go test -v github.com/go-gophers/gophers/examples/github-go
# 	gophers test --debug github.com/go-gophers/gophers/examples/placehold-go
# 	gophers load --debug --weighted github.com/go-gophers/gophers/examples/placehold-go

# test-race: install-race
# 	go test -v -race github.com/go-gophers/gophers/examples/github-go
# 	gophers test --debug --race github.com/go-gophers/gophers/examples/placehold-go
# 	gophers load --debug --race --weighted github.com/go-gophers/gophers/examples/placehold-go

check: install
	go vet $(PACKAGES)
	- errcheck $(PACKAGES)
	for package in $(PACKAGES) ; do \
		golint $$package ; \
	done

aglio:
	# npm install -g aglio
	aglio -i examples/github-go/github.apib -t flatly -o examples/github-go/github.html

prometheus:
	cd examples/prometheus && prometheus -config.file=prometheus.yml -web.listen-address=127.0.0.1:9090
