package net

import (
	"net"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	mReadCalls = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gophers",
		Subsystem: "net",
		Name:      "read_calls",
		Help:      "Read calls",
	}, []string{"error"})
	mReadBytes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gophers",
		Subsystem: "net",
		Name:      "read_bytes",
		Help:      "Read bytes",
	}, []string{"error"})

	mWriteCalls = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gophers",
		Subsystem: "net",
		Name:      "write_calls",
		Help:      "Write calls",
	}, []string{"error"})
	mWriteBytes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gophers",
		Subsystem: "net",
		Name:      "write_bytes",
		Help:      "Write bytes",
	}, []string{"error"})

	mCloseCalls = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gophers",
		Subsystem: "net",
		Name:      "close_calls",
		Help:      "Close calls",
	}, []string{"error"})
)

func init() {
	prometheus.MustRegister(mReadCalls, mReadBytes)
	prometheus.MustRegister(mWriteCalls, mWriteBytes)
	prometheus.MustRegister(mCloseCalls)
}

type Conn struct {
	net.Conn
}

// check interface
var _ net.Conn = new(Conn)

// Read reads data from the connection.
// Read can be made to time out and return a Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.
func (c *Conn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	mReadCalls.WithLabelValues(errorLabelValue(err)).Inc()
	mReadBytes.WithLabelValues(errorLabelValue(err)).Add(float64(n))
	return
}

// Write writes data to the connection.
// Write can be made to time out and return a Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline.
func (c *Conn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	mWriteCalls.WithLabelValues(errorLabelValue(err)).Inc()
	mWriteBytes.WithLabelValues(errorLabelValue(err)).Add(float64(n))
	return
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c *Conn) Close() (err error) {
	err = c.Conn.Close()
	mCloseCalls.WithLabelValues(errorLabelValue(err)).Inc()
	return
}
