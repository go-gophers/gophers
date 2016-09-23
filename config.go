package gophers

import (
	"sync"

	"github.com/fatih/color"
)

// Config defines configuration for various methods.
type Config struct {
	// Disable requests and responses coloring,
	NoColors bool

	// Disable requests and responses recorder.
	NoRecorders bool

	// Log to stderr full requests and responses.
	Verbose bool
}

// DefaultConfig contain default shared configuration.
var DefaultConfig = &Config{}

var initColorsOnce sync.Once

func initColors() {
	initColorsOnce.Do(func() {
		if DefaultConfig.NoColors {
			color.NoColor = true
		}
	})
}
