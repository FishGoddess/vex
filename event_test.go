// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import "testing"

// go test -v -cover -run=^TestNewLogEventListener$
func TestNewLogEventListener(t *testing.T) {
	listener := NewLogEventListener()
	if listener.OnServerStart == nil {
		t.Error("listener.OnServerStart == nil")
	}

	if listener.OnServerShutdown == nil {
		t.Error("listener.OnServerShutdown == nil")
	}

	if listener.OnServerGotConnected == nil {
		t.Error("listener.OnServerGotConnected == nil")
	}

	if listener.OnServerGotDisconnected == nil {
		t.Error("listener.OnServerGotDisconnected == nil")
	}
}

// go test -v -cover -run=^TestLogEventListener$
func TestLogEventListener(t *testing.T) {
	listener := EventListener{}
	listener.CallOnServerStart(ServerStartEvent{})
	listener.CallOnServerShutdown(ServerShutdownEvent{})
	listener.CallOnServerGotConnected(ServerGotConnectedEvent{})
	listener.CallOnServerGotDisconnected(ServerGotDisconnectedEvent{})
}
