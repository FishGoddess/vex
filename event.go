// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

const (
	eventServing      Event = 1
	eventShutdown     Event = 2
	eventConnected    Event = 3
	eventDisconnected Event = 4
)

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

// DefaultEventHandler is the default event handler.
type DefaultEventHandler struct {
	name string
}

// NewDefaultEventHandler returns a new default event handler with given name.
func NewDefaultEventHandler(name string) *DefaultEventHandler {
	return &DefaultEventHandler{
		name: name,
	}
}

// HandleEvent handles events.
func (deh *DefaultEventHandler) HandleEvent(e Event) {
	if e.Serving() {
		if deh.name == "" {
			log("vex: server is serving...")
		} else {
			log("vex: server %s is serving...", deh.name)
		}
	}

	if e.Shutdown() {
		if deh.name == "" {
			log("vex: server is shutdown...")
		} else {
			log("vex: server %s is shutdown...", deh.name)
		}
	}
}
