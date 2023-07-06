// Copyright 2023 FishGoddess. All rights reserved.
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
		t.Error("disable debug failed")
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
		t.Error("disable info failed")
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
		t.Error("disable error failed")
	}
}

// go test -v -cover -run=^TestDebug$
func TestDebug(t *testing.T) {
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

// go test -v -cover -run=^TestInfo$
func TestInfo(t *testing.T) {
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

// go test -v -cover -run=^TestError$
func TestError(t *testing.T) {
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
