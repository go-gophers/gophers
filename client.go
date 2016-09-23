package gophers

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
)

const (
	defaultUserAgent = "github.com/go-gophers/gophers"
)

// Client wraps base API URL with default headers and cookies.
// Base URL can contain scheme, user info, host, path prefix and default query parameters.
type Client struct {
	Base           url.URL
	HTTPClient     *http.Client
	DefaultHeaders http.Header
	DefaultCookies []http.Cookie
}

// TODO use io.MultiReader and http.DetectContentType to sent ContentType?

// NewClient creates new client with given base URL.
func NewClient(base url.URL) *Client {
	return &Client{
		Base:       base,
		HTTPClient: http.DefaultClient,
		DefaultHeaders: http.Header{
			"User-Agent": []string{defaultUserAgent},
		},
	}
}

// NewRequest creates new request with given method, URL and body.
// It adds URL's path and query parameters to client's base URL.
// It also adds default headers and cookies from client.
// In case of error if fails test.
func (c *Client) NewRequest(t TestingT, method string, urlStr string, body fmt.Stringer) *Request {
	initColors()

	r, err := http.NewRequest(method, urlStr, nil)
	if err != nil {
		t.Fatalf("can't create request: %s", err)
	}

	req := &Request{Request: r}
	req.SetBodyStringer(body)

	newURL := c.Base

	// update request URL path, check for '//'
	if strings.HasSuffix(newURL.Path, "/") && strings.HasPrefix(req.URL.Path, "/") {
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/")
	}
	newURL.Path += req.URL.Path

	// update request URL query
	q := newURL.Query()
	for k, vs := range req.URL.Query() {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	newURL.RawQuery = q.Encode()

	req.URL = &newURL

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

	return req
}

// Do makes request and returns response.
// It also logs and records them and checks that response status code is equal to one of the provided.
// Request and response Body fields are filled, inner *http.(Request|Response).Body fields
// are replaced by stubs.
// In case of error it fails test.
func (c *Client) Do(t TestingT, req *Request, expectedStatuses ...int) *Response {
	status, headers, body, err := dumpRequest(req.Request)
	if err != nil {
		t.Fatalf("can't dump request: %s", err)
	}

	repr := bodyRepr(req.Header.Get("Content-Type"), body)

	colorF := func(b []byte) string { return color.BlueString("%s", string(b)) }
	if DefaultConfig.Verbose {
		t.Logf("\n%s\n%s\n\n%s\n", colorF(status), colorF(headers), colorF(repr))
	} else {
		t.Logf("%s\n", colorF(status))
	}

	if req.Recorder != nil && req.RequestWC != nil {
		err = req.Recorder.RecordRequest(req.Request, status, headers, repr, req.RequestWC)
		if err != nil {
			t.Fatalf("can't record request: %s", err)
		}
		if f, ok := req.RequestWC.(*os.File); ok {
			t.Logf("request recorded to %s", f.Name())
		} else {
			t.Logf("request recorded")
		}
	}

	r, err := c.HTTPClient.Do(req.Request)
	if r != nil {
		origBody := r.Body
		defer func() {
			err = origBody.Close()
			if err != nil {
				t.Fatalf("can't close response body: %s", err)
			}
		}()
	}
	if err != nil {
		t.Fatalf("can't make request: %s", err)
	}

	// put dumped request body back
	req.Body = body
	req.Request.Body = errorReadCloser{}

	resp := &Response{Response: r}

	status, headers, body, err = dumpResponse(resp.Response)
	if err != nil {
		t.Fatalf("can't dump response: %s", err)
	}

	// put dumped response body back
	resp.Body = body
	resp.Response.Body = errorReadCloser{}

	repr = bodyRepr(resp.Header.Get("Content-Type"), body)

	switch {
	case resp.StatusCode >= 400:
		colorF = func(b []byte) string { return color.RedString("%s", string(b)) }
	case resp.StatusCode >= 300:
		colorF = func(b []byte) string { return color.YellowString("%s", string(b)) }
	default:
		colorF = func(b []byte) string { return color.GreenString("%s", string(b)) }
	}

	if DefaultConfig.Verbose {
		t.Logf("\n%s\n%s\n\n%s\n", colorF(status), colorF(headers), colorF(repr))
	} else {
		t.Logf("%s\n", colorF(status))
	}

	if req.Recorder != nil && req.ResponseWC != nil {
		err = req.Recorder.RecordResponse(resp.Response, status, headers, repr, req.ResponseWC)
		if err != nil {
			t.Fatalf("can't record response: %s", err)
		}
		if f, ok := req.ResponseWC.(*os.File); ok {
			t.Logf("response recorded to %s", f.Name())
		} else {
			t.Logf("response recorded")
		}
	}

	if len(expectedStatuses) > 0 {
		var found bool
		for _, s := range expectedStatuses {
			if resp.StatusCode == s {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("%s %s: expected status code to be in %v, got %s", req.Method, req.URL.String(), expectedStatuses, resp.Status)
		}
	}
	return resp
}

// Head makes HEAD request. See Do for more details.
func (c *Client) Head(t TestingT, urlStr string, expectedStatuses ...int) *Response {
	return c.Do(t, c.NewRequest(t, "HEAD", urlStr, nil), expectedStatuses...)
}

// Get makes GET request. See Do for more details.
func (c *Client) Get(t TestingT, urlStr string, expectedStatuses ...int) *Response {
	return c.Do(t, c.NewRequest(t, "GET", urlStr, nil), expectedStatuses...)
}

// Post makes POST request. See Do for more details.
func (c *Client) Post(t TestingT, urlStr string, body fmt.Stringer, expectedStatuses ...int) *Response {
	return c.Do(t, c.NewRequest(t, "POST", urlStr, body), expectedStatuses...)
}

// Put makes PUT request. See Do for more details.
func (c *Client) Put(t TestingT, urlStr string, body fmt.Stringer, expectedStatuses ...int) *Response {
	return c.Do(t, c.NewRequest(t, "PUT", urlStr, body), expectedStatuses...)
}

// Patch makes PATCH request. See Do for more details.
func (c *Client) Patch(t TestingT, urlStr string, body fmt.Stringer, expectedStatuses ...int) *Response {
	return c.Do(t, c.NewRequest(t, "PATCH", urlStr, body), expectedStatuses...)
}

// Delete makes DELETE request. See Do for more details.
func (c *Client) Delete(t TestingT, urlStr string, expectedStatuses ...int) *Response {
	return c.Do(t, c.NewRequest(t, "DELETE", urlStr, nil), expectedStatuses...)
}
