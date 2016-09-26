// Package placehold contains Gophers examples for placehold.it to be used with gophers tool.
package placehold

import (
	_ "image/jpeg" // register decoder
	_ "image/png"  // register decoder
	"net/http"
	"net/url"
	"time"

	"github.com/go-gophers/gophers"
	"github.com/go-gophers/gophers/config"
	"github.com/go-gophers/gophers/net"
)

var (
	Client *gophers.Client
)

func init() {
	config.Default.DisableIPv6 = true

	u, err := url.Parse("http://placehold.it/")
	if err != nil {
		panic(err)
	}
	Client = gophers.NewClient(*u)
	Client.HTTPClient = &http.Client{
		Transport: &http.Transport{
			Dial:                net.Dial,
			MaxIdleConnsPerHost: 1000,
		},
		Timeout: 10 * time.Second,
	}
}
