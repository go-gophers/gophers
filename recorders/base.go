// Package recorders providers request/response recorders for Gophers tool.
package recorders

import (
	"io"
	"net/http"
)

// Interface is a common interface of all request/response recorders.
type Interface interface {
	// RecordRequest writes request's status, headers and body.
	RecordRequest(req *http.Request, status, headers, body []byte, wc io.WriteCloser) (err error)

	// RecordResponse writes response's status, headers and body.
	RecordResponse(resp *http.Response, status, headers, body []byte, wc io.WriteCloser) (err error)
}
