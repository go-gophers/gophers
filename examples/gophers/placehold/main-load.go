// +build ignore

package main

// generated with https://github.com/go-gophers/gophers

import (
	"flag"
	"os"
	"time"

	"github.com/go-gophers/gophers/config"
	"github.com/go-gophers/gophers/gophers/runner"
	"github.com/go-gophers/gophers/utils/log"

	"github.com/go-gophers/gophers/examples/gophers/placehold"
)

func main() {
	config.Flags.VisitAll(func(f *flag.Flag) {
		flag.Var(f.Value, f.Name, f.Usage)
	})
	flag.Parse()

	exitCode := 2
	defer func() {
		if p := recover(); p != nil {
			panic(p)
		}

		os.Exit(exitCode)
	}()

	r := runner.New("", log.New(os.Stderr, "", 0))
	r.Add("TestBasic", placehold.TestBasic, 1)

	l, err := runner.NewStepLoader(5, 10, 1, 1*time.Second)
	if err != nil {
		panic(err)
	}

	exitCode = r.Load(nil, l, runner.FailStep)
}
