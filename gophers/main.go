// Package gophers implements gophers tool.
package main

import (
	"github.com/go-gophers/gophers/gophers/cmd"
	"github.com/go-gophers/gophers/utils/log"
)

func main() {
	log.Default.SetPrefix("gophers: ")

	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
