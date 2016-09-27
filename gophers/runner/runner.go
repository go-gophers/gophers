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

var (
	mConcurrency = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gophers",
		Subsystem: "load",
		Name:      "concurrency",
		Help:      "Current concurrency",
	}, []string{"test"})
	mDuration = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "gophers",
		Subsystem: "load",
		Name:      "duration",
		Help:      "Load test duration in seconds",
		MaxAge:    15 * time.Second,
	}, []string{"test", "state"})
)

func init() {
	prometheus.MustRegister(mConcurrency, mDuration)
}

type addedTest struct {
	name   string
	test   TestFunc
	weight int
}

// Runner contains test functions.
type Runner struct {
	l     *log.Logger
	tests []addedTest
}

// New creates new runner with given logger.
func New(l *log.Logger) *Runner {
	if config.Default.HTTPAddr != "" {
		http.Handle("/metrics", prometheus.Handler())
		l.Printf("Prometheus: http://%s/metrics", config.Default.HTTPAddr)
		l.Printf("expvar    : http://%s/debug/vars", config.Default.HTTPAddr)
		l.Printf("pprof     : http://%s/debug/pprof/", config.Default.HTTPAddr)
		go func() { l.Fatal(http.ListenAndServe(config.Default.HTTPAddr, nil)) }()
	}

	return &Runner{l: l}
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

// Test runs tests matching regexp in random order.
func (r *Runner) Test(re *regexp.Regexp) {
	var failedTests, skippedTests []string
	for _, p := range rand.Perm(len(r.tests)) {
		test := r.tests[p]
		if re != nil && !re.MatchString(test.name) {
			continue
		}

		r.l.Printf("=== TEST %s", test.name)

		state := run(test.test, log.New(os.Stderr, "", 0))

		switch state {
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
		os.Exit(1)
	}
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

// Load runs tests matching regexp in random order.
func (r *Runner) Load(re *regexp.Regexp, loader Loader) {
	reporter := newReporter(r.l)
	var failedTests, skippedTests []string
	for _, p := range rand.Perm(len(r.tests)) {
		test := r.tests[p]
		if re != nil && !re.MatchString(test.name) {
			continue
		}

		reporter.setName(test.name)
		allStart := time.Now()
		r.l.Printf("=== LOAD %s (%s)", test.name, loader)

		pool := taskpool.New(taskRun, 0)
		var out *taskOutput
		ticker := time.NewTicker(time.Second)
		inputCh := pool.Input
	For:
		for {
			select {
			case inputCh <- &taskInput{test.name, test.test}:
				// nothing

			case o, ok := <-pool.Output:
				if !ok {
					break For
				}
				out = o.(*taskOutput)
				mDuration.WithLabelValues(test.name, out.state.String()).Observe(out.duration.Seconds())

			case t := <-ticker.C:
				c := loader.Count(t.Sub(allStart))
				mConcurrency.WithLabelValues(test.name).Set(float64(c))
				if c < 0 {
					r.l.Print("stopping...")
					ticker.Stop()
					close(pool.Input)
					inputCh = nil
					go pool.Wait()
					continue
				}
				if pool.Resize(uint(c)) {
					r.l.Printf("concurrency changed to %d", c)
				}
			}
		}

		r.l.Print(out.buf.String())

		switch out.state {
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
		os.Exit(1)
	}
}
