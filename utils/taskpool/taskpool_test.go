package taskpool

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type TaskPoolSuite struct {
	suite.Suite
	numGoroutine int
	p            *TaskPool
	logf         func(format string, args ...interface{})
}

func TestTaskPoolSuite(t *testing.T) {
	suite.Run(t, new(TaskPoolSuite))
}

func (suite *TaskPoolSuite) SetupTest() {
	suite.numGoroutine = runtime.NumGoroutine()
	suite.logf = suite.T().Logf
}

func (suite *TaskPoolSuite) TearDownTest() {
	close(suite.p.Input)
	suite.p.Wait()
	time.Sleep(time.Millisecond) // let readOutput() exit
	suite.Require().Equal(suite.numGoroutine, runtime.NumGoroutine(), "leaked goroutine")
}

func (suite *TaskPoolSuite) task(in interface{}) interface{} {
	suite.logf("task(%v)", in)
	return nil
}

func (suite *TaskPoolSuite) readOutput() {
	for out := range suite.p.Output {
		suite.logf("output: %v", out)
	}
}

func (suite *TaskPoolSuite) TestSimple() {
	suite.p = New(suite.task, 2)
	go suite.readOutput()
	for i := 0; i < 5; i++ {
		suite.p.Input <- i
	}
}

func (suite *TaskPoolSuite) TestSequence() {
	suite.p = New(suite.task, 1)
	go suite.readOutput()
	for i := 0; i < 5; i++ {
		suite.p.Input <- i
	}
}

func (suite *TaskPoolSuite) TestEmpty() {
	suite.p = New(suite.task, 0)
}
