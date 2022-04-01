// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"encoding/binary"
	"errors"
	"net"
	"testing"
	"time"
)

const (
	packetTypeTest PacketType = 1
)

var (
	errTestRequestFailed = errors.New("vex: test request failed")
)

// go test -v -cover -run=^TestNewServer$
func TestNewServer(t *testing.T) {
	address := "127.0.0.1:5837"

	server := NewServer()
	server.RegisterPacketHandler(packetTypeTest, func(requestBody []byte) (responseBody []byte, err error) {
		if len(requestBody) <= 0 {
			return nil, errTestRequestFailed
		}
		return requestBody, nil
	})
	defer server.Close()

	go func() {
		err := server.ListenAndServe("tcp", address)
		if err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(time.Second)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	// failed test
	_, err = conn.Write([]byte{
		protocolVersion, packetTypeTest, 0, 0, 0, 0,
	})
	if err != nil {
		t.Error(err)
	}

	buffer := make([]byte, 128)
	n, err := conn.Read(buffer)
	if err != nil {
		t.Error(err)
	}

	if buffer[0] != protocolVersion {
		t.Errorf("protocol version %d is wrong!", buffer[0])
	}

	if buffer[1] != packetTypeTest {
		t.Errorf("tag %d is wrong!", buffer[1])
	}

	bodySize := binary.BigEndian.Uint32(buffer[2:6])
	if bodySize != uint32(len([]byte(errTestRequestFailed.Error()))) {
		t.Errorf("body size %d is wrong!", bodySize)
	}

	buffer = buffer[6:n]
	if string(buffer) != errTestRequestFailed.Error() {
		t.Errorf("body %s is wrong!", string(buffer))
	}

	// successful test
	_, err = conn.Write([]byte{
		protocolVersion, packetTypeTest, 0, 0, 0, 9, 'k', 'e', 'y', ' ', 'v', 'a', 'l', 'u', 'e',
	})
	if err != nil {
		t.Error(err)
	}

	buffer = make([]byte, 128)
	n, err = conn.Read(buffer)
	if err != nil {
		t.Error(err)
	}

	if buffer[0] != protocolVersion {
		t.Errorf("protocol version %d is wrong!", buffer[0])
	}

	if buffer[1] != packetTypeTest {
		t.Errorf("tag %d is wrong!", buffer[1])
	}

	bodySize = binary.BigEndian.Uint32(buffer[2:6])
	if bodySize != uint32(9) {
		t.Errorf("body size %d is wrong!", bodySize)
	}

	buffer = buffer[6:n]
	if string(buffer) != "key value" {
		t.Errorf("body %s is wrong!", string(buffer))
	}
}
