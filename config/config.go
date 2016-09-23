// Package config contains configuration for various Gophers methods.
package config

import (
	"time"
)

// Config defines configuration for various Gophers methods.
type Config struct {
	// Disable requests and responses coloring,
	NoColors bool

	// Disable requests and responses recorder.
	NoRecorders bool

	// Log to stderr full requests and responses.
	Verbose bool

	// Maximum amount of time a gophers/net.Dial will wait for a connect to complete.
	// Default is no timeout.
	DialTimeout time.Duration

	// Disable usage of IPv6.
	DisableIPv6 bool
}

// Default contain default shared configuration.
var Default = &Config{}
