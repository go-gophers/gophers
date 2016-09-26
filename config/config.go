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
	DialTimeout time.Duration

	// Disable usage of IPv6.
	DisableIPv6 bool

	// HTTP listen address for Prometheus, etc.
	HTTPAddr string
}

var (
	// Default contain default shared configuration.
	Default = &Config{
		HTTPAddr: "127.0.0.1:10311",
	}

	// Flags contains command-line flags for Default.
	Flags *flag.FlagSet
)

func init() {
	Flags = flag.NewFlagSet("gophers", flag.ContinueOnError)
	Flags.BoolVar(&Default.NoColors, "no-color", Default.NoColors, "Disable requests and responses coloring.")
	Flags.BoolVar(&Default.NoRecorders, "no-recorders", Default.NoRecorders, "Disable requests and responses recorders.")
	Flags.BoolVar(&Default.Verbose, "verbose", Default.Verbose, "Log to stderr full requests and responses.")
	Flags.DurationVar(&Default.DialTimeout, "dial-timeout", Default.DialTimeout, "Maximum amount of time a gophers/net.Dial will wait for a connect to complete.")
	Flags.BoolVar(&Default.DisableIPv6, "disable-ipv6", Default.DisableIPv6, "Disable usage of IPv6.")
	Flags.StringVar(&Default.HTTPAddr, "http-addr", Default.HTTPAddr, "HTTP listen address for Prometheus, etc.")
}
