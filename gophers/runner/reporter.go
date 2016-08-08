package runner

import (
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"

	"github.com/go-gophers/gophers/utils/log"
)

var infoSignals = []os.Signal{syscall.SIGUSR1}

type reporter struct {
	info   chan os.Signal
	logger *log.Logger
	name   *atomic.Value
}

func newReporter(logger *log.Logger) *reporter {
	r := &reporter{
		info:   make(chan os.Signal, 1),
		logger: logger,
		name:   new(atomic.Value),
	}
	r.setName("")

	go func() {
		for range r.info {
			r.report()
		}
	}()
	signal.Notify(r.info, infoSignals...)

	return r
}

func (r *reporter) setName(name string) {
	r.name.Store(name)
}

func (r *reporter) stop() {
	signal.Stop(r.info)
	close(r.info)
}

func (r *reporter) report() {
	var collector prometheus.Collector = mDuration // TODO make collector argument?

	metrics := make(chan prometheus.Metric)
	done := make(chan struct{})
	go func() {
	Next:
		for m := range metrics {
			var metric dto.Metric
			err := m.Write(&metric)
			if err != nil {
				r.logger.Fatal(err)
			}

			for _, l := range metric.Label {
				if l.GetName() != "test" {
					continue
				}
				name := r.name.Load().(string)
				if name != "" && l.GetValue() != name {
					continue Next
				}
				r.logger.Printf("%s stats:", l.GetValue())
			}
			r.logger.Printf("\t%d observations, sum %s",
				metric.Summary.GetSampleCount(), time.Duration(metric.Summary.GetSampleSum()*float64(time.Second)))
			for _, q := range metric.Summary.Quantile {
				r.logger.Printf("\t%.2fφ: %s",
					q.GetQuantile(), time.Duration(q.GetValue()*float64(time.Second)))
			}
		}
		close(done)
	}()

	collector.Collect(metrics)
	close(metrics)
	<-done
}
