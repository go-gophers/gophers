package gophers

import (
	"io"
)

// Implemented by *bytes.Buffer, *bytes.Reader, *strings.Reader.
type LenReader interface {
	io.Reader
	Len() int
}
