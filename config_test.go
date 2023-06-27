// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"testing"
	"time"
)

// go test -v -cover -run=^TestNewClientConfig$
func TestNewClientConfig(t *testing.T) {
	conf := newClientConfig("127.0.0.1:5837")

	if conf.address != "127.0.0.1:5837" {
		t.Errorf("conf.address %s is wrong", conf.address)
	}

	if conf.ReadTimeout != 10*time.Minute {
		t.Errorf("conf.ReadTimeout %d is wrong", conf.ReadTimeout)
	}

	if conf.WriteTimeout != 10*time.Minute {
		t.Errorf("conf.WriteTimeout %d is wrong", conf.WriteTimeout)
	}

	if conf.ReadBufferSize != 16*1024 {
		t.Errorf("conf.ReadBufferSize %d is wrong", conf.ReadBufferSize)
	}

	if conf.WriteBufferSize != 16*1024 {
		t.Errorf("conf.WriteBufferSize %d is wrong", conf.WriteBufferSize)
	}
}

// go test -v -cover -run=^TestNewServerConfig$
func TestNewServerConfig(t *testing.T) {
	conf := newServerConfig("127.0.0.1:5837")

	if conf.address != "127.0.0.1:5837" {
		t.Errorf("conf.address %s is wrong", conf.address)
	}

	if conf.Name != "127.0.0.1:5837" {
		t.Errorf("conf.Name %s is wrong", conf.Name)
	}

	if conf.ReadTimeout != 10*time.Minute {
		t.Errorf("conf.ReadTimeout %d is wrong", conf.ReadTimeout)
	}

	if conf.WriteTimeout != 10*time.Minute {
		t.Errorf("conf.WriteTimeout %d is wrong", conf.WriteTimeout)
	}

	if conf.CloseTimeout != time.Minute {
		t.Errorf("conf.CloseTimeout %d is wrong", conf.CloseTimeout)
	}

	if conf.ReadBufferSize != 16*1024 {
		t.Errorf("conf.ReadBufferSize %d is wrong", conf.ReadBufferSize)
	}

	if conf.WriteBufferSize != 16*1024 {
		t.Errorf("conf.WriteBufferSize %d is wrong", conf.WriteBufferSize)
	}
}

// go test -v -cover -run=^TestConfigApplyClientOptions$
func TestConfigApplyClientOptions(t *testing.T) {
	conf := newClientConfig("127.0.0.1:5837")

	conf.ApplyClientOptions([]ClientOption{
		WithClientReadTimeout(time.Second),
		WithClientWriteTimeout(3 * time.Second),
		WithClientReadBufferSize(64),
		WithClientWriteBufferSize(512),
	})

	if conf.address != "127.0.0.1:5837" {
		t.Errorf("conf.address %s is wrong", conf.address)
	}

	if conf.ReadTimeout != time.Second {
		t.Errorf("conf.ReadTimeout %d is wrong", conf.ReadTimeout)
	}

	if conf.WriteTimeout != 3*time.Second {
		t.Errorf("conf.WriteTimeout %d is wrong", conf.WriteTimeout)
	}

	if conf.ReadBufferSize != 64 {
		t.Errorf("config.ReadBufferSize %d is wrong", conf.ReadBufferSize)
	}

	if conf.WriteBufferSize != 512 {
		t.Errorf("config.WriteBufferSize %d is wrong", conf.WriteBufferSize)
	}
}

// go test -v -cover -run=^TestConfigApplyServerOptions$
func TestConfigApplyServerOptions(t *testing.T) {
	conf := newServerConfig("127.0.0.1:5837")

	conf.ApplyServerOptions([]ServerOption{
		WithServerName("test-name"),
		WithServerReadTimeout(time.Second),
		WithServerWriteTimeout(3 * time.Second),
		WithServerCloseTimeout(time.Second),
		WithServerReadBufferSize(64),
		WithServerWriteBufferSize(512),
	})

	if conf.address != "127.0.0.1:5837" {
		t.Errorf("conf.address %s is wrong", conf.address)
	}

	if conf.Name != "test-name" {
		t.Errorf("config.Name %s is wrong", conf.Name)
	}

	if conf.ReadTimeout != time.Second {
		t.Errorf("conf.ReadTimeout %d is wrong", conf.ReadTimeout)
	}

	if conf.WriteTimeout != 3*time.Second {
		t.Errorf("conf.WriteTimeout %d is wrong", conf.WriteTimeout)
	}

	if conf.CloseTimeout != time.Second {
		t.Errorf("config.CloseTimeout %d is wrong", conf.CloseTimeout)
	}

	if conf.ReadBufferSize != 64 {
		t.Errorf("config.ReadBufferSize %d is wrong", conf.ReadBufferSize)
	}

	if conf.WriteBufferSize != 512 {
		t.Errorf("config.WriteBufferSize %d is wrong", conf.WriteBufferSize)
	}
}
