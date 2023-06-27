// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"testing"
	"time"
)

// go test -v -cover -run=^TestWithClientReadWriteTimeout$
func TestWithClientReadWriteTimeout(t *testing.T) {
	conf := &Config{ReadTimeout: 0, WriteTimeout: 0}

	WithClientReadTimeout(time.Hour)(conf)
	WithClientWriteTimeout(time.Minute)(conf)

	if conf.ReadTimeout != time.Hour {
		t.Errorf("c.ReadTimeout %d is wrong", conf.ReadTimeout)
	}

	if conf.WriteTimeout != time.Minute {
		t.Errorf("c.WriteTimeout %d is wrong", conf.WriteTimeout)
	}
}

// go test -v -cover -run=^TestWithClientBufferSize$
func TestWithClientBufferSize(t *testing.T) {
	conf := &Config{ReadBufferSize: 0, WriteBufferSize: 0}

	WithClientReadBufferSize(64)(conf)
	WithClientWriteBufferSize(512)(conf)

	if conf.ReadBufferSize != 64 {
		t.Errorf("c.ReadBufferSize %d is wrong", conf.ReadBufferSize)
	}

	if conf.WriteBufferSize != 512 {
		t.Errorf("c.WriteBufferSize %d is wrong", conf.WriteBufferSize)
	}
}

// go test -v -cover -run=^TestWithServerName$
func TestWithServerName(t *testing.T) {
	name := "test-name"

	conf := &Config{Name: ""}
	WithServerName(name)(conf)

	if conf.Name != name {
		t.Errorf("c.Name %s is wrong", conf.Name)
	}
}

// go test -v -cover -run=^TestWithServerReadWriteTimeout$
func TestWithServerReadWriteTimeout(t *testing.T) {
	conf := &Config{ReadTimeout: 0, WriteTimeout: 0}

	WithServerReadTimeout(time.Hour)(conf)
	WithServerWriteTimeout(time.Minute)(conf)

	if conf.ReadTimeout != time.Hour {
		t.Errorf("c.ReadTimeout %d is wrong", conf.ReadTimeout)
	}

	if conf.WriteTimeout != time.Minute {
		t.Errorf("c.WriteTimeout %d is wrong", conf.WriteTimeout)
	}
}

// go test -v -cover -run=^TestWithServerCloseTimeout$
func TestWithServerCloseTimeout(t *testing.T) {
	conf := &Config{CloseTimeout: 0}

	WithServerCloseTimeout(time.Hour)(conf)

	if conf.CloseTimeout != time.Hour {
		t.Errorf("c.CloseTimeout %d is wrong", conf.CloseTimeout)
	}
}

// go test -v -cover -run=^TestWithServerBufferSize$
func TestWithServerBufferSize(t *testing.T) {
	conf := &Config{ReadBufferSize: 0, WriteBufferSize: 0}

	WithServerReadBufferSize(64)(conf)
	WithServerWriteBufferSize(512)(conf)

	if conf.ReadBufferSize != 64 {
		t.Errorf("c.ReadBufferSize %d is wrong", conf.ReadBufferSize)
	}

	if conf.WriteBufferSize != 512 {
		t.Errorf("c.WriteBufferSize %d is wrong", conf.WriteBufferSize)
	}
}
