// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"testing"
	"time"
)

// go test -v -cover -run=^TestWithName$
func TestWithName(t *testing.T) {
	name := "test-name"

	conf := &ServerConfig{Name: ""}
	WithName(name)(conf)

	if conf.Name != name {
		t.Errorf("c.Name %s is wrong", conf.Name)
	}
}

// go test -v -cover -run=^TestWithReadWriteTimeout$
func TestWithReadWriteTimeout(t *testing.T) {
	conf := &ServerConfig{ReadTimeout: 0, WriteTimeout: 0}

	WithReadTimeout(time.Hour)(conf)
	WithWriteTimeout(time.Minute)(conf)

	if conf.ReadTimeout != time.Hour {
		t.Errorf("c.ReadTimeout %d is wrong", conf.ReadTimeout)
	}

	if conf.WriteTimeout != time.Minute {
		t.Errorf("c.WriteTimeout %d is wrong", conf.WriteTimeout)
	}
}

// go test -v -cover -run=^TestWithCloseTimeout$
func TestWithCloseTimeout(t *testing.T) {
	conf := &ServerConfig{CloseTimeout: 0}

	WithCloseTimeout(time.Hour)(conf)

	if conf.CloseTimeout != time.Hour {
		t.Errorf("c.CloseTimeout %d is wrong", conf.CloseTimeout)
	}
}

// go test -v -cover -run=^TestWithBufferSize$
func TestWithBufferSize(t *testing.T) {
	conf := &ServerConfig{ReadBufferSize: 0, WriteBufferSize: 0}

	WithReadBufferSize(64)(conf)
	WithWriteBufferSize(512)(conf)

	if conf.ReadBufferSize != 64 {
		t.Errorf("c.ReadBufferSize %d is wrong", conf.ReadBufferSize)
	}

	if conf.WriteBufferSize != 512 {
		t.Errorf("c.WriteBufferSize %d is wrong", conf.WriteBufferSize)
	}
}
