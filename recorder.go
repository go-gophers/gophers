package gophers

import (
	"io"
	"text/template"
)

var (
	apibRequestTemplate = template.Must(template.New("apibRequest").Parse(`
+ Request ({{ .ContentType }})

        {{ .Body }}
`))

	apibResponseTemplate = template.Must(template.New("apibResponse").Parse(`
+ Response {{ .StatusCode }} ({{ .ContentType }})

    + Headers

            {{ .Headers }}

    + Body

            {{ .Body }}
`))
)

type Recorder interface {
	RecordRequest(req *Request, status, headers, body []byte, wc io.WriteCloser) (err error)
	RecordResponse(resp *Response, status, headers, body []byte, wc io.WriteCloser) (err error)
}

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

// check interfaces
var (
	_ Recorder = new(PlainRecorder)
	// _ Recorder = new(APIBRecorder)
)
