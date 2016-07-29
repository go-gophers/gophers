package gophers

import (
	"flag"
	"sync"

	"github.com/fatih/color"
)

var (
	noColorsF        = flag.Bool("gophers.no-colors", false, "Disable requests and responses coloring")
	vF               = flag.Bool("gophers.v", false, "Log to stderr full requests and responses")
	disableRecorderF = flag.Bool("gophers.disable-recorder", false, "Disable requests and responses recorder")
)

var initColorOnce sync.Once

func initColor() {
	initColorOnce.Do(func() {
		if *noColorsF {
			color.NoColor = true
		}
	})
}
