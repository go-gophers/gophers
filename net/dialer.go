package net

import (
	"net"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/go-gophers/gophers/config"
	"github.com/go-gophers/gophers/utils/log"
)

// shared Prometheus metrics for all Dial calls
var (
	mDials = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gophers",
		Subsystem: "net",
		Name:      "dials",
		Help:      "Dial calls",
	}, []string{"network"})
	mDialResults = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gophers",
		Subsystem: "net",
		Name:      "dial_results",
		Help:      "Dial results",
	}, []string{"network", "error"})
)

func init() {
	prometheus.MustRegister(mDials, mDialResults)
}

// Dial wraps net.Dial with DNS lookup cache and Prometheus metrics.
func Dial(network, addr string) (net.Conn, error) {
	if network != "tcp" {
		panic("only tcp is supported for now")
	}

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		log.Printf("net.SplitHostPort(%q): %s", addr, err)
		return nil, err
	}

	host, err = lookupHost(host)
	if err != nil {
		log.Printf("lookupHost(%q): %s", host, err)
		return nil, err
	}
	addr = net.JoinHostPort(host, port)

	mDials.WithLabelValues(network).Inc()
	log.Debugf("gophers/net.Dial(%q, %q)", network, addr)
	start := time.Now()
	conn, err := net.DialTimeout(network, addr, config.Default.DialTimeout)
	mDialResults.WithLabelValues(network, errorLabelValue(err)).Inc()
	if err == nil {
		log.Debugf("gophers/net.Dial(%q, %q): connection established (in %s)", network, addr, time.Now().Sub(start))
	} else {
		log.Printf("gophers/net.Dial(%q, %q): %s (in %s)", network, addr, err, time.Now().Sub(start))
	}
	return &Conn{conn}, err
}
