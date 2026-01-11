// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
)

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
