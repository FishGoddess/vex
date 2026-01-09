// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packet

import (
	"slices"
	"testing"
)

// go test -v -cover -run=^TestPacket$
func TestPacket(t *testing.T) {
	sequence := uint64(12580)
	data := []byte{1, 2, 1, 3, 8}

	var packet Packet
	packet.magic = magic
	packet.sequence = sequence
	packet.SetType(PacketTypeError)
	packet.SetData(data)

	if packet.magic != magic {
		t.Fatalf("got %d != want %d", packet.magic, magic)
	}

	if packet.Type() != packet.ptype {
		t.Fatalf("got %d != want %d", packet.Type(), packet.ptype)
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

	if !slices.Equal(packet.Data(), packet.data) {
		t.Fatalf("got %+v != want %+v", packet.data, packet.data)
	}

	if !slices.Equal(packet.data, data) {
		t.Fatalf("got %+v != want %+v", packet.data, data)
	}
}
