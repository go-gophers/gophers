// Package runner implements test runner used by gophers tool.
package runner

import (
	"bytes"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"time"

	_ "expvar"         // for side-effects
	_ "net/http/pprof" // for side-effects

	"github.com/prometheus/client_golang/prometheus"

	"github.com/go-gophers/gophers"
	"github.com/go-gophers/gophers/config"
	"github.com/go-gophers/gophers/utils/log"
	"github.com/go-gophers/gophers/utils/taskpool"
)

// TestFunc is a test function.
type TestFunc func(gophers.TestingT)

// FailMode defines how early load test should fail if test fails.
type FailMode int

const (
	// FailEarly terminates failed load test as fast as possible.
	FailEarly FailMode = iota

	// FailStep terminates failed load test before next load step.
	FailStep

	// FailContinue doesn't terminate failed load test.
	FailContinue
)

// shared Prometheus metrics for all Runners
var (
	mTestRun = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gophers",
		Subsystem: "test",
		Name:      "run",
		Help:      "Test run count",
	}, []string{"suite", "test", "state"})
	mLoadConcurrency = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gophers",
		Subsystem: "load",
		Name:      "concurrency",
		Help:      "Load test current concurrency",
	}, []string{"suite", "test"})
	mLoadDuration = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "gophers",
		Subsystem: "load",
		Name:      "duration",
		Help:      "Load test duration in seconds",
		MaxAge:    15 * time.Second,
	}, []string{"suite", "test", "state"})
)

func init() {
	prometheus.MustRegister(mTestRun, mLoadConcurrency, mLoadDuration)
}

type addedTest struct {
	name   string
	test   TestFunc
	weight int
}

// Runner contains test functions.
type Runner struct {
	suite string
	l     *log.Logger
	tests []addedTest
}

// New creates new runner with given suite name and logger.
func New(suite string, l *log.Logger) *Runner {
	// start Prometheus and debug endpoints
	if config.Default.HTTPAddr != "" {
		http.Handle("/metrics", prometheus.Handler())
		l.Printf("Prometheus: http://%s/metrics", config.Default.HTTPAddr)
		l.Printf("expvar    : http://%s/debug/vars", config.Default.HTTPAddr)
		l.Printf("pprof     : http://%s/debug/pprof/", config.Default.HTTPAddr)
		go func() {
			l.Fatal(http.ListenAndServe(config.Default.HTTPAddr, nil))
		}()
	}

	return &Runner{
		suite: suite,
		l:     l,
	}
}

// Add registers test function under given name.
func (r *Runner) Add(name string, test TestFunc, weight int) {
	r.tests = append(r.tests, addedTest{name, test, weight})
}

func errorStack(t *testingT) {
	pc := make([]uintptr, 100)
	n := runtime.Callers(5, pc)
	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i])
		if f == nil {
			t.Error("-")
		} else {
			file, line := f.FileLine(pc[i] - 1)
			t.Errorf("%s (%s:%d)", f.Name(), file, line)
		}
	}
}

// run runs single test with given logger.
func run(test TestFunc, l *log.Logger) state {
	t := &testingT{l: l}
	result := make(chan state)

	go func() {
		var finished bool

		defer func() {
			if p := recover(); p != nil {
				t.Errorf("panic: %v", p)
				errorStack(t)
				t.panic()
			}
			if t.passed() && !finished {
				t.Error("test executed panic(nil) or runtime.Goexit()")
				errorStack(t)
				t.panic()
			}

			result <- t.state()
		}()

		test(t)
		finished = true
	}()

	return <-result
}

func shuffle(n int, seed int64) []int {
	if seed == 0 {
		res := make([]int, n)
		for i := 0; i < n; i++ {
			res[i] = i
		}
		return res
	}

	return rand.New(rand.NewSource(seed)).Perm(n)
}

// Test runs tests matching regexp in order defined by seed.
// If seed is 0, tests are run sequentially.
// Otherwise, seed defines random order of tests.
func (r *Runner) Test(re *regexp.Regexp, seed int64) int {
	var failedTests, skippedTests []string
	for _, p := range shuffle(len(r.tests), seed) {
		test := r.tests[p]
		if re != nil && !re.MatchString(test.name) {
			continue
		}

		r.l.Printf("=== TEST %s", test.name)

		start := time.Now()
		state := run(test.test, log.New(os.Stderr, "", 0))
		duration := time.Since(start)
		mTestRun.WithLabelValues(r.suite, test.name, state.String()).Inc()

		switch state {
		case failed, panicked:
			r.l.Printf("--- FAIL %s (%.2fs)", test.name, duration.Seconds())
			failedTests = append(failedTests, test.name)
		case skipped:
			r.l.Printf("--- SKIP %s (%.2fs)", test.name, duration.Seconds())
			skippedTests = append(skippedTests, test.name)
		case passed:
			r.l.Printf("--- PASS %s (%.2fs)", test.name, duration.Seconds())
		}
	}

	r.l.Printf("%d tests run, %d passed, %d skipped, %d failed.",
		len(r.tests), len(r.tests)-len(skippedTests)-len(failedTests), len(skippedTests), len(failedTests))
	if len(skippedTests) > 0 {
		r.l.Print("Skipped tests:")
		for _, name := range skippedTests {
			r.l.Printf("\t%s", name)
		}
	}
	if len(failedTests) > 0 {
		r.l.Print("Failed tests:")
		for _, name := range failedTests {
			r.l.Printf("\t%s", name)
		}
		return 1
	}
	return 0
}

type taskInput struct {
	name string
	test TestFunc
}

type taskOutput struct {
	name     string
	state    state
	duration time.Duration
	buf      *bytes.Buffer
}

func taskRun(input interface{}) interface{} {
	in := input.(*taskInput)
	buf := new(bytes.Buffer) // TODO use sync.Pool?
	l := log.New(buf, "", 0)
	start := time.Now()
	state := run(in.test, l)
	return &taskOutput{in.name, state, time.Now().Sub(start), buf}
}

func (r *Runner) load(test *addedTest, loader Loader, failMode FailMode) (worstOut *taskOutput) {
	start := time.Now()
	ticker := time.NewTicker(time.Second)
	pool := taskpool.New(taskRun, 0)
	inputCh := pool.Input
	stop := func() {
		if inputCh != nil {
			r.l.Print("stopping...")
			ticker.Stop()
			close(pool.Input)
			inputCh = nil
			go pool.Wait()
		}
	}

	for {
		select {
		case inputCh <- &taskInput{test.name, test.test}:
			// nothing

		case o, ok := <-pool.Output:
			if !ok {
				return
			}

			out := o.(*taskOutput)
			mLoadDuration.WithLabelValues(r.suite, test.name, out.state.String()).Observe(out.duration.Seconds())

			if worstOut == nil || worstOut.state == passed {
				worstOut = out
			}

			if out.state != passed && failMode == FailEarly {
				stop()
				continue
			}

		case t := <-ticker.C:
			c := loader.Count(t.Sub(start))
			switch {
			case uint(c) == pool.Size():
				// nothing

			case c < 0:
				fallthrough

			case worstOut != nil && worstOut.state != passed && failMode == FailStep:
				mLoadConcurrency.WithLabelValues(r.suite, test.name).Set(0)
				stop()
				continue

			default:
				mLoadConcurrency.WithLabelValues(r.suite, test.name).Set(float64(c))
				pool.Resize(uint(c))
				r.l.Printf("concurrency changed to %d", c)
			}
		}
	}
}

// Test runs tests matching regexp in order defined by seed in load test mode.
// If seed is 0, tests are run sequentially.
// Otherwise, seed defines random order of tests.
func (r *Runner) Load(re *regexp.Regexp, seed int64, loader Loader, failMode FailMode) int {
	reporter := newReporter(r.l)
	var failedTests, skippedTests []string
	for _, p := range shuffle(len(r.tests), seed) {
		test := r.tests[p]
		if re != nil && !re.MatchString(test.name) {
			continue
		}

		reporter.setName(test.name)
		r.l.Printf("=== LOAD %s (%s)", test.name, loader)

		worstOut := r.load(&test, loader, failMode)

		r.l.Print(worstOut.buf.String())

		switch worstOut.state {
		case failed, panicked:
			r.l.Printf("--- FAIL %s", test.name)
			failedTests = append(failedTests, test.name)
		case skipped:
			r.l.Printf("--- SKIP %s", test.name)
			skippedTests = append(skippedTests, test.name)
		case passed:
			r.l.Printf("--- PASS %s", test.name)
		}
	}

	reporter.setName("")
	reporter.report()
	reporter.stop()

	r.l.Printf("%d tests run, %d passed, %d skipped, %d failed.",
		len(r.tests), len(r.tests)-len(skippedTests)-len(failedTests), len(skippedTests), len(failedTests))
	if len(skippedTests) > 0 {
		r.l.Print("Skipped tests:")
		for _, name := range skippedTests {
			r.l.Printf("\t%s", name)
		}
	}
	if len(failedTests) > 0 {
		r.l.Print("Failed tests:")
		for _, name := range failedTests {
			r.l.Printf("\t%s", name)
		}
		return 1
	}
	return 0
}
