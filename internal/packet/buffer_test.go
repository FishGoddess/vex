// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packet

import "testing"

// go test -v -cover -run=^TestBufferData$
func TestBufferData(t *testing.T) {
	type testCase struct {
		length   int
		capacity int
	}

	testCases := []testCase{
		{length: 0, capacity: 0},
		{length: 0, capacity: 16},
		{length: 16, capacity: 0},
		{length: 16, capacity: 16},
	}

	for _, testCase := range testCases {
		buff := buffer{data: make([]byte, 0)}
		data := buff.Data(testCase.length, testCase.capacity)

		if len(data) < testCase.length {
			t.Fatalf("got %d < want %d", len(data), testCase.length)
		}

		if cap(data) < testCase.capacity {
			t.Fatalf("got %d < want %d", cap(data), testCase.capacity)
		}
	}
}

// go test -v -cover -run=^TestAcquireBuffer$
func TestAcquireBuffer(t *testing.T) {
	buff := acquireBuffer()

	if len(buff.data) != 0 {
		t.Fatalf("got %d != want %d", len(buff.data), 0)
	}

	if cap(buff.data) < 1*KB {
		t.Fatalf("got %d != want %d", cap(buff.data), 1*KB)
	}
}

// go test -v -cover -run=^TestReleaseBuffer$
func TestReleaseBuffer(t *testing.T) {
	buff := acquireBuffer()

	if !releaseBuffer(buff) {
		t.Fatalf("buff %d not released", cap(buff.data))
	}

	buff.data = make([]byte, 64*KB)

	if releaseBuffer(buff) {
		t.Fatalf("buff %d released", cap(buff.data))
	}
}
