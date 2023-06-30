// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

import (
	"testing"
)

// go test -v -cover -run=^TestWithMaxConnected$
func TestWithMaxConnected(t *testing.T) {
	conf := &Config{MaxConnected: 0}
	WithMaxConnected(64)(conf)

	if conf.MaxConnected != 64 {
		t.Errorf("conf.MaxConnected %d is wrong", conf.MaxConnected)
	}
}

// go test -v -cover -run=^TestWithMaxIdle$
func TestWithMaxIdle(t *testing.T) {
	conf := &Config{MaxIdle: 0}
	WithMaxIdle(16)(conf)

	if conf.MaxIdle != 16 {
		t.Errorf("conf.MaxIdle %d is wrong", conf.MaxIdle)
	}
}

// go test -v -cover -run=^TestWithConnections$
func TestWithConnections(t *testing.T) {
	conf := &Config{MaxConnected: 0, MaxIdle: 0}
	WithConnections(64)(conf)

	if conf.MaxConnected != 64 {
		t.Errorf("conf.MaxConnected %d is wrong", conf.MaxConnected)
	}

	if conf.MaxIdle != 64 {
		t.Errorf("conf.MaxIdle %d is wrong", conf.MaxIdle)
	}
}

// go test -v -cover -run=^TestWithNonBlockOnFull$
func TestWithNonBlockOnFull(t *testing.T) {
	conf := &Config{BlockOnFull: true}
	WithNonBlockOnFull()(conf)

	if conf.BlockOnFull {
		t.Errorf("conf.BlockOnFull %+v is wrong", conf.BlockOnFull)
	}
}
