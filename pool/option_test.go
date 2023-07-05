// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

import (
	"testing"
)

// go test -v -cover -run=^TestWithMaxConnected$
func TestWithMaxConnected(t *testing.T) {
	conf := &Config{maxConnected: 0}
	WithMaxConnected(64)(conf)

	if conf.maxConnected != 64 {
		t.Errorf("conf.maxConnected %d is wrong", conf.maxConnected)
	}
}

// go test -v -cover -run=^TestWithMaxIdle$
func TestWithMaxIdle(t *testing.T) {
	conf := &Config{maxIdle: 0}
	WithMaxIdle(16)(conf)

	if conf.maxIdle != 16 {
		t.Errorf("conf.maxIdle %d is wrong", conf.maxIdle)
	}
}

// go test -v -cover -run=^TestWithConnections$
func TestWithConnections(t *testing.T) {
	conf := &Config{maxConnected: 0, maxIdle: 0}
	WithConnections(64)(conf)

	if conf.maxConnected != 64 {
		t.Errorf("conf.maxConnected %d is wrong", conf.maxConnected)
	}

	if conf.maxIdle != 64 {
		t.Errorf("conf.maxIdle %d is wrong", conf.maxIdle)
	}
}

// go test -v -cover -run=^TestWithNonBlockOnFull$
func TestWithNonBlockOnFull(t *testing.T) {
	conf := &Config{blockOnFull: true}
	WithNonBlockOnFull()(conf)

	if conf.blockOnFull {
		t.Errorf("conf.blockOnFull %+v is wrong", conf.blockOnFull)
	}
}
