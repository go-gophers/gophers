package net

import (
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/go-gophers/gophers/utils/log"
)

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

var (
	dns   = make(map[string][]string)
	dnsRW sync.RWMutex
)

func lookupHost(host string) (string, error) {
	ip := net.ParseIP(host)
	if ip != nil {
		return ip.String(), nil
	}

	dnsRW.RLock()
	addrs := dns[host]
	dnsRW.RUnlock()
	if addrs != nil {
		return addrs[rand.Intn(len(addrs))], nil
	}

	dnsRW.Lock()
	defer dnsRW.Unlock()
	ips, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}

	// keep only IPv4 addresses
	// FIXME make it configurable
	addrs = make([]string, 0, len(ips))
	for _, ip := range ips {
		ip4 := ip.To4()
		if ip4 != nil {
			addrs = append(addrs, ip4.String())
		}
	}

	dns[host] = addrs
	return addrs[rand.Intn(len(addrs))], nil
}

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
	conn, err := net.Dial(network, addr)
	if err == nil {
		mDialResults.WithLabelValues(network, "0").Inc()
		log.Debugf("gophers/net.Dial(%q, %q): connection established (in %s)", network, addr, time.Now().Sub(start))
	} else {
		mDialResults.WithLabelValues(network, "1").Inc()
		log.Printf("gophers/net.Dial(%q, %q): %s (in %s)", network, addr, err, time.Now().Sub(start))
	}
	return &Conn{conn}, err
}
