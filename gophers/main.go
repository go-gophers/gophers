package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/yuin/gopher-lua"

	"github.com/gophergala2016/gophers/glua"
)

func main() {
	log.SetFlags(0)
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
