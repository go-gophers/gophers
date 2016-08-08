package runner

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

func infinity(d uint) string {
	if d == 0 {
		return "∞"
	}
	return strconv.FormatUint(uint64(d), 10)
}

type Loader interface {
	Count(duration time.Duration) int
	String() string
}

type StepLoader struct {
	minC         uint
	maxC         uint
	stepC        uint
	stepDuration time.Duration
	lastStep     uint
}

func NewStepLoader(minC uint, maxC uint, stepC uint, stepDuration time.Duration) (*StepLoader, error) {
	if minC == 0 {
		return nil, errors.New("minC must be positive")
	}
	if stepC == 0 {
		return nil, errors.New("stepC must be positive")
	}
	if stepDuration == 0 {
		return nil, errors.New("stepDuration must be positive")
	}

	lastStep := (maxC - minC) / stepC
	if (maxC-minC)%stepC > 0 {
		lastStep++
	}
	if maxC == 0 {
		lastStep = 1<<64 - 1
	}
	return &StepLoader{
		minC:         minC,
		maxC:         maxC,
		stepC:        stepC,
		stepDuration: stepDuration,
		lastStep:     lastStep,
	}, nil
}

func (s *StepLoader) Count(duration time.Duration) int {
	stepN := uint(duration / s.stepDuration)
	switch {
	case stepN > s.lastStep:
		return -1
	case stepN == s.lastStep:
		return int(s.maxC)
	default:
		return int(s.minC + s.stepC*stepN)
	}
}

func (s *StepLoader) String() string {
	m := fmt.Sprintf("%d, %s, %d, %s", s.minC, infinity(s.maxC), s.stepC, s.stepDuration)
	if s.maxC == 0 {
		return fmt.Sprintf("StepLoader(%s)", m)
	} else {
		return fmt.Sprintf("StepLoader(%s: %d steps in %s)", m, s.lastStep+1, s.stepDuration*time.Duration(s.lastStep+1))
	}
}

// check interfaces
var (
	_ Loader = new(StepLoader)
)
