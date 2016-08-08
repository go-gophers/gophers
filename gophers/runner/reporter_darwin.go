package runner

import (
	"syscall"
)

func init() {
	infoSignals = append(infoSignals, syscall.SIGINFO) // ^T
}
