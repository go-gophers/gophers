package gophers

import (
	"flag"
)

var (
	vF = flag.Bool("gophers.v", false, "Log to stderr full requests and responses")
)
