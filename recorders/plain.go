package recorders

import (
	"io"
	"net/http"
)

// Plain records request and response to writers in plain text format.
type Plain struct{}

func (r *Plain) record(status, headers, body []byte, wc io.WriteCloser) (err error) {
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

// RecordRequest writes request's status, headers and body.
func (r *Plain) RecordRequest(req *http.Request, status, headers, body []byte, wc io.WriteCloser) (err error) {
	return r.record(status, headers, body, wc)
}

// RecordResponse writes response's status, headers and body.
func (r *Plain) RecordResponse(resp *http.Response, status, headers, body []byte, wc io.WriteCloser) (err error) {
	return r.record(status, headers, body, wc)
}

// check interface
var _ Interface = new(Plain)
