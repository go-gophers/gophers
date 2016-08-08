package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Logger is our custom logger with Debugf method.
type Logger struct {
	*log.Logger
	Debug bool
}

// New creates a new Logger.
func New(out io.Writer, prefix string, flag int) *Logger {
	return &Logger{log.New(out, prefix, flag), false}
}

// Debugf prints message only when Logger Debug flag is set to true.
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.Debug {
		l.Output(2, fmt.Sprintf(format, args...))
	}
}

// Default shared instance.
var Default = New(os.Stderr, "", 0)

// Debugf is an alias for Default.Debugf.
func Debugf(format string, v ...interface{}) {
	Default.Debugf(format, v...)
}

// Fatal is an alias for Default.Fatal.
func Fatal(v ...interface{}) {
	Default.Fatal(v...)
}

// Fatalf is an alias for Default.Fatalf.
func Fatalf(format string, v ...interface{}) {
	Default.Fatalf(format, v...)
}

// Print is an alias for Default.Print.
func Print(v ...interface{}) {
	Default.Print(v...)
}

// Printf is an alias for Default.Printf.
func Printf(format string, v ...interface{}) {
	Default.Printf(format, v...)
}
