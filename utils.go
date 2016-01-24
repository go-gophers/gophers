package gophers

import (
	"io"
)

// Implemented by *bytes.Buffer, *bytes.Reader, *strings.Reader.
type lenReader interface {
	io.Reader
	Len() int
}

// TestingTB is a subset of testing.TB interface.
// It can be used to hook other testing librariers and frameworks.
type TestingTB interface {
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}
