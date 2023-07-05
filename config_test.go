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

	if conf.readTimeout != 10*time.Minute {
		t.Errorf("conf.readTimeout %d is wrong", conf.readTimeout)
	}

	if conf.writeTimeout != 10*time.Minute {
		t.Errorf("conf.writeTimeout %d is wrong", conf.writeTimeout)
	}

	if conf.readBufferSize != 64*1024 {
		t.Errorf("conf.readBufferSize %d is wrong", conf.readBufferSize)
	}

	if conf.writeBufferSize != 64*1024 {
		t.Errorf("conf.writeBufferSize %d is wrong", conf.writeBufferSize)
	}
}

// go test -v -cover -run=^TestNewServerConfig$
func TestNewServerConfig(t *testing.T) {
	conf := newServerConfig("127.0.0.1:5837")

	if conf.address != "127.0.0.1:5837" {
		t.Errorf("conf.address %s is wrong", conf.address)
	}

	if conf.name != "127.0.0.1:5837" {
		t.Errorf("conf.name %s is wrong", conf.name)
	}

	if conf.readTimeout != 10*time.Minute {
		t.Errorf("conf.readTimeout %d is wrong", conf.readTimeout)
	}

	if conf.writeTimeout != 10*time.Minute {
		t.Errorf("conf.writeTimeout %d is wrong", conf.writeTimeout)
	}

	if conf.closeTimeout != time.Minute {
		t.Errorf("conf.closeTimeout %d is wrong", conf.closeTimeout)
	}

	if conf.readBufferSize != 16*1024 {
		t.Errorf("conf.readBufferSize %d is wrong", conf.readBufferSize)
	}

	if conf.writeBufferSize != 16*1024 {
		t.Errorf("conf.writeBufferSize %d is wrong", conf.writeBufferSize)
	}
}

// go test -v -cover -run=^TestConfigApplyOptions$
func TestConfigApplyOptions(t *testing.T) {
	conf := &Config{address: "127.0.0.1:5837"}

	conf.ApplyOptions([]Option{
		WithReadTimeout(time.Second),
		WithWriteTimeout(3 * time.Second),
		WithReadBufferSize(64),
		WithWriteBufferSize(512),
	})

	if conf.address != "127.0.0.1:5837" {
		t.Errorf("conf.address %s is wrong", conf.address)
	}

	if conf.readTimeout != time.Second {
		t.Errorf("conf.readTimeout %d is wrong", conf.readTimeout)
	}

	if conf.writeTimeout != 3*time.Second {
		t.Errorf("conf.writeTimeout %d is wrong", conf.writeTimeout)
	}

	if conf.readBufferSize != 64 {
		t.Errorf("config.readBufferSize %d is wrong", conf.readBufferSize)
	}

	if conf.writeBufferSize != 512 {
		t.Errorf("config.writeBufferSize %d is wrong", conf.writeBufferSize)
	}
}
