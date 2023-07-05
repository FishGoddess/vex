// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"testing"
	"time"
	"unsafe"
)

// go test -v -cover -run=^TestWithName$
func TestWithName(t *testing.T) {
	name := "test-name"

	conf := &Config{name: ""}
	WithName(name)(conf)

	if conf.name != name {
		t.Errorf("c.name %s is wrong", conf.name)
	}
}

// go test -v -cover -run=^TestWithReadTimeout$
func TestWithReadTimeout(t *testing.T) {
	conf := &Config{readTimeout: 0}
	WithReadTimeout(time.Hour)(conf)

	if conf.readTimeout != time.Hour {
		t.Errorf("c.readTimeout %d is wrong", conf.readTimeout)
	}
}

// go test -v -cover -run=^TestWithWriteTimeout$
func TestWithWriteTimeout(t *testing.T) {
	conf := &Config{writeTimeout: 0}
	WithWriteTimeout(time.Minute)(conf)

	if conf.writeTimeout != time.Minute {
		t.Errorf("c.writeTimeout %d is wrong", conf.writeTimeout)
	}
}

// go test -v -cover -run=^TestWithCloseTimeout$
func TestWithCloseTimeout(t *testing.T) {
	conf := &Config{closeTimeout: 0}

	WithCloseTimeout(time.Hour)(conf)

	if conf.closeTimeout != time.Hour {
		t.Errorf("c.closeTimeout %d is wrong", conf.closeTimeout)
	}
}

// go test -v -cover -run=^TestWithBufferSize$
func TestWithBufferSize(t *testing.T) {
	conf := &Config{readBufferSize: 0}
	WithReadBufferSize(64)(conf)

	if conf.readBufferSize != 64 {
		t.Errorf("c.readBufferSize %d is wrong", conf.readBufferSize)
	}
}

// go test -v -cover -run=^TestWithWriteBufferSize$
func TestWithWriteBufferSize(t *testing.T) {
	conf := &Config{writeBufferSize: 0}
	WithWriteBufferSize(512)(conf)

	if conf.writeBufferSize != 512 {
		t.Errorf("c.writeBufferSize %d is wrong", conf.writeBufferSize)
	}
}

// go test -v -cover -run=^TestWithOnConnected$
func TestWithOnConnected(t *testing.T) {
	onConnectedFunc := func(clientAddress string, serverAddress string) {}

	conf := &Config{onConnectedFunc: nil}
	WithOnConnected(onConnectedFunc)(conf)

	if unsafe.Sizeof(conf.onConnectedFunc) != unsafe.Sizeof(onConnectedFunc) {
		t.Errorf("c.onConnectedFunc %d is wrong", unsafe.Sizeof(conf.onConnectedFunc))
	}
}

// go test -v -cover -run=^TestWithOnDisconnected$
func TestWithOnDisconnected(t *testing.T) {
	onDisconnectedFunc := func(clientAddress string, serverAddress string) {}

	conf := &Config{onDisconnectedFunc: nil}
	WithOnDisconnected(onDisconnectedFunc)(conf)

	if unsafe.Sizeof(conf.onDisconnectedFunc) != unsafe.Sizeof(onDisconnectedFunc) {
		t.Errorf("c.onDisconnectedFunc %d is wrong", unsafe.Sizeof(conf.onDisconnectedFunc))
	}
}

// go test -v -cover -run=^TestWithBeforeServing$
func TestWithBeforeServing(t *testing.T) {
	beforeServingFunc := func(address string) {}

	conf := &Config{beforeServingFunc: nil}
	WithBeforeServing(beforeServingFunc)(conf)

	if unsafe.Sizeof(conf.beforeServingFunc) != unsafe.Sizeof(beforeServingFunc) {
		t.Errorf("c.beforeServingFunc %d is wrong", unsafe.Sizeof(conf.beforeServingFunc))
	}
}

// go test -v -cover -run=^TestWithAfterServing$
func TestWithAfterServing(t *testing.T) {
	afterServingFunc := func(address string, err error) {}

	conf := &Config{afterServingFunc: nil}
	WithAfterServing(afterServingFunc)(conf)

	if unsafe.Sizeof(conf.afterServingFunc) != unsafe.Sizeof(afterServingFunc) {
		t.Errorf("c.afterServingFunc %d is wrong", unsafe.Sizeof(conf.afterServingFunc))
	}
}

// go test -v -cover -run=^TestWithBeforeHandling$
func TestWithBeforeHandling(t *testing.T) {
	beforeHandlingFunc := func(ctx *Context) {}

	conf := &Config{beforeHandlingFunc: nil}
	WithBeforeHandling(beforeHandlingFunc)(conf)

	if unsafe.Sizeof(conf.beforeHandlingFunc) != unsafe.Sizeof(beforeHandlingFunc) {
		t.Errorf("c.beforeHandlingFunc %d is wrong", unsafe.Sizeof(conf.beforeHandlingFunc))
	}
}

// go test -v -cover -run=^TestWithAfterHandling$
func TestWithAfterHandling(t *testing.T) {
	afterHandlingFunc := func(ctx *Context) {}

	conf := &Config{afterHandlingFunc: nil}
	WithAfterHandling(afterHandlingFunc)(conf)

	if unsafe.Sizeof(conf.afterHandlingFunc) != unsafe.Sizeof(afterHandlingFunc) {
		t.Errorf("c.afterHandlingFunc %d is wrong", unsafe.Sizeof(conf.afterHandlingFunc))
	}
}

// go test -v -cover -run=^TestWithBeforeClosing$
func TestWithBeforeClosing(t *testing.T) {
	beforeClosingFunc := func(address string) {}

	conf := &Config{beforeClosingFunc: nil}
	WithBeforeClosing(beforeClosingFunc)(conf)

	if unsafe.Sizeof(conf.beforeClosingFunc) != unsafe.Sizeof(beforeClosingFunc) {
		t.Errorf("c.beforeClosingFunc %d is wrong", unsafe.Sizeof(conf.beforeClosingFunc))
	}
}

// go test -v -cover -run=^TestWithAfterClosing$
func TestWithAfterClosing(t *testing.T) {
	afterClosingFunc := func(address string, err error) {}

	conf := &Config{afterClosingFunc: nil}
	WithAfterClosing(afterClosingFunc)(conf)

	if unsafe.Sizeof(conf.afterClosingFunc) != unsafe.Sizeof(afterClosingFunc) {
		t.Errorf("c.afterClosingFunc %d is wrong", unsafe.Sizeof(conf.afterClosingFunc))
	}
}
