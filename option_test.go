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

	conf := &Config{Name: ""}
	WithName(name)(conf)

	if conf.Name != name {
		t.Errorf("c.Name %s is wrong", conf.Name)
	}
}

// go test -v -cover -run=^TestWithReadTimeout$
func TestWithReadTimeout(t *testing.T) {
	conf := &Config{ReadTimeout: 0}
	WithReadTimeout(time.Hour)(conf)

	if conf.ReadTimeout != time.Hour {
		t.Errorf("c.ReadTimeout %d is wrong", conf.ReadTimeout)
	}
}

// go test -v -cover -run=^TestWithWriteTimeout$
func TestWithWriteTimeout(t *testing.T) {
	conf := &Config{WriteTimeout: 0}
	WithWriteTimeout(time.Minute)(conf)

	if conf.WriteTimeout != time.Minute {
		t.Errorf("c.WriteTimeout %d is wrong", conf.WriteTimeout)
	}
}

// go test -v -cover -run=^TestWithCloseTimeout$
func TestWithCloseTimeout(t *testing.T) {
	conf := &Config{CloseTimeout: 0}

	WithCloseTimeout(time.Hour)(conf)

	if conf.CloseTimeout != time.Hour {
		t.Errorf("c.CloseTimeout %d is wrong", conf.CloseTimeout)
	}
}

// go test -v -cover -run=^TestWithBufferSize$
func TestWithBufferSize(t *testing.T) {
	conf := &Config{ReadBufferSize: 0}
	WithReadBufferSize(64)(conf)

	if conf.ReadBufferSize != 64 {
		t.Errorf("c.ReadBufferSize %d is wrong", conf.ReadBufferSize)
	}
}

// go test -v -cover -run=^TestWithWriteBufferSize$
func TestWithWriteBufferSize(t *testing.T) {
	conf := &Config{WriteBufferSize: 0}
	WithWriteBufferSize(512)(conf)

	if conf.WriteBufferSize != 512 {
		t.Errorf("c.WriteBufferSize %d is wrong", conf.WriteBufferSize)
	}
}
