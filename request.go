package gophers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gophergala2016/gophers/recorders"
)

// Request represents HTTP request and recording parameters.
type Request struct {
	*http.Request

	Recorder   recorders.Interface
	RequestWC  io.WriteCloser
	ResponseWC io.WriteCloser
}

// SetBodyReader sets request body with given reader.
// It also try to set Content-Length header.
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
	lr, ok := r.(lenReader)
	if ok {
		req.ContentLength = int64(lr.Len())
	}

	return req
}

// SetBodyString sets request body with given string.
// It also sets Content-Length header.
func (req *Request) SetBodyString(s string) *Request {
	return req.SetBodyReader(strings.NewReader(s))
}

// SetBodyStringer sets request body with given Stringer.
// It also sets Content-Length header.
func (req *Request) SetBodyStringer(s fmt.Stringer) *Request {
	if s == nil {
		return req
	}
	return req.SetBodyString(s.String())
}

// AddHeaders adds headers to request.
func (req *Request) AddHeaders(h http.Header) *Request {
	for k, vs := range h {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	return req
}

// AddCookies adds cookies to request.
func (req *Request) AddCookies(c []http.Cookie) *Request {
	for _, e := range c {
		req.AddCookie(&e)
	}
	return req
}

// EnableRecording enables recording of this request and following response
// to files with given base name. Recorder type is selected by extension:
// recorders.APIB for ".apib", recorders.Plain for any other.
func (req *Request) EnableRecording(baseFileName string) *Request {
	ext := filepath.Ext(baseFileName)
	base := strings.TrimSuffix(baseFileName, ext)

	reqF, err := os.Create(base + "_request" + ext)
	if err != nil {
		panic(err)
	}
	req.RequestWC = reqF

	resF, err := os.Create(base + "_response" + ext)
	if err != nil {
		panic(err)
	}
	req.ResponseWC = resF

	switch ext {
	case ".apib":
		req.Recorder = new(recorders.APIB)
	default:
		req.Recorder = new(recorders.Plain)
	}

	return req
}
