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
	DefaultCookies []http.Cookie
}

// TODO use io.MultiReader and http.DetectContentType to sent ContentType?

func NewClient(base url.URL) *Client {
	return &Client{
		Base:       base,
		HTTPClient: http.DefaultClient,
		DefaultHeaders: http.Header{
			"User-Agent": []string{defaultUserAgent},
		},
	}
}

func (c *Client) NewRequest(t testing.TB, method string, urlStr string) *Request {
	req, err := http.NewRequest(method, urlStr, nil)
	if err != nil {
		t.Fatal(err)
	}

	newUrl := c.Base

	// update request URL path, check for '//'
	if strings.HasSuffix(newUrl.Path, "/") && strings.HasPrefix(req.URL.Path, "/") {
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/")
	}
	newUrl.Path += req.URL.Path

	// update request URL query
	q := newUrl.Query()
	for k, vs := range req.URL.Query() {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	newUrl.RawQuery = q.Encode()

	req.URL = &newUrl

	// add headers
	for k, vs := range c.DefaultHeaders {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	// add cookies
	for _, c := range c.DefaultCookies {
		req.AddCookie(&c)
	}

	return &Request{Request: req}
}

func (c *Client) Do(t testing.TB, req *Request, expectedStatusCode int) *Response {
	headers, body, err := DumpRequest(req.Request)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%s\n\n%s\n", headers, body)

	resp, err := c.HTTPClient.Do(req.Request)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		t.Fatal(err)
	}

	headers, body, err = DumpResponse(resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%s\n\n%s\n", headers, body)

	if resp.StatusCode != expectedStatusCode {
		t.Errorf("%s %s: expected %d, got %s", req.Method, req.URL.String(), expectedStatusCode, resp.Status)
	}
	return &Response{Response: resp}
}
