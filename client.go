package gophers

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

const (
	defaultUserAgent = "github.com/gophergala2016/gophers"
)

type Client struct {
	Base           url.URL
	HTTPClient     *http.Client
	DefaultHeaders http.Header
}

func NewClient(base url.URL) *Client {
	return &Client{
		Base:           base,
		HTTPClient:     http.DefaultClient,
		DefaultHeaders: http.Header{},
	}
}

func (c *Client) NewRequest(t testing.TB, method string, urlStr string) *http.Request {
	suffix, err := url.Parse(urlStr)
	if err != nil {
		t.Fatal(err)
	}

	// add path, check for '//'
	u := c.Base
	if strings.HasSuffix(u.Path, "/") && strings.HasPrefix(suffix.Path, "/") {
		suffix.Path = strings.TrimPrefix(suffix.Path, "/")
	}
	u.Path += suffix.Path

	// add query
	q := u.Query()
	for k, vs := range suffix.Query() {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	// add headers
	for k, vs := range c.DefaultHeaders {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", defaultUserAgent)
	}

	return req
}
