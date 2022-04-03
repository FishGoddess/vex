// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.¬

package vex

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
			input: []byte{0x7, 0x55, 0xDD, 0x8C, protocolVersion, packetTypeTest, 0, 0, 0, 2, 'o', 'k'},
			expect: expect{
				packetType: packetTypeTest,
				body:       []byte{'o', 'k'},
				err:        nil,
			},
		},
		{
			input: []byte{0x7, 0x55, 0xDD, 0x8C + 1, protocolVersion, packetTypeTest, 0, 0, 0, 2, 'o', 'k'},
			expect: expect{
				packetType: 0,
				body:       nil,
				err:        errMagicMismatch,
			},
		},
		{
			input: []byte{0x7, 0x55, 0xDD, 0x8C, protocolVersion + 1, packetTypeTest, 0, 0, 0, 2, 'o', 'k'},
			expect: expect{
				packetType: 0,
				body:       nil,
				err:        errProtocolMismatch,
			},
		},
		{
			input: []byte{0x7, 0x55, 0xDD, 0x8C, protocolVersion, packetTypeTest, 0, 0, 0, 3, 'o', 'k'},
			expect: expect{
				packetType: packetTypeTest,
				body:       nil,
				err:        errReadSizeMismatch,
			},
		},
	}

	for _, oneCase := range cases {
		packetType, body, err := readPacket(bytes.NewReader(oneCase.input))
		if packetType != oneCase.expect.packetType {
			t.Errorf("packetType %+v != oneCase.expect.packetType %+v", packetType, oneCase.expect.packetType)
		}

		if bytes.Compare(body, oneCase.expect.body) != 0 {
			t.Errorf("body %+v != oneCase.expect.body %+v", body, oneCase.expect.body)
		}

		if err != oneCase.expect.err {
			t.Errorf("err %+v != oneCase.expect.err %+v", err, oneCase.expect.err)
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
				packetType: packetTypeOK,
				body:       []byte{'o', 'k'},
			},
			expect: expect{
				err:    nil,
				packet: []byte{0x7, 0x55, 0xDD, 0x8C, protocolVersion, packetTypeOK, 0, 0, 0, 2, 'o', 'k'},
			},
		},
	}

	buffer := bytes.NewBuffer(make([]byte, 0, 64))
	for _, oneCase := range cases {
		buffer.Reset()

		err := writePacket(buffer, oneCase.input.packetType, oneCase.input.body)
		if err != oneCase.expect.err {
			t.Errorf("err == nil, err %+v != oneCase.expect.err %+v", err, oneCase.expect.err)
		}

		if bytes.Compare(buffer.Bytes(), oneCase.expect.packet) != 0 {
			t.Errorf("buffer %+v != oneCase.expect.packet %+v", buffer.Bytes(), oneCase.expect.packet)
		}
	}
}
