// +build ignore

package main

// generated with https://github.com/go-gophers/gophers

import (
	"flag"
	"os"
	"time"

	"github.com/go-gophers/gophers/gophers/runner"
	"github.com/go-gophers/gophers/utils/log"

	"github.com/go-gophers/gophers/examples/placehold-go"
)

func main() {
	flag.Parse()

	r := runner.New(log.New(os.Stderr, "", 0), "127.0.0.1:10311")
	r.Add("TestBasic", placehold.TestBasic, 1)
	r.Add("TestFail", placehold.TestFail, 1)
	r.Add("TestPanic", placehold.TestPanic, 1)

	l, err := runner.NewStepLoader(5, 10, 1, 1*time.Second)
	if err != nil {
		panic(err)
	}

	r.LoadWeighted(l)
}
