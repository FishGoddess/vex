// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packet

import (
	"slices"
	"testing"
)

// go test -v -cover -run=^TestNew$
func TestNew(t *testing.T) {
	sequence := uint64(12580)
	data := []byte{1, 2, 1, 3, 8}

	packet := New(PacketTypeError, sequence, data)
	if packet.magic != magic {
		t.Fatalf("got %d != want %d", packet.magic, magic)
	}

	if packet.ptype != PacketTypeError {
		t.Fatalf("got %d != want %d", packet.ptype, PacketTypeError)
	}

	length := uint32(len(data))
	if packet.length != length {
		t.Fatalf("got %d != want %d", packet.length, length)
	}

	if packet.sequence != sequence {
		t.Fatalf("got %d != want %d", packet.sequence, sequence)
	}

	if !slices.Equal(packet.data, data) {
		t.Fatalf("got %+v != want %+v", packet.data, data)
	}
}
