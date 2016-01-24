package github

import (
	"net/url"
	"os"

	"github.com/gophergala2016/gophers"
)

var (
	TestPrefix = "test-gophers-"

	Login  string
	Client *gophers.Client
)

func init() {
	urlStr := "https://api.github.com/?access_token=" + os.Getenv("GOPHERS_GITHUB_TOKEN")
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	Client = gophers.NewClient(*u)
}
