// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import "testing"

// go test -v -cover -run=^TestEvent$
func TestEvent(t *testing.T) {
	if !eventServing.Serving() {
		t.Error("eventServing.Serving() returns false")
	}

	if !eventShutdown.Shutdown() {
		t.Error("eventShutdown.Shutdown() returns false")
	}

	if !eventConnected.Connected() {
		t.Error("eventConnected.Connected() returns false")
	}

	if !eventDisconnected.Disconnected() {
		t.Error("eventDisconnected.Disconnected() returns false")
	}
}

// go test -v -cover -run=^TestNewDefaultEventHandler$
func TestNewDefaultEventHandler(t *testing.T) {
	name := "xxx"
	handler := NewDefaultEventHandler(name)
	if handler.name != name {
		t.Errorf("handler.name %s != name %s", handler.name, name)
	}
}
