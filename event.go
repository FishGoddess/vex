// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

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

// EventHandler is the handler of event.
type EventHandler interface {
	// HandleEvent handles events.
	HandleEvent(e Event)
}
