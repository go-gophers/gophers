// Package config contains configuration for various Gophers methods.
package config

import (
	"flag"
	"time"
)

// Config defines configuration for various Gophers methods.
type Config struct {
	// Disable requests and responses coloring.
	NoColors bool

	// Disable requests and responses recorders.
	NoRecorders bool

	// Log to stderr full requests and responses.
	Verbose bool

	// Maximum amount of time a gophers/net.Dial will wait for a connect to complete.
	// Default is no timeout.
	DialTimeout time.Duration

	// Disable usage of IPv6.
	DisableIPv6 bool
}

var (
	// Default contain default shared configuration.
	Default = &Config{}

	// Flags contains command-line flags for Default.
	Flags *flag.FlagSet
)

func init() {
	Flags = flag.NewFlagSet("gophers", flag.ContinueOnError)
	Flags.BoolVar(&Default.NoColors, "no-color", false, "Disable requests and responses coloring.")
	Flags.BoolVar(&Default.NoRecorders, "no-recorders", false, "Disable requests and responses recorders.")
	Flags.BoolVar(&Default.Verbose, "verbose", false, "Log to stderr full requests and responses.")
	Flags.DurationVar(&Default.DialTimeout, "dial-timeout", 0, "Maximum amount of time a gophers/net.Dial will wait for a connect to complete.")
	Flags.BoolVar(&Default.DisableIPv6, "disable-ipv6", false, "Disable usage of IPv6.")
}
