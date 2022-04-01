package vex

import (
	"log"
	"net"
	"os"
)

var (
	// Log logs some messages.
	Log = log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags).Printf

	// Dial dials a net connection.
	Dial = net.Dial

	// MakeBytes makes a new byte slice.
	MakeBytes = makeBytes
)

// makeBytes makes a new byte slice.
func makeBytes(initial int32) []byte {
	return make([]byte, initial)
}
