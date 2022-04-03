// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import "testing"

// go test -v -cover -run=^TestLog$
func TestLog(t *testing.T) {
	if Log == nil {
		t.Error("Log == nil")
	}

	defer func() {
		r := recover()
		if r != nil {
			t.Error(t)
		}
	}()

	Log = func(format string, v ...interface{}) {}
	log("...")
}

// go test -v -cover -run=^TestDial$
func TestDial(t *testing.T) {
	if Dial == nil {
		t.Error("Dial == nil")
	}

	defer func() {
		r := recover()
		if r != nil {
			t.Error(t)
		}
	}()

	_, err := dial("tcp", "...")
	if err == nil {
		t.Error("err == nil")
	}
}

// go test -v -cover -run=^TestMakeBytes$
func TestMakeBytes(t *testing.T) {
	initialized := 64

	result := makeBytes(int32(initialized))
	if len(result) != initialized {
		t.Errorf("len(result) %d != initialized %d", len(result), initialized)
	}

	if cap(result) != initialized {
		t.Errorf("cap(result) %d != initialized %d", cap(result), initialized)
	}
}
