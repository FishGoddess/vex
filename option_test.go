// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"
)

// go test -v -cover -run=^TestConfigApply$
func TestConfigApply(t *testing.T) {
	var conf config

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	opt1 := func(c *config) { c.logger = logger }
	opt2 := func(c *config) { c.dialTimeout = 2 }

	got := *conf.apply(opt1, opt2)
	want := config{logger: logger, dialTimeout: 2}
	if got != want {
		t.Fatalf("got %+v != want %+v", got, want)
	}
}

// go test -v -cover -run=^TestWithLogger$
func TestWithLogger(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	conf := &config{logger: nil}
	WithLogger(logger)(conf)

	got := fmt.Sprintf("%p", conf.logger)
	want := fmt.Sprintf("%p", logger)
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

// go test -v -cover -run=^TestWithDialTimeout$
func TestWithDialTimeout(t *testing.T) {
	timeout := time.Millisecond

	conf := &config{dialTimeout: 0}
	WithDialTimeout(timeout)(conf)

	got := conf.dialTimeout
	want := timeout
	if got != want {
		t.Fatalf("got %d != want %d", got, want)
	}
}
