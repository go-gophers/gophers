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
func (c *Client) NewRequest(t TestingTB, method string, urlStr string, body fmt.Stringer) *Request {
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
// It also logs and records them and checks response status code.
// Request and response Body fields are filled, inner *http.(Request|Response).Body fields
// are replaced by stubs.
// In case of error it fails test.
func (c *Client) Do(t TestingTB, req *Request, expectedStatusCode int) *Response {
	status, headers, body, err := dumpRequest(req.Request)
	if err != nil {
		t.Fatalf("can't dump request: %s", err)
	}

	repr := bodyRepr(req.Header.Get("Content-Type"), body)

	// color methods *String accepts format as first argument
	// any format-like strings passed to it will be treated as format string
	// ex. URL encoded strings /content?from=2016-03-31T08%3A00%3A00%2B03%3A00
	// we should explicitly pass `%s` format as first argument
	colorF := func(b []byte) string { return color.BlueString("%s", string(b)) }
	if *vF {
		t.Logf("\n%s\n%s\n\n%s\n", colorF(status), colorF(headers), colorF(repr))
	} else {
		t.Logf("\n%s\n", colorF(status))
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

	if *vF {
		t.Logf("\n%s\n%s\n\n%s\n", colorF(status), colorF(headers), colorF(repr))
	} else {
		t.Logf("\n%s\n", colorF(status))
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

	if resp.StatusCode != expectedStatusCode {
		t.Errorf("%s %s: expected %d, got %s", req.Method, req.URL.String(), expectedStatusCode, resp.Status)
	}
	return resp
}

// Head makes HEAD request. See Do for more details.
func (c *Client) Head(t TestingTB, urlStr string, expectedStatusCode int) *Response {
	return c.Do(t, c.NewRequest(t, "HEAD", urlStr, nil), expectedStatusCode)
}

// Get makes GET request. See Do for more details.
func (c *Client) Get(t TestingTB, urlStr string, expectedStatusCode int) *Response {
	return c.Do(t, c.NewRequest(t, "GET", urlStr, nil), expectedStatusCode)
}

// Post makes POST request. See Do for more details.
func (c *Client) Post(t TestingTB, urlStr string, body fmt.Stringer, expectedStatusCode int) *Response {
	return c.Do(t, c.NewRequest(t, "POST", urlStr, body), expectedStatusCode)
}

// Put makes PUT request. See Do for more details.
func (c *Client) Put(t TestingTB, urlStr string, body fmt.Stringer, expectedStatusCode int) *Response {
	return c.Do(t, c.NewRequest(t, "PUT", urlStr, body), expectedStatusCode)
}

// Patch makes PATCH request. See Do for more details.
func (c *Client) Patch(t TestingTB, urlStr string, body fmt.Stringer, expectedStatusCode int) *Response {
	return c.Do(t, c.NewRequest(t, "PATCH", urlStr, body), expectedStatusCode)
}

// Delete makes DELETE request. See Do for more details.
func (c *Client) Delete(t TestingTB, urlStr string, expectedStatusCode int) *Response {
	return c.Do(t, c.NewRequest(t, "DELETE", urlStr, nil), expectedStatusCode)
}
