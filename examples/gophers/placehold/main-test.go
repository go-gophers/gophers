// +build ignore

package main

// generated with https://github.com/go-gophers/gophers

import (
	"flag"
	"os"

	"github.com/go-gophers/gophers/gophers/runner"
	"github.com/go-gophers/gophers/utils/log"

	"github.com/go-gophers/gophers/examples/gophers/placehold"
)

func main() {
	flag.Parse()

	r := runner.New(log.New(os.Stderr, "", 0), "127.0.0.1:10311")
	r.Add("TestBasic", placehold.TestBasic, 1)

	r.Test(nil)
}
