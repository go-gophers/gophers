package recorders

import (
	"io"
	"net/http"
)

// Interface is a common interface of all request/response recorders.
type Interface interface {
	RecordRequest(req *http.Request, status, headers, body []byte, wc io.WriteCloser) (err error)
	RecordResponse(resp *http.Response, status, headers, body []byte, wc io.WriteCloser) (err error)
}
