package gophers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// Implemented by *bytes.Buffer, *bytes.Reader, *strings.Reader.
type LenReader interface {
	io.Reader
	Len() int
}

type Request struct {
	*http.Request
}

func (req *Request) SetBodyReader(r io.Reader) *Request {
	if r == nil {
		return req
	}

	rc, ok := r.(io.ReadCloser)
	if !ok && r != nil {
		rc = ioutil.NopCloser(r)
	}
	req.Body = rc

	req.ContentLength = 0
	lr, ok := r.(LenReader)
	if ok {
		req.ContentLength = int64(lr.Len())
	}

	return req
}

func (req *Request) SetBodyString(s string) *Request {
	return req.SetBodyReader(strings.NewReader(s))
}

func (req *Request) SetBodyStringer(s fmt.Stringer) *Request {
	if s == nil {
		return req
	}
	return req.SetBodyString(s.String())
}

func (req *Request) AddHeaders(h http.Header) *Request {
	for k, vs := range h {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	return req
}

func (req *Request) AddCookies(c []http.Cookie) *Request {
	for _, e := range c {
		req.AddCookie(&e)
	}
	return req
}
