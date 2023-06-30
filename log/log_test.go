// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
	"io"
	"testing"
)

// go test -v -cover -run=^TestDebugFunc$
func TestDebugFunc(t *testing.T) {
	if DebugFunc == nil {
		t.Error("DebugFunc == nil")
	}

	Debug("...%d...", 1)

	defer func() {
		r := recover()
		if r != nil {
			t.Error(t)
		}
	}()

	DebugFunc = nil
	Debug("...%d...", 1)
}

// go test -v -cover -run=^TestInfoFunc$
func TestInfoFunc(t *testing.T) {
	if InfoFunc == nil {
		t.Error("InfoFunc == nil")
	}

	Info("...%d...", 1)

	defer func() {
		r := recover()
		if r != nil {
			t.Error(t)
		}
	}()

	InfoFunc = nil
	Info("...%d...", 1)
}

// go test -v -cover -run=^TestErrorFunc$
func TestErrorFunc(t *testing.T) {
	if ErrorFunc == nil {
		t.Error("ErrorFunc == nil")
	}

	Error(io.EOF, "...%d...", 1)

	defer func() {
		r := recover()
		if r != nil {
			t.Error(t)
		}
	}()

	ErrorFunc = nil
	Error(io.ErrUnexpectedEOF, "...%d...", 1)
}
