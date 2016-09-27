package net

import (
	"io"
)

// errorLabelValue returns value for Prometheus metric label "error":
// "ok" for no error, "EOF" for io.EOF, "error" otherwise.
func errorLabelValue(err error) string {
	if err == nil {
		return "ok"
	}
	if err == io.EOF {
		return "EOF"
	}
	return "error"
}
