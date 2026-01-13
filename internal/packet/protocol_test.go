// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packet

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"testing"
)

// go test -v -cover -run=^TestDecode$
func TestDecode(t *testing.T) {
	type testCase struct {
		packetBytes []byte
		packet      Packet
		err         error
	}

	testCases := []testCase{
		{
			packetBytes: []byte{},
			packet:      Packet{},
			err:         io.EOF,
		},
		{
			packetBytes: make([]byte, HeaderBytes),
			packet:      Packet{},
			err:         ErrWrongMagic,
		},
		{
			packetBytes: []byte{0xC, 0x63, 0x8B, PacketTypeError, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
			packet:      Packet{},
			err:         io.ErrUnexpectedEOF,
		},
		{
			packetBytes: []byte{0xC, 0x63, 0x8B, PacketTypeError, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9},
			packet:      Packet{Magic: Magic, Type: PacketTypeError, Length: 0, Sequence: 9},
			err:         nil,
		},
		{
			packetBytes: []byte{0xC, 0x63, 0x8B, PacketTypeError, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 9, 6, 7, 8},
			packet:      Packet{Magic: Magic, Type: PacketTypeError, Length: 3, Sequence: 9, Data: []byte{6, 7, 8}},
			err:         nil,
		},
	}

	for _, testCase := range testCases {
		reader := bytes.NewBuffer(testCase.packetBytes)

		packet, err := Decode(reader)
		if err != testCase.err {
			t.Fatalf("input %+v: got %+v != want %+v", testCase.packetBytes, err, testCase.err)
		}

		got := fmt.Sprintf("%+v", packet)
		want := fmt.Sprintf("%+v", testCase.packet)
		if got != want {
			t.Fatalf("input %+v: got %+v != want %+v", testCase.packetBytes, got, want)
		}
	}
}

// go test -v -cover -run=^TestEncode$
func TestEncode(t *testing.T) {
	type testCase struct {
		packetBytes []byte
		packet      Packet
		err         error
	}

	testCases := []testCase{
		{
			packetBytes: []byte{},
			packet:      Packet{},
			err:         ErrWrongMagic,
		},
		{
			packetBytes: []byte{},
			packet:      Packet{Magic: Magic, Type: PacketTypeError, Length: 1, Sequence: 6},
			err:         ErrWrongLength,
		},
		{
			packetBytes: []byte{0xC, 0x63, 0x8B, PacketTypeError, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9},
			packet:      Packet{Magic: Magic, Type: PacketTypeError, Length: 0, Sequence: 9},
			err:         nil,
		},
		{
			packetBytes: []byte{0xC, 0x63, 0x8B, PacketTypeError, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 9, 6, 7, 8},
			packet:      Packet{Magic: Magic, Type: PacketTypeError, Length: 3, Sequence: 9, Data: []byte{6, 7, 8}},
			err:         nil,
		},
	}

	for _, testCase := range testCases {
		writer := bytes.NewBuffer(make([]byte, 0, 128))

		err := Encode(writer, testCase.packet)
		if err != testCase.err {
			t.Fatalf("input %+v: got %+v != want %+v", testCase.packet, err, testCase.err)
		}

		got := writer.Bytes()
		want := testCase.packetBytes
		if !slices.Equal(got, want) {
			t.Fatalf("input %+v: got %+v != want %+v", testCase.packet, got, want)
		}
	}
}
