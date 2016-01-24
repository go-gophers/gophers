package gophers

import (
	"io"
)

type Recorder interface {
	Setup(req *Request, wc io.WriteCloser)
	Record(status, headers, body []byte) (err error)
}

type PlainRecorder struct {
	req *Request
	wc  io.WriteCloser
}

func (r *PlainRecorder) Setup(req *Request, wc io.WriteCloser) {
	r.req = req
	r.wc = wc
}

func (r *PlainRecorder) Record(status, headers, body []byte) (err error) {
	write := func(b []byte) {
		_, err = r.wc.Write(b)
		if err != nil {
			return
		}
	}

	if r.req.RecordStatusLine {
		write(status)
	}
	if r.req.RecordHeaders {
		write(headers)
		write([]byte("\n"))
	}
	write(body)

	err = r.wc.Close()
	return
}

type APIBRecorder struct {
	req *Request
}

// check interfaces
var (
	_ Recorder = new(PlainRecorder)
	// _ Recorder = new(APIBRecorder)
)
