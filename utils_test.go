package gophers

import (
	"fmt"
)

type FakeTB struct {
	Logs   []string
	Errors []string
	Fatals []string
}

func (f *FakeTB) Logf(format string, a ...interface{}) {
	f.Logs = append(f.Logs, fmt.Sprintf(format, a))
}

func (f *FakeTB) Errorf(format string, a ...interface{}) {
	f.Errors = append(f.Errors, fmt.Sprintf(format, a))
}

func (f *FakeTB) Fatalf(format string, a ...interface{}) {
	f.Fatals = append(f.Fatals, fmt.Sprintf(format, a))
}

// check interface
var _ TestingTB = new(FakeTB)
