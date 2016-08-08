// Package github contains Gophers examples for github.com to be used with go test tool.
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
		msg := "To run github tests you should first get personal github.com token here: https://github.com/settings/tokens\n" +
			"Required permissions: read:org, public_repo, delete_repo. \n" +
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
