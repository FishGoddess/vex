// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packet

import (
	"io"
	"slices"
	"testing"
)

// go test -v -cover -run=^TestNew$
func TestNew(t *testing.T) {
	got := New(1234567890)

	if got.id != 1234567890 {
		t.Fatalf("got %d is wrong", got.id)
	}

	if got.magic != magic {
		t.Fatalf("got %d is wrong", got.magic)
	}
}

// go test -v -cover -run=^TestPacketFlagSet$
func TestPacketFlagSet(t *testing.T) {
	flag1 := uint64(2)
	flag2 := uint64(16)

	packet := Packet{flags: flag1 + flag2}
	if !packet.flagSet(flag1) {
		t.Fatalf("flag %d not set", flag1)
	}

	if !packet.flagSet(flag2) {
		t.Fatalf("flag %d not set", flag2)
	}
}

// go test -v -cover -run=^TestPacketID$
func TestPacketID(t *testing.T) {
	packet := Packet{id: 1234567890}

	if packet.id != 1234567890 {
		t.Fatalf("got %+v is wrong", packet.id)
	}
}

// go test -v -cover -run=^TestPacketData$
func TestPacketData(t *testing.T) {
	data := []byte("欲买桂花同载酒")
	packet := Packet{flags: 0, data: data}

	got, err := packet.Data()
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(got, data) {
		t.Fatalf("got %+v != want %+v", got, data)
	}

	data = []byte(io.EOF.Error())
	packet = Packet{flags: flagError, data: data}

	_, err = packet.Data()
	if err == nil {
		t.Fatalf("packet data returns nil error")
	}

	got = []byte(err.Error())
	if !slices.Equal(got, data) {
		t.Fatalf("got %+v != want %+v", got, data)
	}
}

// go test -v -cover -run=^TestPacketSetFlag$
func TestPacketSetFlag(t *testing.T) {
	flag1 := uint64(2)
	flag2 := uint64(16)

	packet := Packet{flags: 0}
	packet.setFlag(flag1)
	packet.setFlag(flag2)

	got := packet.flags
	want := flag1 + flag2
	if got != want {
		t.Fatalf("got %d != want %d", got, want)
	}
}

// go test -v -cover -run=^TestPacketSetData$
func TestPacketSetData(t *testing.T) {
	data := []byte("终不似少年游")

	packet := Packet{length: 0, data: nil}
	packet.SetData(data)

	got := packet.data
	if !slices.Equal(got, data) {
		t.Fatalf("got %+v != want %+v", got, data)
	}
}

// go test -v -cover -run=^TestPacketSetError$
func TestPacketSetError(t *testing.T) {
	err := io.EOF

	packet := Packet{length: 0, data: nil}
	packet.SetError(err)

	got := packet.data
	want := []byte(err.Error())
	if !slices.Equal(got, want) {
		t.Fatalf("got %+v != want %+v", got, want)
	}
}
