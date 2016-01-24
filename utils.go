package gophers

import (
	"io"
)

// Implemented by *bytes.Buffer, *bytes.Reader, *strings.Reader.
type lenReader interface {
	io.Reader
	Len() int
}

type TestingTB interface {
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}
