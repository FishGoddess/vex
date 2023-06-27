// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"io"
	"testing"
)

// go test -v -cover -run=^TestLogDebug$
func TestLogDebug(t *testing.T) {
	if LogDebug == nil {
		t.Error("LogDebug == nil")
	}

	logDebug("...%d...", 1)

	defer func() {
		r := recover()
		if r != nil {
			t.Error(t)
		}
	}()

	LogDebug = nil
	logDebug("...%d...", 1)
}

// go test -v -cover -run=^TestLogInfo$
func TestLogInfo(t *testing.T) {
	if LogInfo == nil {
		t.Error("LogInfo == nil")
	}

	logInfo("...%d...", 1)

	defer func() {
		r := recover()
		if r != nil {
			t.Error(t)
		}
	}()

	LogInfo = nil
	logInfo("...%d...", 1)
}

// go test -v -cover -run=^TestLogError$
func TestLogError(t *testing.T) {
	if LogError == nil {
		t.Error("LogError == nil")
	}

	logError(io.EOF, "...%d...", 1)

	defer func() {
		r := recover()
		if r != nil {
			t.Error(t)
		}
	}()

	LogError = nil
	logError(io.ErrUnexpectedEOF, "...%d...", 1)
}
