package gophers

import (
	"flag"

	"github.com/fatih/color"
)

var (
	noColorsF = flag.Bool("gophers.no-colors", false, "Disable requests and responses coloring")
	vF        = flag.Bool("gophers.v", false, "Log to stderr full requests and responses")
)

func init() {
	if *noColorsF {
		color.NoColor = true
	}
}
