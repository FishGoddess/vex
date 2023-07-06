// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

import "testing"

// go test -v -cover -run=^TestNewDefaultConfig$
func TestNewDefaultConfig(t *testing.T) {
	conf := newDefaultConfig()

	if conf.maxConnected != 256 {
		t.Errorf("conf.maxConnected %d is wrong", conf.maxConnected)
	}

	if conf.maxIdle != 256 {
		t.Errorf("conf.maxIdle %d is wrong", conf.maxIdle)
	}

	if conf.blockOnFull != true {
		t.Errorf("conf.LimitStrategy %+v is wrong", conf.blockOnFull)
	}
}

// go test -v -cover -run=^TestConfigApplyOptions$
func TestConfigApplyOptions(t *testing.T) {
	conf := newDefaultConfig()

	conf.ApplyOptions([]Option{
		WithMaxConnected(128),
		WithMaxIdle(32),
		WithNonBlockOnFull(),
	})

	if conf.maxConnected != 128 {
		t.Errorf("conf.maxConnected %d is wrong", conf.maxConnected)
	}

	if conf.maxIdle != 32 {
		t.Errorf("conf.maxIdle %d is wrong", conf.maxIdle)
	}

	if conf.blockOnFull {
		t.Errorf("conf.blockOnFull %+v is wrong", conf.blockOnFull)
	}
}
