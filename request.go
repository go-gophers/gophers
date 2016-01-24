package gophers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Request struct {
	*http.Request

	RequestRecorder  io.WriteCloser
	ResponseRecorder io.WriteCloser
	RecordStatusLine bool
	RecordHeaders    bool
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

func (req *Request) EnableRecording(baseFileName string) *Request {
	ext := filepath.Ext(baseFileName)
	base := strings.TrimSuffix(baseFileName, ext)

	reqF, err := os.Create(base + "_request" + ext)
	if err != nil {
		panic(err)
	}
	req.RequestRecorder = reqF

	resF, err := os.Create(base + "_response" + ext)
	if err != nil {
		panic(err)
	}
	req.ResponseRecorder = resF

	return req
}

func (req *Request) record(wc io.WriteCloser, status, headers, body []byte) bool {
	if wc == nil {
		return false
	}

	write := func(b []byte) {
		_, err := wc.Write(b)
		if err != nil {
			panic(err)
		}
	}

	if req.RecordStatusLine {
		write(status)
	}
	if req.RecordHeaders {
		write(headers)
		write([]byte("\n"))
	}
	write(body)

	err := wc.Close()
	if err != nil {
		panic(err)
	}

	return true
}
