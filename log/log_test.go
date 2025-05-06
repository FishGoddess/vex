// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
	"io"
	"testing"
)

// go test -v -cover -run=^TestDisableDebug$
func TestDisableDebug(t *testing.T) {
	old := DebugFunc
	defer func() {
		DebugFunc = old
	}()

	DisableDebug()
	if DebugFunc != nil {
		t.Fatal("disable debug failed")
	}
}

// go test -v -cover -run=^TestDisableInfo$
func TestDisableInfo(t *testing.T) {
	old := InfoFunc
	defer func() {
		InfoFunc = old
	}()

	DisableInfo()
	if InfoFunc != nil {
		t.Fatal("disable info failed")
	}
}

// go test -v -cover -run=^TestDisableError$
func TestDisableError(t *testing.T) {
	old := ErrorFunc
	defer func() {
		ErrorFunc = old
	}()

	DisableError()
	if ErrorFunc != nil {
		t.Fatal("disable error failed")
	}
}

// go test -v -cover -run=^TestDebug$
func TestDebug(t *testing.T) {
	if DebugFunc == nil {
		t.Fatal("DebugFunc == nil")
	}

	Debug("...%d...", 1)

	defer func() {
		r := recover()
		if r != nil {
			t.Fatal(t)
		}
	}()

	DebugFunc = nil
	Debug("...%d...", 1)
}

// go test -v -cover -run=^TestInfo$
func TestInfo(t *testing.T) {
	if InfoFunc == nil {
		t.Fatal("InfoFunc == nil")
	}

	Info("...%d...", 1)

	defer func() {
		r := recover()
		if r != nil {
			t.Fatal(t)
		}
	}()

	InfoFunc = nil
	Info("...%d...", 1)
}

// go test -v -cover -run=^TestError$
func TestError(t *testing.T) {
	if ErrorFunc == nil {
		t.Fatal("ErrorFunc == nil")
	}

	Error(io.EOF, "...%d...", 1)

	defer func() {
		r := recover()
		if r != nil {
			t.Fatal(t)
		}
	}()

	ErrorFunc = nil
	Error(io.ErrUnexpectedEOF, "...%d...", 1)
}
