package runner

import (
	"fmt"
	"runtime"
	"sync/atomic"

	"github.com/go-gophers/gophers"
	"github.com/go-gophers/gophers/utils/log"
)

type state int64

const (
	passed state = iota
	failed
	skipped
	panicked
)

func (s state) String() string {
	switch s {
	case passed:
		return "passed"
	case failed:
		return "failed"
	case skipped:
		return "skipped"
	case panicked:
		return "panicked"
	default:
		panic(fmt.Sprintf("unexpected state %d", s))
	}
}

// check interface
var _ fmt.Stringer = panicked

// testingT implements gophers.TestingT.
type testingT struct {
	l *log.Logger
	s int64
}

func (t *testingT) state() state {
	return state(atomic.LoadInt64(&t.s))
}

func (t *testingT) passed() bool {
	return t.state() == passed
}

func (t *testingT) Log(args ...interface{}) {
	t.l.Println(append([]interface{}{"        "}, args...)...)
}

func (t *testingT) Logf(format string, args ...interface{}) {
	t.l.Printf("         "+format, args...)
}

func (t *testingT) Failed() bool {
	return t.state() == failed
}

func (t *testingT) Fail() {
	atomic.StoreInt64(&t.s, int64(failed))
}

func (t *testingT) Error(args ...interface{}) {
	t.Log(args...)
	t.Fail()
}

func (t *testingT) Errorf(format string, args ...interface{}) {
	t.Logf(format, args...)
	t.Fail()
}

func (t *testingT) FailNow() {
	t.Fail()
	runtime.Goexit()
}

func (t *testingT) Fatal(args ...interface{}) {
	t.Log(args...)
	t.FailNow()
}

func (t *testingT) Fatalf(format string, args ...interface{}) {
	t.Logf(format, args...)
	t.FailNow()
}

func (t *testingT) Skipped() bool {
	return t.state() == skipped
}

func (t *testingT) SkipNow() {
	atomic.StoreInt64(&t.s, int64(skipped))
	runtime.Goexit()
}

func (t *testingT) Skip(args ...interface{}) {
	t.Log(args...)
	t.SkipNow()
}

func (t *testingT) Skipf(format string, args ...interface{}) {
	t.Logf(format, args...)
	t.SkipNow()
}

func (t *testingT) panic() {
	atomic.StoreInt64(&t.s, int64(panicked))
}

// check interface
var _ gophers.TestingT = new(testingT)
