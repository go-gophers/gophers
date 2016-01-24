package gophers

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
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

func (c *Client) NewRequest(t TestingTB, method string, urlStr string, body io.Reader) *Request {
	r, err := http.NewRequest(method, urlStr, nil)
	if err != nil {
		t.Fatalf("can't create request: %s", err)
	}

	req := &Request{Request: r}
	req.SetBodyReader(body)

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

	return req
}

func (c *Client) Do(t TestingTB, req *Request, expectedStatusCode int) *Response {
	status, headers, body, err := DumpRequest(req.Request)
	if err != nil {
		t.Fatalf("can't dump request: %s", err)
	}

	colorF := func(b []byte) string { return color.BlueString(string(b)) }
	if *vF {
		t.Logf("\n%s\n%s\n\n%s\n", colorF(status), colorF(headers), colorF(body))
	} else {
		t.Logf("\n%s\n", colorF(status))
	}

	if req.record(req.RequestRecorder, status, headers, body) {
		if f, ok := req.RequestRecorder.(*os.File); ok {
			t.Logf("request recorded to %s", f.Name())
		} else {
			t.Logf("request recorded")
		}
	}

	resp, err := c.HTTPClient.Do(req.Request)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		t.Fatalf("can't make request: %s", err)
	}

	status, headers, body, err = DumpResponse(resp)
	if err != nil {
		t.Fatalf("can't dump response: %s", err)
	}

	switch {
	case resp.StatusCode >= 400:
		colorF = func(b []byte) string { return color.RedString(string(b)) }
	case resp.StatusCode >= 300:
		colorF = func(b []byte) string { return color.YellowString(string(b)) }
	default:
		colorF = func(b []byte) string { return color.GreenString(string(b)) }
	}

	if *vF {
		t.Logf("\n%s\n%s\n\n%s\n", colorF(status), colorF(headers), colorF(body))
	} else {
		t.Logf("\n%s\n", colorF(status))
	}

	if req.record(req.ResponseRecorder, status, headers, body) {
		if f, ok := req.ResponseRecorder.(*os.File); ok {
			t.Logf("response recorded to %s", f.Name())
		} else {
			t.Logf("response recorded")
		}
	}

	if resp.StatusCode != expectedStatusCode {
		t.Errorf("%s %s: expected %d, got %s", req.Method, req.URL.String(), expectedStatusCode, resp.Status)
	}
	return &Response{Response: resp}
}

func (c *Client) Head(t TestingTB, urlStr string, expectedStatusCode int) *Response {
	return c.Do(t, c.NewRequest(t, "GET", urlStr, nil), expectedStatusCode)
}

func (c *Client) Get(t TestingTB, urlStr string, expectedStatusCode int) *Response {
	return c.Do(t, c.NewRequest(t, "GET", urlStr, nil), expectedStatusCode)
}

func (c *Client) Post(t TestingTB, urlStr string, body io.Reader, expectedStatusCode int) *Response {
	return c.Do(t, c.NewRequest(t, "POST", urlStr, body), expectedStatusCode)
}

func (c *Client) Put(t TestingTB, urlStr string, body io.Reader, expectedStatusCode int) *Response {
	return c.Do(t, c.NewRequest(t, "PUT", urlStr, body), expectedStatusCode)
}

func (c *Client) Patch(t TestingTB, urlStr string, body io.Reader, expectedStatusCode int) *Response {
	return c.Do(t, c.NewRequest(t, "PATCH", urlStr, body), expectedStatusCode)
}

func (c *Client) Delete(t TestingTB, urlStr string, expectedStatusCode int) *Response {
	return c.Do(t, c.NewRequest(t, "DELETE", urlStr, nil), expectedStatusCode)
}
