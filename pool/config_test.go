// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

import "testing"

// go test -v -cover -run=^TestNewDefaultConfig$
func TestNewDefaultConfig(t *testing.T) {
	conf := newDefaultConfig()

	if conf.MaxConnected != 256 {
		t.Errorf("conf.MaxConnected %d is wrong", conf.MaxConnected)
	}

	if conf.MaxIdle != 256 {
		t.Errorf("conf.MaxIdle %d is wrong", conf.MaxIdle)
	}

	if conf.BlockOnFull != true {
		t.Errorf("conf.LimitStrategy %+v is wrong", conf.BlockOnFull)
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

	if conf.MaxConnected != 128 {
		t.Errorf("conf.MaxConnected %d is wrong", conf.MaxConnected)
	}

	if conf.MaxIdle != 32 {
		t.Errorf("conf.MaxIdle %d is wrong", conf.MaxIdle)
	}

	if conf.BlockOnFull {
		t.Errorf("conf.BlockOnFull %+v is wrong", conf.BlockOnFull)
	}
}
