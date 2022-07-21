// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

import (
	"testing"
)

// go test -v -cover -run=^TestWithMaxConnected$
func TestWithMaxConnected(t *testing.T) {
	c := &config{MaxConnected: 0}
	WithMaxConnected(64)(c)
	if c.MaxConnected != 64 {
		t.Errorf("c.MaxConnected %d != 64", c.MaxConnected)
	}
}

// go test -v -cover -run=^TestWithMaxIdle$
func TestWithMaxIdle(t *testing.T) {
	c := &config{MaxIdle: 0}
	WithMaxIdle(64)(c)
	if c.MaxIdle != 64 {
		t.Errorf("c.MaxIdle %d != 64", c.MaxIdle)
	}
}

// go test -v -cover -run=^TestWithNonBlockOnLimit$
func TestWithNonBlockOnLimit(t *testing.T) {
	c := &config{BlockOnFull: true}
	WithNonBlockOnFull()(c)
	if c.BlockOnFull {
		t.Errorf("c.BlockOnFull %+v != false", c.BlockOnFull)
	}
}
