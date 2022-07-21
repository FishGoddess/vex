// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"testing"
	"time"
)

// go test -v -cover -run=^TestNewDefaultConfig$
func TestNewDefaultConfig(t *testing.T) {
	config := newDefaultConfig()
	if config.ConnTimeout != 8*time.Hour {
		t.Errorf("config.ConnTimeout %d != 8*time.Hour", config.ConnTimeout)
	}

	if config.ReadBufferSize != 4*1024*1024 {
		t.Errorf("config.ReadBufferSize %d != 4*1024*1024", config.ReadBufferSize)
	}

	if config.WriteBufferSize != 4*1024*1024 {
		t.Errorf("config.WriteBufferSize %d != 4*1024*1024", config.WriteBufferSize)
	}

	if config.MaxConnected != 4096 {
		t.Errorf("config.MaxConnected %d != 4096", config.MaxConnected)
	}
}

// go test -v -cover -run=^TestConfigApplyOptions$
func TestConfigApplyOptions(t *testing.T) {
	config := newDefaultConfig()
	config.ApplyOptions([]Option{
		WithConnTimeout(time.Hour),
		WithReadBufferSize(64),
		WithWriteBufferSize(512),
		WithMaxConnected(128),
	})

	if config.ConnTimeout != time.Hour {
		t.Errorf("config.ConnTimeout %d != time.Hour", config.ConnTimeout)
	}

	if config.ReadBufferSize != 64 {
		t.Errorf("config.ReadBufferSize %d != 64", config.ReadBufferSize)
	}

	if config.WriteBufferSize != 512 {
		t.Errorf("config.WriteBufferSize %d != 512", config.WriteBufferSize)
	}

	if config.MaxConnected != 128 {
		t.Errorf("config.MaxConnected %d != 128", config.MaxConnected)
	}
}
