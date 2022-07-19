// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import "testing"

// go test -v -cover -run=^TestNewDefaultPoolConfig$
func TestNewDefaultPoolConfig(t *testing.T) {
	config := NewDefaultConfig()
	if config.MaxConnected != 4096 {
		t.Errorf("config.MaxConnected %d != 4096", config.MaxConnected)
	}

	if config.MaxIdle != 4096 {
		t.Errorf("config.MaxIdle %d != 4096", config.MaxIdle)
	}

	if config.LimitStrategy != limitStrategyBlock {
		t.Errorf("config.LimitStrategy %+v != %d", config.LimitStrategy, limitStrategyBlock)
	}

	if config.ReadBufferSize != 4*1024*1024 {
		t.Errorf("config.ReadBufferSize %d != 4*1024*1024", config.ReadBufferSize)
	}

	if config.WriteBufferSize != 4*1024*1024 {
		t.Errorf("config.WriteBufferSize %d != 4*1024*1024", config.WriteBufferSize)
	}
}

// go test -v -cover -run=^TestConfigApplyOptions$
func TestConfigApplyOptions(t *testing.T) {
	config := NewDefaultConfig()
	config.ApplyOptions([]Option{
		WithMaxConnected(128),
		WithMaxIdle(32),
		WithNewOnLimit(),
		WithReadBufferSize(64),
		WithWriteBufferSize(512),
	})

	if config.MaxConnected != 128 {
		t.Errorf("config.MaxConnected %d != 128", config.MaxConnected)
	}

	if config.MaxIdle != 32 {
		t.Errorf("config.MaxIdle %d != 32", config.MaxIdle)
	}

	if config.LimitStrategy != limitStrategyNew {
		t.Errorf("config.LimitStrategy %+v != %d", config.LimitStrategy, limitStrategyNew)
	}

	if config.ReadBufferSize != 64 {
		t.Errorf("config.ReadBufferSize %d != 64", config.ReadBufferSize)
	}

	if config.WriteBufferSize != 512 {
		t.Errorf("config.WriteBufferSize %d != 512", config.WriteBufferSize)
	}
}
