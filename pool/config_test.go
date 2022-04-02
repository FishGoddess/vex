// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

import "testing"

// go test -v -cover -run=^TestNewDefaultPoolConfig$
func TestNewDefaultPoolConfig(t *testing.T) {
	config := newDefaultConfig()
	if config.maxConnected != 64 {
		t.Errorf("config.maxConnected %d != 64", config.maxConnected)
	}

	if config.maxIdle != 64 {
		t.Errorf("config.maxIdle %d != 64", config.maxIdle)
	}

	if config.fullStrategy != FullStrategyBlock {
		t.Errorf("config.fullStrategy %+v != %d", config.fullStrategy, FullStrategyBlock)
	}
}

// go test -v -cover -run=^TestConfigApplyOptions$
func TestConfigApplyOptions(t *testing.T) {
	config := newDefaultConfig()
	config.applyOptions([]Option{
		WithMaxOpened(128),
		WithMaxIdle(32),
		WithFullStrategy(FullStrategyNew),
	})

	if config.maxConnected != 128 {
		t.Errorf("config.maxConnected %d != 64", config.maxConnected)
	}

	if config.maxIdle != 32 {
		t.Errorf("config.maxIdle %d != 64", config.maxIdle)
	}

	if config.fullStrategy != FullStrategyNew {
		t.Errorf("config.fullStrategy %+v != %d", config.fullStrategy, FullStrategyNew)
	}
}
