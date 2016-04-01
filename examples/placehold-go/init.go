package placehold

import (
	_ "image/jpeg"
	_ "image/png"
	"net/url"

	"github.com/go-gophers/gophers"
)

var (
	Client *gophers.Client
)

func init() {
	u, err := url.Parse("http://placehold.it/")
	if err != nil {
		panic(err)
	}
	Client = gophers.NewClient(*u)
}
