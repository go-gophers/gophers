package github

import (
	"net/url"
	"os"

	"github.com/go-gophers/gophers"
)

var (
	TestPrefix = "test-gophers-"

	Login  string
	Client *gophers.Client
)

func init() {
	token := os.Getenv("GOPHERS_GITHUB_TOKEN")
	if token == "" {
		msg := "To run tests you should first get persoinal github.com token here: https://github.com/settings/tokens\n" +
			"Set it to environment variable GOPHERS_GITHUB_TOKEN."
		panic(msg)
	}

	urlStr := "https://api.github.com/?access_token=" + token
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	Client = gophers.NewClient(*u)
	Client.DefaultHeaders.Set("Content-Type", "application/json")
}
