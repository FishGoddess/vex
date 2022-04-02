package vex

import (
	stdlog "log"
	"net"
)

const (
	eventServing      = 1
	eventShutdown     = 2
	eventConnected    = 3
	eventDisconnected = 4
)

var (
	// Log logs some messages.
	Log = stdlog.Printf

	// Dial dials a net connection.
	Dial = net.Dial

	// Notify notifies an event.
	Notify = defaultNotify
)

// makeBytes makes a new byte slice.
func makeBytes(initial int32) []byte {
	return make([]byte, initial)
}

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

// notify publishes events.
func notify(e Event) {
	if Notify != nil {
		Notify(e)
	}
}

func defaultNotify(e Event) {
	if e.Serving() {
		log("vex: server is serving...")
	}

	if e.Shutdown() {
		log("vex: server is shutdown...")
	}
}

// Event is the type of server actions.
type Event int8

// Serving returns if event is server serving.
func (e Event) Serving() bool {
	return e == eventServing
}

// Shutdown returns if event is server shutdown.
func (e Event) Shutdown() bool {
	return e == eventShutdown
}

// Connected returns if event is client connected.
func (e Event) Connected() bool {
	return e == eventConnected
}

// Disconnected returns if event is client disconnected.
func (e Event) Disconnected() bool {
	return e == eventDisconnected
}
