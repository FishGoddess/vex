// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

import "testing"

// go test -v -cover -run=^TestNewDefaultConfig$
func TestNewDefaultConfig(t *testing.T) {
	config := newDefaultConfig()
	if config.MaxConnected != 256 {
		t.Errorf("config.MaxConnected %d != 256", config.MaxConnected)
	}

	if config.MaxIdle != 256 {
		t.Errorf("config.MaxIdle %d != 256", config.MaxIdle)
	}

	if config.BlockOnFull != true {
		t.Errorf("config.LimitStrategy %+v != true", config.BlockOnFull)
	}
}

// go test -v -cover -run=^TestConfigApplyOptions$
func TestConfigApplyOptions(t *testing.T) {
	config := newDefaultConfig()
	config.ApplyOptions([]Option{
		WithMaxConnected(128),
		WithMaxIdle(32),
		WithNonBlockOnFull(),
	})

	if config.MaxConnected != 128 {
		t.Errorf("config.MaxConnected %d != 128", config.MaxConnected)
	}

	if config.MaxIdle != 32 {
		t.Errorf("config.MaxIdle %d != 32", config.MaxIdle)
	}

	if config.BlockOnFull {
		t.Errorf("config.BlockOnFull %+v != true", config.BlockOnFull)
	}
}
