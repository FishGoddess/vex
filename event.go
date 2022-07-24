// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"net"
)

// ServerStartEvent is event representing start of a server.
type ServerStartEvent struct {
	// Server is the server starting.
	Server *Server
}

// ServerShutdownEvent is event representing shutdown of a server.
type ServerShutdownEvent struct {
	// Server is the server shutting down.
	Server *Server
}

// ServerGotConnectedEvent is event representing a new connection to a server.
type ServerGotConnectedEvent struct {
	// Server is the server got connected.
	Server *Server

	// LocalAddr is the local address of the new connection.
	LocalAddr net.Addr

	// RemoteAddr is the remote address of the new connection.
	RemoteAddr net.Addr
}

// ServerGotDisconnectedEvent is event representing a connection disconnected from a server.
type ServerGotDisconnectedEvent struct {
	// Server is the server got disconnected.
	Server *Server

	// LocalAddr is the local address of the disconnected connection.
	LocalAddr net.Addr

	// RemoteAddr is the remote address of the disconnected connection.
	RemoteAddr net.Addr
}

// EventListener is the listener for events.
type EventListener struct {
	// OnServerStart will be called when receives ServerStartEvents.
	OnServerStart func(event ServerStartEvent)

	// OnServerShutdown will be called when receives ServerShutdownEvents.
	OnServerShutdown func(event ServerShutdownEvent)

	// OnServerGotConnected will be called when receives ServerGotConnectedEvents.
	OnServerGotConnected func(event ServerGotConnectedEvent)

	// OnServerGotDisconnected will be called when receives ServerGotDisconnectedEvents.
	OnServerGotDisconnected func(event ServerGotDisconnectedEvent)
}

// NewLogEventListener returns an event listener which will log all events received.
func NewLogEventListener() EventListener {
	return EventListener{
		OnServerStart: func(event ServerStartEvent) {
			log("vex: server %s is starting...", event.Server.Name())
		},
		OnServerShutdown: func(event ServerShutdownEvent) {
			log("vex: server %s is shutdown...", event.Server.Name())
		},
		OnServerGotConnected: func(event ServerGotConnectedEvent) {
			log("vex: %s connected to server %s...", event.RemoteAddr.String(), event.Server.Name())
		},
		OnServerGotDisconnected: func(event ServerGotDisconnectedEvent) {
			log("vex: %s disconnected from server %s...", event.RemoteAddr.String(), event.Server.Name())
		},
	}
}

// CallOnServerStart calls OnServerStart safely.
func (el *EventListener) CallOnServerStart(event ServerStartEvent) {
	if el.OnServerStart != nil {
		el.OnServerStart(event)
	}
}

// CallOnServerShutdown calls OnServerShutdown safely.
func (el *EventListener) CallOnServerShutdown(event ServerShutdownEvent) {
	if el.OnServerShutdown != nil {
		el.OnServerShutdown(event)
	}
}

// CallOnServerGotConnected calls OnServerGotConnected safely.
func (el *EventListener) CallOnServerGotConnected(event ServerGotConnectedEvent) {
	if el.OnServerGotConnected != nil {
		el.OnServerGotConnected(event)
	}
}

// CallOnServerGotDisconnected calls OnServerGotDisconnected safely.
func (el *EventListener) CallOnServerGotDisconnected(event ServerGotDisconnectedEvent) {
	if el.OnServerGotDisconnected != nil {
		el.OnServerGotDisconnected(event)
	}
}
