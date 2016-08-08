package runner

import (
	"bytes"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-gophers/gophers"
	"github.com/go-gophers/gophers/utils/log"
)

func testLogger() (*log.Logger, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	return log.New(buf, "", 0), buf
}

func TestRunOk(t *testing.T) {
	l, buf := testLogger()
	state := run(func(tt gophers.TestingT) {}, l)
	assert.Equal(t, passed, state)
	assert.Equal(t, "", buf.String())
}

func TestRunFatal(t *testing.T) {
	l, buf := testLogger()
	state := run(func(tt gophers.TestingT) { tt.Fatal("fatal error") }, l)
	assert.Equal(t, failed, state)
	assert.Equal(t, "         fatal error\n", buf.String())
}

func TestRunSkip(t *testing.T) {
	l, buf := testLogger()
	state := run(func(tt gophers.TestingT) { tt.Skip("skip") }, l)
	assert.Equal(t, skipped, state)
	assert.Equal(t, "         skip\n", buf.String())
}

func TestRunPanic(t *testing.T) {
	l, buf := testLogger()
	state := run(func(tt gophers.TestingT) { panic("PANIC!") }, l)
	assert.Equal(t, panicked, state)
	assert.True(t, strings.HasPrefix(buf.String(),
		`         panic: PANIC!
         github.com/go-gophers/gophers/gophers/runner.TestRunPanic.func1`),
		"%s", buf.String())
}

func TestRunPanicNil(t *testing.T) {
	l, buf := testLogger()
	state := run(func(tt gophers.TestingT) { panic(nil) }, l)
	assert.Equal(t, panicked, state)
	assert.True(t, strings.HasPrefix(buf.String(),
		`         test executed panic(nil) or runtime.Goexit()
         github.com/go-gophers/gophers/gophers/runner.TestRunPanicNil.func1`),
		"%s", buf.String())
}

func TestRunBug(t *testing.T) {
	l, buf := testLogger()
	state := run(func(tt gophers.TestingT) { var p *int; _ = *p }, l)
	assert.Equal(t, panicked, state)
	assert.True(t, strings.HasPrefix(buf.String(),
		`         panic: runtime error: invalid memory address or nil pointer dereference
         runtime.panicmem`),
		"%s", buf.String())
}

func TestRunExit(t *testing.T) {
	l, buf := testLogger()
	state := run(func(tt gophers.TestingT) { runtime.Goexit() }, l)
	assert.Equal(t, panicked, state)
	assert.True(t, strings.HasPrefix(buf.String(),
		`         test executed panic(nil) or runtime.Goexit()
         github.com/go-gophers/gophers/gophers/runner.TestRunExit.func1`),
		"%s", buf.String())
}
