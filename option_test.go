// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"testing"
	"time"
)

// go test -v -cover -run=^TestWithConnTimeout$
func TestWithConnTimeout(t *testing.T) {
	c := &config{ConnTimeout: 0}
	WithConnTimeout(time.Hour)(c)
	if c.ConnTimeout != time.Hour {
		t.Errorf("c.ConnTimeout %d != time.Hour", c.ConnTimeout)
	}
}

// go test -v -cover -run=^TestWithBufferSize$
func TestWithBufferSize(t *testing.T) {
	c := &config{ReadBufferSize: 0, WriteBufferSize: 0}
	WithReadBufferSize(64)(c)
	WithWriteBufferSize(512)(c)

	if c.ReadBufferSize != 64 {
		t.Errorf("c.ReadBufferSize %d != 64", c.ReadBufferSize)
	}

	if c.WriteBufferSize != 512 {
		t.Errorf("c.WriteBufferSize %d != 512", c.WriteBufferSize)
	}
}

// go test -v -cover -run=^TestWithMaxConnected$
func TestWithMaxConnected(t *testing.T) {
	c := &config{MaxConnected: 0}
	WithMaxConnected(64)(c)
	if c.MaxConnected != 64 {
		t.Errorf("c.MaxConnected %d != 64", c.MaxConnected)
	}
}

// go test -v -cover -run=^TestWithEventHandler$
func TestWithEventHandler(t *testing.T) {
	c := &config{EventHandler: nil}
	handler := NewDefaultEventHandler("")
	WithEventHandler(handler)(c)
	if c.EventHandler != handler {
		t.Errorf("c.EventHandler %p != %p", c.EventHandler, handler)
	}
}
