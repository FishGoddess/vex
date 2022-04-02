// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import "testing"

// go test -v -cover -run=^TestWithMaxConnected$
func TestWithMaxConnected(t *testing.T) {
	c := &Config{MaxConnected: 0}
	WithMaxConnected(64)(c)
	if c.MaxConnected != 64 {
		t.Errorf("c.MaxConnected %d != 64", c.MaxConnected)
	}
}

// go test -v -cover -run=^TestWithMaxIdle$
func TestWithMaxIdle(t *testing.T) {
	c := &Config{MaxIdle: 0}
	WithMaxIdle(64)(c)
	if c.MaxIdle != 64 {
		t.Errorf("c.MaxIdle %d != 64", c.MaxIdle)
	}
}

// go test -v -cover -run=^TestWithLimitStrategy$
func TestWithLimitStrategy(t *testing.T) {
	c := &Config{LimitStrategy: 0}
	WithBlockOnLimit()(c)
	if c.LimitStrategy != limitStrategyBlock {
		t.Errorf("c.LimitStrategy %d != %d", c.LimitStrategy, limitStrategyBlock)
	}

	WithFailedOnLimit()(c)
	if c.LimitStrategy != limitStrategyFailed {
		t.Errorf("c.LimitStrategy %d != %d", c.LimitStrategy, limitStrategyFailed)
	}

	WithNewOnLimit()(c)
	if c.LimitStrategy != limitStrategyNew {
		t.Errorf("c.LimitStrategy %d != %d", c.LimitStrategy, limitStrategyNew)
	}
}

// go test -v -cover -run=^TestWithBufferSize$
func TestWithBufferSize(t *testing.T) {
	c := &Config{ReadBufferSize: 0, WriteBufferSize: 0}
	WithReadBufferSize(64)(c)
	WithWriteBufferSize(512)(c)

	if c.ReadBufferSize != 64 {
		t.Errorf("c.ReadBufferSize %d != 64", c.ReadBufferSize)
	}

	if c.WriteBufferSize != 512 {
		t.Errorf("c.WriteBufferSize %d != 512", c.WriteBufferSize)
	}
}
