// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pack

import (
	"bytes"
	"testing"
)

const (
	packetTypeTest PacketType = 1
)

// go test -v -cover -run=^TestReadPacket$
func TestReadPacket(t *testing.T) {
	type expect struct {
		packetType PacketType
		body       []byte
		err        error
	}

	cases := []struct {
		input  []byte
		expect expect
	}{
		{
			input: []byte{0xC, 0x63, 0x8B, packetTypeTest, 0, 0, 0, 2, 'o', 'k'},
			expect: expect{
				packetType: packetTypeTest,
				body:       []byte{'o', 'k'},
				err:        nil,
			},
		},
		{
			input: []byte{0xC, 0x63, 0x8B + 1, packetTypeTest, 0, 0, 0, 2, 'o', 'k'},
			expect: expect{
				packetType: 0,
				body:       nil,
				err:        ErrWrongMagicNumber,
			},
		},
		{
			input: []byte{0xC, 0x63, 0x8B, packetTypeTest, 0, 0, 0, 0},
			expect: expect{
				packetType: packetTypeTest,
				body:       nil,
				err:        nil,
			},
		},
		{
			input: []byte{0xC, 0x63, 0x8B, packetTypeTest, 0, 0, 0, 3, 'o', 'k'},
			expect: expect{
				packetType: packetTypeTest,
				body:       nil,
				err:        ErrReadSizeMismatch,
			},
		},
	}

	for i, oneCase := range cases {
		reader := bytes.NewReader(oneCase.input)

		packetType, body, err := readPacket(reader)
		if err != oneCase.expect.err {
			t.Fatalf("i %d, err %+v != oneCase.expect.err %+v", i, err, oneCase.expect.err)
			break
		}

		if packetType != oneCase.expect.packetType {
			t.Fatalf("i %d, packetType %+v != oneCase.expect.packetType %+v", i, packetType, oneCase.expect.packetType)
			break
		}

		if bytes.Compare(body, oneCase.expect.body) != 0 {
			t.Fatalf("i %d, body %+v != oneCase.expect.body %+v", i, body, oneCase.expect.body)
			break
		}
	}
}

// go test -v -cover -run=^TestWritePacket$
func TestWritePacket(t *testing.T) {
	type input struct {
		packetType PacketType
		body       []byte
	}

	type expect struct {
		err    error
		packet []byte
	}

	cases := []struct {
		input  input
		expect expect
	}{
		{
			input: input{
				packetType: packetTypeStandard,
				body:       []byte{'o', 'k'},
			},
			expect: expect{
				err:    nil,
				packet: []byte{0xC, 0x63, 0x8B, packetTypeStandard, 0, 0, 0, 2, 'o', 'k'},
			},
		},
		{
			input: input{
				packetType: packetTypeError,
				body:       []byte{'e', 'r', 'r'},
			},
			expect: expect{
				err:    nil,
				packet: []byte{0xC, 0x63, 0x8B, packetTypeError, 0, 0, 0, 3, 'e', 'r', 'r'},
			},
		},
	}

	buffer := bytes.NewBuffer(make([]byte, 0, 64))
	for i, oneCase := range cases {
		buffer.Reset()

		err := writePacket(buffer, oneCase.input.packetType, oneCase.input.body)
		if err != oneCase.expect.err {
			t.Fatalf("i %d, err == nil, err %+v != oneCase.expect.err %+v", i, err, oneCase.expect.err)
			break
		}

		if bytes.Compare(buffer.Bytes(), oneCase.expect.packet) != 0 {
			t.Fatalf("i %d, buffer %+v != oneCase.expect.packet %+v", i, buffer.Bytes(), oneCase.expect.packet)
			break
		}
	}
}
