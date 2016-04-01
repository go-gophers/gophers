package gophers

import (
	"errors"
	"io"
)

// lenReader is implemented by *bytes.Buffer, *bytes.Reader, *strings.Reader.
type lenReader interface {
	io.Reader
	Len() int
}

// errorReadCloser is a io.ReadCloser, which returns error for any operation.
type errorReadCloser struct {
}

var rcError = errors.New(
	"gophers: do not use Request.Request.Body or Response.Response.Body (io.ReadCloser), " +
		"use Request.Body or Response.Body ([]byte)",
)

func (errorReadCloser) Read(p []byte) (int, error) {
	return 0, rcError
}

func (errorReadCloser) Close() error {
	return rcError
}

// check interface
var _ io.ReadCloser = errorReadCloser{}

// TestingTB is a subset of testing.TB interface.
// It can be used to hook other testing librariers and frameworks.
type TestingTB interface {
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}
