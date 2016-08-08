package gophers

import (
	"errors"
	"io"
)

// TestingT is a copy of testing.TB interface without private methods.
// It can be used to hook other testing libraries, frameworks, and runners, such as gophers tool.
type TestingT interface {
	// Log formats its arguments using default formatting, analogous to Println, and records the text in the error log.
	Log(args ...interface{})

	// Logf formats its arguments according to the format, analogous to Printf, and records the text in the error log.
	Logf(format string, args ...interface{})

	// Failed reports whether the function has failed.
	Failed() bool

	// Fail marks the function as having failed but continues execution.
	Fail()

	// Error is equivalent to Log followed by Fail.
	Error(args ...interface{})

	// Errorf is equivalent to Logf followed by Fail.
	Errorf(format string, args ...interface{})

	// FailNow marks the function as having failed and stops its execution. Execution will continue at the next test.
	// FailNow must be called from the goroutine running the test function, not from other goroutines created during the test.
	// Calling FailNow does not stop those other goroutines.
	FailNow()

	// Fatal is equivalent to Log followed by FailNow.
	Fatal(args ...interface{})

	// Fatalf is equivalent to Logf followed by FailNow.
	Fatalf(format string, args ...interface{})

	// Skipped reports whether the test was skipped.
	Skipped() bool

	// SkipNow marks the test as having been skipped and stops its execution. Execution will continue at the next test. See also FailNow.
	// SkipNow must be called from the goroutine running the test, not from other goroutines created during the test.
	// Calling SkipNow does not stop those other goroutines.
	SkipNow()

	// Skip is equivalent to Log followed by SkipNow.
	Skip(args ...interface{})

	// Skipf is equivalent to Logf followed by SkipNow.
	Skipf(format string, args ...interface{})
}

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
