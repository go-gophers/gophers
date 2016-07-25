package gophers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	"github.com/go-gophers/gophers/jsons"
)

func isChunked(te []string) bool {
	for _, v := range te {
		if v == "chunked" {
			return true
		}
	}
	return false
}

func dump(b []byte, te []string) (status, headers, body []byte, err error) {
	p := bytes.SplitN(b, []byte("\r\n\r\n"), 2)
	headers, body = p[0], p[1]
	p = bytes.SplitN(headers, []byte("\r\n"), 2)
	status, headers = p[0], p[1]

	if len(body) > 0 && isChunked(te) {
		r := httputil.NewChunkedReader(bytes.NewReader(body))
		body, err = ioutil.ReadAll(r)
		if err != nil {
			return
		}
	}

	return
}

// dumpRequest returns representation of req with status line, headers and body.
// It uses httputil.DumpRequestOut and additionally converts body from chunked encoding
// to identity so it can be used in documentation.
func dumpRequest(req *http.Request) (status, headers, body []byte, err error) {
	var b []byte
	b, err = httputil.DumpRequestOut(req, true)
	if err != nil {
		return
	}
	return dump(b, req.TransferEncoding)
}

// dumpRequest returns representation of res with status line, headers and body.
// It uses httputil.DumpResponse and additionally converts body from chunked encoding
// to identity so it can be used in documentation.
func dumpResponse(res *http.Response) (status, headers, body []byte, err error) {
	var b []byte
	b, err = httputil.DumpResponse(res, true)
	if err != nil {
		return
	}
	return dump(b, res.TransferEncoding)
}

// bodyRepr returns representation of body depending on content type.
// It may be indented, shortened or returned as is.
// It returns nil for empty body.
func bodyRepr(contentType string, body []byte) []byte {
	if len(body) == 0 {
		return nil
	}

	switch {
	case strings.Contains(contentType, "json"):
		return []byte(jsons.ParseBytes(body).Indent())

	default:
		for _, r := range string(body) {
			if !strconv.IsPrint(r) {
				return []byte(fmt.Sprintf("[%d bytes data]", len(body)))
			}
		}
		return body
	}
}
