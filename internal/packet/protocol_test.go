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

// go test -v -cover -run=^TestReadPacket$
func TestReadPacket(t *testing.T) {
	maxDataBytes = 8

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
			packetBytes: make([]byte, headerBytes),
			packet:      Packet{},
			err:         errWrongMagic,
		},
		{
			packetBytes: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0x77, 0x14, 0x30, 0xCB, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			packet:      Packet{},
			err:         io.ErrUnexpectedEOF,
		},
		{
			packetBytes: []byte{0, 0, 0, 0, 0, 0, 0, 5, 0x77, 0x14, 0x30, 0xCB, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 10, '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'},
			packet:      Packet{id: 5, magic: magic, flags: 1, length: 10, data: []byte{}},
			err:         errDataTooLarge,
		},
		{
			packetBytes: []byte{0, 0, 0, 0, 0, 0, 0, 5, 0x77, 0x14, 0x30, 0xCB, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
			packet:      Packet{id: 5, magic: magic, flags: 1, length: 0},
			err:         nil,
		},
		{
			packetBytes: []byte{0, 0, 0, 0, 0, 0, 0, 5, 0x77, 0x14, 0x30, 0xCB, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 3, 'A', 'B', 'C'},
			packet:      Packet{id: 5, magic: magic, flags: 1, length: 3, data: []byte("ABC")},
			err:         nil,
		},
		{
			packetBytes: []byte{0, 0, 0, 0, 0, 0, 0, 5, 0x77, 0x14, 0x30, 0xCB, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 'A', 'B', 'C'},
			packet:      Packet{id: 5, magic: magic, flags: 0, length: 3, data: []byte("ABC")},
			err:         nil,
		},
	}

	for _, testCase := range testCases {
		reader := bytes.NewBuffer(testCase.packetBytes)

		packet, err := ReadPacket(reader)
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

// go test -v -cover -run=^TestWritePacket$
func TestWritePacket(t *testing.T) {
	maxDataBytes = 8

	type testCase struct {
		packetBytes []byte
		packet      Packet
		err         error
	}

	testCases := []testCase{
		{
			packetBytes: []byte{},
			packet:      Packet{},
			err:         errWrongMagic,
		},
		{
			packetBytes: []byte{},
			packet:      Packet{id: 5, magic: magic, flags: 1, length: 3},
			err:         errWrongLength,
		},
		{
			packetBytes: []byte{},
			packet:      Packet{id: 5, magic: magic, flags: 1, length: 10, data: []byte("0123456789")},
			err:         errDataTooLarge,
		},
		{
			packetBytes: []byte{0, 0, 0, 0, 0, 0, 0, 5, 0x77, 0x14, 0x30, 0xCB, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
			packet:      Packet{id: 5, magic: magic, flags: 1, length: 0},
			err:         nil,
		},
		{
			packetBytes: []byte{0, 0, 0, 0, 0, 0, 0, 5, 0x77, 0x14, 0x30, 0xCB, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 3, 'A', 'B', 'C'},
			packet:      Packet{id: 5, magic: magic, flags: 1, length: 3, data: []byte("ABC")},
			err:         nil,
		},
		{
			packetBytes: []byte{0, 0, 0, 0, 0, 0, 0, 5, 0x77, 0x14, 0x30, 0xCB, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 'A', 'B', 'C'},
			packet:      Packet{id: 5, magic: magic, flags: 0, length: 3, data: []byte("ABC")},
			err:         nil,
		},
	}

	for _, testCase := range testCases {
		writer := bytes.NewBuffer(make([]byte, 0, 128))

		err := WritePacket(writer, testCase.packet)
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
