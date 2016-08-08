package main

import (
	"flag"
	"net/http"

	"github.com/yuin/gopher-lua"

	"github.com/go-gophers/gophers/exp/glua"
	"github.com/go-gophers/gophers/utils/log"
)

func main() {
	flag.Parse()

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("gophers", glua.NewClient(&http.Client{}).Loader)

	for _, arg := range flag.Args() {
		err := L.DoFile(arg)
		if err != nil {
			log.Fatal(err)
		}
	}
}
