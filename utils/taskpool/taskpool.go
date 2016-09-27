// Package taskpool provides task pool with variable amount of worker goroutines.
package taskpool

import (
	"sync"
)

// Task is a work function with one input parameter and one output.
type Task func(input interface{}) (output interface{})

// TaskPool is task pool with variable amount of worker goroutines.
// It should be created with New().
// There should be only one sender to Input and one receiver from Output.
// Input must be closed by user before invoking Wait().
// Output is closed by Wait().
type TaskPool struct {
	Input  chan interface{}
	Output chan interface{}

	task Task
	stop chan struct{}
	size uint
	wg   sync.WaitGroup
}

// New creates new TaskPool with given task and initial size.
func New(task Task, size uint) *TaskPool {
	p := &TaskPool{
		Input:  make(chan interface{}),
		Output: make(chan interface{}),
		task:   task,
		stop:   make(chan struct{}),
	}
	p.Resize(size)
	return p
}

// Wait waits for worker goroutines to finish and closes Output.
// User should call this method after closing Input.
// User should continue to receive from Output until it's closed.
//
// This method is not thread-safe.
func (p *TaskPool) Wait() {
	p.wg.Wait()
	close(p.Output)
}

// Size returns current amount of worker goroutines.
//
// This method is not thread-safe.
func (p *TaskPool) Size() uint {
	return p.size
}

// Resize changes amount of worker goroutines and returns true if amount was changed.
//
// This method is not thread-safe.
func (p *TaskPool) Resize(size uint) bool {
	if p.size == size {
		return false
	}

	for p.size < size {
		p.size++
		p.wg.Add(1)
		go p.worker()
	}

	for p.size > size {
		p.size--
		p.stop <- struct{}{}
	}

	return true
}

func (p *TaskPool) worker() {
	defer p.wg.Done()

	for {
		select {
		case i, ok := <-p.Input:
			if !ok {
				return
			}
			p.Output <- p.task(i)
		case <-p.stop:
			return
		}
	}
}
