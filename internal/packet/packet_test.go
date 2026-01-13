// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packet

import (
	"slices"
	"testing"
)

// go test -v -cover -run=^TestPacketWith$
func TestPacketWith(t *testing.T) {
	data := []byte{1, 2, 5, 8, 0}

	packet := new(Packet)
	packet.With(data)

	if packet.Length != uint32(len(data)) {
		t.Fatalf("got %d != want %d", packet.Length, uint32(len(data)))
	}

	if !slices.Equal(packet.Data, data) {
		t.Fatalf("got %d != want %d", packet.Length, uint32(len(data)))
	}
}
