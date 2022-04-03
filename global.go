package vex

import (
	stdlog "log"
	"net"
)

var (
	// Log logs some messages.
	Log = stdlog.Printf

	// Dial dials a net connection.
	Dial = net.Dial
)

// log records logs with format and v.
func log(format string, v ...interface{}) {
	if Log != nil {
		Log(format, v...)
	}
}

// dial records logs with format and v.
func dial(network string, address string) (net.Conn, error) {
	if Dial == nil {
		panic("vex: Dial == nil")
	}

	return Dial(network, address)
}

// makeBytes makes a new byte slice.
func makeBytes(initial int32) []byte {
	return make([]byte, initial)
}
