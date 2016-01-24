package gophers

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"text/template"
)

var (
	apibRequestTemplate = template.Must(template.New("apibRequest").Parse(strings.TrimSpace(`
+ Request ({{ .ContentType }})

{{ .Body }}
`)))

	apibResponseTemplate = template.Must(template.New("apibResponse").Parse(strings.TrimSpace(`
+ Response {{ .StatusCode }}

    + Headers

{{ .Headers }}

    + Body

{{ .Body }}
`)))
)

// Recorder is a common interface of all request/response recorders.
type Recorder interface {
	RecordRequest(req *Request, status, headers, body []byte, wc io.WriteCloser) (err error)
	RecordResponse(resp *Response, status, headers, body []byte, wc io.WriteCloser) (err error)
}

// PlainRecorder writes request and response to plain text file.
type PlainRecorder struct{}

func (r *PlainRecorder) record(status, headers, body []byte, wc io.WriteCloser) (err error) {
	write := func(b []byte) {
		_, err = wc.Write(b)
		if err != nil {
			return
		}
	}

	write(status)
	write(headers)
	write([]byte("\n\n"))
	write(body)
	return wc.Close()
}

func (r *PlainRecorder) RecordRequest(req *Request, status, headers, body []byte, wc io.WriteCloser) (err error) {
	return r.record(status, headers, body, wc)
}

func (r *PlainRecorder) RecordResponse(resp *Response, status, headers, body []byte, wc io.WriteCloser) (err error) {
	return r.record(status, headers, body, wc)
}

// PlainRecorder writes request and response to file in API Blueprint format.
type APIBRecorder struct{}

func (r *APIBRecorder) RecordRequest(req *Request, status, headers, body []byte, wc io.WriteCloser) (err error) {
	indent := strings.Repeat(" ", 13)

	// indent body
	var bodyS []string
	s := bufio.NewScanner(bytes.NewReader(body))
	for s.Scan() {
		bodyS = append(bodyS, indent+s.Text())
	}
	if err = s.Err(); err != nil {
		return
	}

	err = apibRequestTemplate.Execute(wc, map[string]interface{}{
		"ContentType": req.Header.Get("Content-Type"),
		"Body":        strings.Join(bodyS, "\n"),
	})
	if err == nil {
		err = wc.Close()
	}
	return
}

func (r *APIBRecorder) RecordResponse(resp *Response, status, headers, body []byte, wc io.WriteCloser) (err error) {
	indent := strings.Repeat(" ", 12)

	// indent headers
	var headersS []string
	s := bufio.NewScanner(bytes.NewReader(headers))
	for s.Scan() {
		headersS = append(headersS, indent+s.Text())
	}
	if err = s.Err(); err != nil {
		return
	}

	// indent body
	var bodyS []string
	s = bufio.NewScanner(bytes.NewReader(body))
	for s.Scan() {
		bodyS = append(bodyS, indent+s.Text())
	}
	if err = s.Err(); err != nil {
		return
	}

	err = apibResponseTemplate.Execute(wc, map[string]interface{}{
		"StatusCode": resp.StatusCode,
		"Headers":    strings.Join(headersS, "\n"),
		"Body":       strings.Join(bodyS, "\n"),
	})
	if err == nil {
		err = wc.Close()
	}
	return
}

// check interfaces
var (
	_ Recorder = new(PlainRecorder)
	_ Recorder = new(APIBRecorder)
)
